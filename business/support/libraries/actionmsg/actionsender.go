package actionmsg

import (
	"business/support/config"
	"business/support/libraries/msg"
)

var sender *msg.Sender

func InitActionSender() error {
	sender = &msg.Sender{}
	return sender.Init(config.Conf.SQSConf.Region,
		config.Conf.SQSConf.AccessKey,
		config.Conf.SQSConf.AccessSecret,
		"",
		config.Conf.SQSConf.ActivityQueueUrl)
}
