package actionmsg

import (
	"business/support/config"
	"business/support/libraries/msg"
)

var broadcaster *msg.Broadcaster

func InitActionBroadcaster() error {
	broadcaster = &msg.Broadcaster{}
	return broadcaster.Init(config.Conf.SNSConf.Region,
		config.Conf.SNSConf.AccessKey,
		config.Conf.SNSConf.AccessSecret,
		"",
		config.Conf.SNSConf.ActionTopicArn)
}
