package message

import (
	"business/support/domain"
	"business/support/libraries/loggers"
	"encoding/json"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSClient struct {
	client *sqs.SQS
	conf   *domain.SQSConf
}

// 创建SQS客户端
func NewSQSClient(conf *domain.SQSConf) (*SQSClient, error) {
	if conf == nil {
		return nil, errors.New("Config is nil")
	}
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(conf.Region),
		Credentials: credentials.NewStaticCredentials(conf.AccessKey, conf.AccessSecret, ""),
	})
	if err != nil {
		return nil, err
	}
	return &SQSClient{sqs.New(sess), conf}, nil
}

// 获取AWS原生客户端
func (self *SQSClient) Client() *sqs.SQS {
	return self.client
}

// 获取SQS配置
func (self *SQSClient) Conf() *domain.SQSConf {
	return self.conf
}

// 从SQS中获取消息
func (self *SQSClient) Recevie(maxNum int64) ([]*domain.Message, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:            &self.conf.QueueUrl,
		MaxNumberOfMessages: &maxNum,
	}
	output, err := self.client.ReceiveMessage(&input)
	if err != nil {
		return nil, err
	}

	var records []*domain.Message
	for _, msg := range output.Messages {
		var record domain.Message
		if err := json.Unmarshal([]byte(*msg.Body), &record); err != nil {
			loggers.Warn.Println("Unmarshal record error:%s", err.Error())
			continue
		}
		record.MessageID = aws.StringValue(msg.MessageId)
		record.ReceiptHandle = aws.StringValue(msg.ReceiptHandle)
		records = append(records, &record)
	}

	return records, nil
}

// 发送消息至SQS(同步)
func (self *SQSClient) Send(msg *domain.Message) error {
	now := time.Now()
	msg.LogTime = &now
	v, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	body := string(v)
	input := sqs.SendMessageInput{
		MessageBody: &body,
		QueueUrl:    &self.conf.QueueUrl,
	}
	// loggers.Info.Printf("Message sqs send type:%s action:%s from_type:%s from:%s start",
	//	msg.Type, msg.Action, msg.FromType, msg.From)
	output, err := self.client.SendMessage(&input)
	if err != nil {
		loggers.Error.Printf("Message sqs send id:%s type:%s action:%s from_type:%s from:%s,input:%s",
			aws.StringValue(output.MessageId), msg.Type, msg.Action, msg.FromType, msg.From, input)
		return err
	}
	return nil
}

// 发送消息至SQS(异步)
func (self *SQSClient) SendAsync(msg *domain.Message) {
	go func(msg *domain.Message) {
		if err := self.Send(msg); err != nil {
			loggers.Info.Printf("Message sqs send type:%s action:%s from_type:%s from:%s error:%s",
				msg.Type, msg.Action, msg.FromType, msg.From, err.Error())
		}
	}(msg)
}

// 从SQS删除一条消息
func (self *SQSClient) Delete(record *domain.Message) {
	input := sqs.DeleteMessageInput{
		QueueUrl:      &self.conf.QueueUrl,
		ReceiptHandle: &record.ReceiptHandle,
	}

	if _, err := self.client.DeleteMessage(&input); err != nil {
		loggers.Warn.Printf("Delete message from sqs error:%s", err.Error())
	}
}

// 从SQS中删除消息(删除多条)
func (self *SQSClient) DeleteMulti(records []domain.Message) {
	for i := range records {
		self.Delete(&records[i])
	}
}
