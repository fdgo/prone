package msg

import "errors"

type Receiver struct {
	receiver *SQSReceiver
}

func (this *Receiver) Init(
	region string,
	accessKey string,
	accessSecret string,
	token string,
	queueUrl string) error {
	client := newSqsClient(region, accessKey, accessSecret, token)
	if nil == client {
		return errors.New("new sqs receiver failed!")
	}
	this.receiver = &SQSReceiver{client: client, queueUrl: queueUrl}
	return nil
}

func (this *Receiver) Receive(waitTimeSeconds int64) ([]*Msg, error) {
	if nil == this.receiver {
		return nil, nil
	}
	return this.receiver.receive(waitTimeSeconds)
}

func (this *Receiver) DeleteMsgs(msgs []*Msg) {
	if nil == this.receiver {
		return
	}
	this.receiver.delete(msgs)
}

func (this *Receiver) DeleteMsg(msg *Msg) {
	if nil == this.receiver {
		return
	}
	this.receiver.deleteOne(msg)
}
