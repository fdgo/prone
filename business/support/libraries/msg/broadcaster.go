package msg

import (
	"business/support/libraries/loggers"
	"encoding/json"
	"errors"
)

type Broadcaster struct {
	sender *SNSSender
}

func (this *Broadcaster) Init(
	region string,
	accessKey string,
	accessSecret string,
	token string,
	topicArn string) error {
	loggers.Debug.Printf(`Broadcaster::init region:%s,
		accessKey:%s,
		accessSecret:%s,
		token:%s,
		topicArn:%s`, region, accessKey, accessSecret, token, topicArn)
	client := newSnsClient(region, accessKey, accessSecret, token)
	if nil == client {
		return errors.New("new sns broadcaster failed!")
	}
	this.sender = &SNSSender{client: client, topicArn: topicArn}
	return nil
}

func (this *Broadcaster) Send(msgId int32, v interface{}) error {
	if nil == this.sender {
		return errors.New("no broadcaster")
	}
	body, err := json.Marshal(v)
	if nil != err {
		return err
	}
	return this.sender.send(&Msg{Id: msgId, Body: body})
}
