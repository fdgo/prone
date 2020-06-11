package msg

import (
	"business/support/libraries/loggers"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQSSender struct {
	client   *sqs.SQS
	queueUrl string
}

func (this *SQSSender) send(msg *Msg) error {
	var delaySeconds = int64(0)
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	msgBody := string(body)
	input := sqs.SendMessageInput{
		DelaySeconds: &delaySeconds,
		MessageBody:  &msgBody,
		QueueUrl:     &(this.queueUrl),
	}
	out, err := this.client.SendMessage(&input)
	if err != nil {
		return err
	}
	loggers.Info.Printf("SQSSender send msg queue id:%s md5:%s", *out.MessageId, *out.MD5OfMessageBody)
	return nil
}
