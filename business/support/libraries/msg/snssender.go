package msg

import (
	"business/support/libraries/loggers"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/sns"
)

type SNSSender struct {
	client   *sns.SNS
	topicArn string
}

func (this *SNSSender) send(msg *Msg) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// loggers.Info.Printf("SNSSender send msg topic arn:%s", this.topicArn)
	msgBody := string(body)
	input := sns.PublishInput{
		Message:  &msgBody,
		TopicArn: &(this.topicArn),
	}
	out, err := this.client.Publish(&input)
	if err != nil {
		loggers.Info.Printf("SNSSender send msg failed out:%#v,err:%s", out, err.Error())
		return err
	}
	loggers.Info.Printf("SQSSender send msg queue id:%s", *out.MessageId)
	return nil
}
