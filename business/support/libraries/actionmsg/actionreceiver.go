package actionmsg

import (
	"business/support/config"
	"business/support/libraries/msg"
)

var receiver *msg.Receiver

func InitActionReceiver() error {
	receiver = &msg.Receiver{}
	return receiver.Init(config.Conf.SQSConf.Region,
		config.Conf.SQSConf.AccessKey,
		config.Conf.SQSConf.AccessSecret,
		"",
		config.Conf.SQSConf.ActivityQueueUrl)
}
