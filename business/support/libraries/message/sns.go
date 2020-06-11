package message

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSClient struct {
	client *sns.SNS
	conf   *domain.SNSConf
}

// 创建SNS Client
func NewSNSClient(conf *domain.SNSConf) (*SNSClient, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.AccessSecret, ""),
	})
	if err != nil {
		return nil, err
	}
	return &SNSClient{sns.New(sess), conf}, nil
}

// 获取AWS原始Client
func (self *SNSClient) Client() *sns.SNS {
	return self.client
}

// 获取SNS配置
func (self *SNSClient) Conf() *domain.SNSConf {
	return self.conf
}

// 发布消息到SNS(同步)
func (self *SNSClient) Publish(msg *domain.Message) error {
	now := time.Now()
	msg.LogTime = &now
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	body := string(v)
	input := sns.PublishInput{
		Message:  &body,
		TopicArn: &self.conf.TopicArn,
	}

	// 订阅者可以通过此这些参数来筛选来自SNS的消息,在AWS SNS控制面板里面可以设置
	attrs := map[string]*sns.MessageAttributeValue{
		"type":      NewSNSMsgAttrStrVal(msg.Type),
		"action":    NewSNSMsgAttrStrVal(msg.Action),
		"from":      NewSNSMsgAttrStrVal(msg.From),
		"from_type": NewSNSMsgAttrStrVal(string(msg.FromType)),
	}
	input.SetMessageAttributes(attrs)

	loggers.Debug.Printf("Message sns publish type:%s action:%s from_type:%s from:%s start",
		msg.Type, msg.Action, msg.FromType, msg.From)
	output, err := self.client.Publish(&input)
	if err != nil {
		return err
	}
	loggers.Debug.Printf("Message sns publish id:%s type:%s action:%s from_type:%s from:%s start",
		aws.StringValue(output.MessageId), msg.Type, msg.Action, msg.FromType, msg.From)
	return nil
}

// 发布消息到SNS(异步)
func (self *SNSClient) PublishAsync(msg *domain.Message) {
	go func(msg *domain.Message) {
		if err := self.Publish(msg); err != nil {
			loggers.Error.Printf("Message sns publish msg:%s, error:%s",
				msg, err.Error())
		}
	}(msg)
}

// 创建SNS消息string类型属性
func NewSNSMsgAttrStrVal(val string) *sns.MessageAttributeValue {
	attr := new(sns.MessageAttributeValue)
	return attr.SetStringValue(val).SetDataType("String")
}
