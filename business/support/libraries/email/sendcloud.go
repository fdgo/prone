package email

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
)

func SendBySendCloud(email *domain.Email, files *map[string][]byte) error {
	RequestURI := config.Conf.SendCloudEmailConf.ApiUrl
	PostParams := url.Values{
		"apiUser":  {config.Conf.SendCloudEmailConf.ApiUser},
		"apiKey":   {config.Conf.SendCloudEmailConf.ApiKey},
		"from":     {email.Sender},
		"fromName": {"no-reply@bifund.com"},
		"to":       {string(email.Recipient)},
		"subject":  {email.Subject},
		"html":     {email.HtmlBody},
	}
	success, respBody, err := sendBySendCloudWithAttachment(RequestURI, PostParams, files)
	if err != nil {
		loggers.Info.Printf("Send email with SendCound subject:%s to:%s error:%s", email.Subject, email.Recipient, err.Error())
	}
	if success {
		loggers.Info.Printf("Send email with SendCound subject:%s to:%s OK %s ", email.Subject, email.Recipient, respBody)
	} else {
		loggers.Info.Printf("Send email with SendCound subject:%s to:%s FAIL %s", email.Subject, email.Recipient, respBody)
	}
	return nil
}

//使用sendcloud发送带附件邮件
func sendBySendCloudWithAttachment(url string, params url.Values, files *map[string][]byte) (bool, string, error) {
	var body = &bytes.Buffer{}
	var writer = multipart.NewWriter(body)

	if files != nil {
		for k, v := range *files {
			fileWriter, err := writer.CreateFormFile("attachments", k)
			if err != nil {
				return false, "", err
			}
			fileWriter.Write(v)
		}
	}

	for key, value := range params {
		_ = writer.WriteField(key, value[0])
	}

	var err = writer.Close()
	if err != nil {
		return false, "", err
	}

	request, err := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	responseHandler, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, "", err
	}
	defer responseHandler.Body.Close()

	bodyByte, err := ioutil.ReadAll(responseHandler.Body)
	if err != nil {
		return false, string(bodyByte), err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyByte, &result)
	return (result["result"] == true), string(bodyByte), err

}
