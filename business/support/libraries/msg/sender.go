package msg

import (
	"encoding/json"
	"errors"
)

type Sender struct {
	sender *SQSSender
}

func (this *Sender) Init(
	region string,
	accessKey string,
	accessSecret string,
	token string,
	queueUrl string) error {
	client := newSqsClient(region, accessKey, accessSecret, token)
	if nil == client {
		return errors.New("new sqs sender failed!")
	}
	this.sender = &SQSSender{client: client, queueUrl: queueUrl}
	return nil
}

func (this *Sender) Send(msgId int32, v interface{}) error {
	if nil == this.sender {
		return errors.New("no sender")
	}
	body, err := json.Marshal(v)
	if nil != err {
		return err
	}
	return this.sender.send(&Msg{Id: msgId, Body: body})
}
