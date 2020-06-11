package msg

import (
	"business/support/libraries/loggers"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSReceiver struct {
	client   *sqs.SQS
	queueUrl string
}

func (this *SQSReceiver) receive(waitTimeSeconds int64) ([]*Msg, error) {
	input := sqs.ReceiveMessageInput{
		QueueUrl:        &(this.queueUrl),
		WaitTimeSeconds: &waitTimeSeconds,
	}
	out, err := this.client.ReceiveMessage(&input)
	if err != nil {
		time.Sleep(time.Second * time.Duration(waitTimeSeconds))
		loggers.Debug.Printf("SQSReceiver::receive error:%s", err.Error())
		return nil, err
	}
	size := len(out.Messages)
	if size <= 0 {
		return nil, nil
	}
	msgs := make([]*Msg, 0, size)
	for _, item := range out.Messages {
		if nil == item.ReceiptHandle {
			continue
		}
		if item.Body == nil {
			this.delHandler(item.ReceiptHandle)
			continue
		}
		var msg Msg
		if err := json.Unmarshal([]byte(*(item.Body)), &msg); err != nil {
			this.delHandler(item.ReceiptHandle)
			loggers.Debug.Printf("SQSReceiver::receive unmarshal error:%s", err.Error())
			continue
		}
		if len(msg.Message) > 0 {
			if err := json.Unmarshal([]byte(msg.Message), &msg); err != nil {
				this.delHandler(item.ReceiptHandle)
				loggers.Debug.Printf("SQSReceiver::receive unmarshal message error:%s", err.Error())
				continue
			}
		}
		msg.Handler = *(item.ReceiptHandle)
		msgs = append(msgs, &msg)
	}
	return msgs, nil
}

func (this *SQSReceiver) delete(msgs []*Msg) {
	input := sqs.DeleteMessageInput{
		QueueUrl: &(this.queueUrl),
	}
	for _, item := range msgs {
		input.ReceiptHandle = &(item.Handler)
		this.client.DeleteMessage(&input)
	}
}

func (this *SQSReceiver) deleteOne(msg *Msg) {
	input := sqs.DeleteMessageInput{
		QueueUrl: &(this.queueUrl),
	}
	input.ReceiptHandle = &(msg.Handler)
	this.client.DeleteMessage(&input)

}

func (this *SQSReceiver) delHandler(handler *string) {
	input := sqs.DeleteMessageInput{
		QueueUrl: &(this.queueUrl),
	}
	input.ReceiptHandle = handler
	this.client.DeleteMessage(&input)

}
