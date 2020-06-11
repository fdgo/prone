package sms

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/shopspring/decimal"
)

var (
	SQSClient *sqs.SQS
	SNSClient *sns.SNS
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Conf.ActionSNSConf.Region),
		Credentials: credentials.NewStaticCredentials(config.Conf.ActionSNSConf.AccessKey, config.Conf.ActionSNSConf.AccessSecret, ""),
	})
	if err != nil {
		panic(err)
	}
	SQSClient = sqs.New(sess)
	SNSClient = sns.New(sess)

}

func SendSMS(msg *domain.SMSMessage, smsTemplate *domain.SMSTemplate) error {
	switch msg.Service {
	case domain.SMS_SERVICE_AWSSNS:
		return SendSMSByAWS(msg, smsTemplate)
	case domain.SMS_SERVICE_NEC:
		msg.TemplateID = smsTemplate.TemplateID
		return SendSMSByNEC(msg)
	default:
		if msg.PhoneNumer.AreaCode() == "+86" {
			msg.TemplateID = smsTemplate.TemplateID
		} else {
			msg.TemplateID = smsTemplate.TemplateInternationalID
		}
		return SendSMSBySendCloud(msg)
	}
}

func SendSMSAsnyc(msg *domain.SMSMessage) error {
	var (
		delaySeconds = int64(0)
		queueUrl     = config.Conf.SQSConf.SMSNotifyUrl
	)
	body, err := json.Marshal(&msg)
	if err != nil {
		return err
	}
	msgBody := string(body)
	input := sqs.SendMessageInput{
		DelaySeconds: &delaySeconds,
		MessageBody:  &msgBody,
		QueueUrl:     &queueUrl,
	}

	out, err := SQSClient.SendMessage(&input)
	if err != nil {
		return err
	}
	loggers.Info.Printf("Send sms message to sms notify queue id:%s md5:%s", *out.MessageId, *out.MD5OfMessageBody)
	return nil
}

func SendSMSBySendCloud(msg *domain.SMSMessage) error {
	RequestURI := config.Conf.SendCloudSMSConf.ApiUrl
	vars, err := json.Marshal(msg.Source)
	if err != nil {
		return err
	}
	smsKey := config.Conf.SendCloudSMSConf.ApiKey
	templateID := fmt.Sprintf("%d", msg.TemplateID)
	params := url.Values{
		"smsUser":    {config.Conf.SendCloudSMSConf.ApiUser},
		"templateId": {templateID},
		"vars":       {string(vars)},
	}
	var signBody string
	if "+86" == msg.PhoneNumer.AreaCode() {
		signBody = fmt.Sprintf("%s&phone=%s&smsUser=%s&templateId=%d&vars=%s&%s",
			smsKey, msg.PhoneNumer.SubNumber(), config.Conf.SendCloudSMSConf.ApiUser, msg.TemplateID, string(vars), smsKey)
		params.Add("phone", string(msg.PhoneNumer.SubNumber()))
	} else {
		signBody = fmt.Sprintf("%s&msgType=2&phone=%s&smsUser=%s&templateId=%d&vars=%s&%s",
			smsKey, msg.PhoneNumer.FullNumber(), config.Conf.SendCloudSMSConf.ApiUser, msg.TemplateID, string(vars), smsKey)
		params.Add("msgType", "2")
		params.Add("phone", string(msg.PhoneNumer.FullNumber()))

	}
	h := md5.New()
	h.Write([]byte(signBody))
	cipherStr := h.Sum(nil)
	signature := hex.EncodeToString(cipherStr)
	params.Add("signature", signature)

	var body = &bytes.Buffer{}
	var writer = multipart.NewWriter(body)
	for key, value := range params {
		_ = writer.WriteField(key, value[0])
	}
	if err := writer.Close(); err != nil {
		return err
	}
	request, err := http.NewRequest("POST", RequestURI, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	responseHandler, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer responseHandler.Body.Close()
	bodyByte, err := ioutil.ReadAll(responseHandler.Body)
	loggers.Info.Printf("Send %s sms to %s by SendCloud result:%s", msg.TemplateName, msg.PhoneNumer.FullNumber(), string(bodyByte))

	return err
}

func SendSMSByAWS(sms *domain.SMSMessage, smsTemplate *domain.SMSTemplate) error {
	body, err := GenSMS(sms, smsTemplate)
	if err != nil {
		loggers.Info.Printf("SMSSender send %s sms to phone:%s error:%s ", sms.TemplateName, sms.PhoneNumer, err.Error())
		return err
	}
	if body == "" {
		loggers.Info.Printf("SMSSender send %s sms to phone:%s body is null ", sms.TemplateName, sms.PhoneNumer)
		return errors.New("sms body is null")
	}

	PhoneNumer := string(sms.PhoneNumer)
	msg := sns.PublishInput{
		Message:     &body,
		PhoneNumber: &PhoneNumer,
	}
	out, err := SNSClient.Publish(&msg)
	if err != nil {
		loggers.Info.Printf("SMSSender send %s sms to phone:%s error:%s ", sms.TemplateName, sms.PhoneNumer, err.Error())
		return err
	}
	loggers.Info.Printf("SMSSender send %s sms to %s by AWS id:%s OK", sms.TemplateName, sms.PhoneNumer, *out.MessageId)
	return nil
}

func GenSMS(sms *domain.SMSMessage, smsTemplate *domain.SMSTemplate) (string, error) {
	if sms.Language == "" {
		sms.Language = "zh-cn"
	}
	if smsTemplate.Body == "" {
		return "", errors.New("Template body is null")
	}

	t := template.New("").Funcs(template.FuncMap{
		"decimal": func(vol decimal.Decimal, places int32) string {
			return strings.TrimRight(vol.StringFixed(places), "0")
		},
	})

	tmpl, err := t.Parse(smsTemplate.Body)
	if err != nil {
		loggers.Error.Printf("SMSSender create template parser error:%s", err.Error())
		return "", err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, sms.Source)
	if err != nil {
		loggers.Error.Printf("SMSSender gen sms error:%s", err.Error())
		return "", err
	}

	return buffer.String(), nil
}

const (
	SecretKey = "7da8f21cebc95d161f3a822f1320b8f5"
	Secret    = "0db8649bc25bb11f6bfb8bc30293c35f"
)

//默认是验证码类短信
var BusinessID = "5209744e9d6d432eb4c9f9a5135a55db"

//根据secretKey和parameters生成签名
func genNECSignature(secretKey string, params url.Values) string {
	var keys []string
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buf := bytes.NewBufferString("")
	for _, key := range keys {
		buf.WriteString(key + params.Get(key))
	}
	buf.WriteString(secretKey)
	has := md5.Sum(buf.Bytes())
	return fmt.Sprintf("%x", has)
}

func SendSMSByNEC(msg *domain.SMSMessage) error {
	var (
		trans = &http.Transport{
			MaxIdleConns:       4,
			IdleConnTimeout:    time.Second * 30,
			DisableCompression: true,
		}
		client = http.Client{
			Transport: trans,
			Timeout:   time.Second * 10,
		}
		param = make(url.Values)
		form  = make(url.Values)
	)
	vars, ok := msg.Source.(map[string]interface{})
	if !ok {
		return errors.New("Invalid params")
	}
	for k, v := range vars {
		param.Add(k, fmt.Sprint(v))
	}
	form.Add("businessId", BusinessID)
	form.Add("mobile", msg.PhoneNumer.SubNumber())
	form.Add("secretId", Secret)
	form.Add("version", "v2")
	form.Add("templateId", fmt.Sprint(msg.TemplateID))
	form.Add("params", param.Encode()) //这里有个bug,param顺序必须给模版里的变量顺序一致否则405错误
	form.Add("needUp", "true")
	form.Add("timestamp", fmt.Sprint(time.Now().Unix()*1000))
	form.Add("nonce", fmt.Sprint(10000000+rand.Int31n(89999999)))
	form.Add("signature", genNECSignature(SecretKey, form))
	req, err := http.NewRequest(http.MethodPost, "https://sms.dun.163yun.com/v2/sendsms", bytes.NewReader([]byte(form.Encode())))
	if err != nil {
		loggers.Warn.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		loggers.Warn.Println(err)
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		loggers.Error.Printf("Request NECCaptch api error:%s", err.Error())
		return err
	}
	loggers.Info.Printf("Send %s sms to %s by NEC code:%d msg:%s data:%v", msg.TemplateName, msg.PhoneNumer.SubNumber(), result.Code, result.Msg, result.Data)

	return nil
}
