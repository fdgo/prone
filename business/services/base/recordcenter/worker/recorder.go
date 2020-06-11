package worker

import (
	"business/support/config"
	"business/support/domain"

	"business/services/base/recordcenter/adaptor"
	"business/support/libraries/loggers"
	"business/support/libraries/message"
	"time"
)

type Recorder struct {
	Interval  time.Duration
	sqsClient *message.SQSClient
}

func (self *Recorder) Start() {
	var err error
	self.sqsClient, err = message.NewSQSClient(config.Conf.RecordSQSConf)
	if err != nil {
		panic(err)
	}
	for {
		self.Run()
	}
}

func (self *Recorder) Run() {
	records, err := self.sqsClient.Recevie(10)
	if err != nil {
		loggers.Error.Printf("AssetRecorder recevie record error:%s", err.Error())
		time.Sleep(self.Interval)
	}

	for i := range records {
		r := records[i]
		loggers.Info.Printf("Handle msg id:%s type:%s action:%s from_type:%s from:%s", r.MessageID, r.Type, r.Action, r.FromType, r.From)
		switch r.Type {
		case domain.MESSAGE_TYPE_WITHDRAW:
			{
				if err := handleSettleRecord(r); err != nil {
					loggers.Error.Printf("handle record error:%s", err.Error())
				} else {
					self.sqsClient.Delete(r)
				}
			}
		case domain.MESSAGE_TYPE_DEPOSIT:
			{
				if err := handleSettleRecord(r); err != nil {
					loggers.Error.Printf("handle record error:%s", err.Error())
				} else {
					self.sqsClient.Delete(r)
				}
			}

		case domain.MESSAGE_TYPE_ACCOUNT:
			{
				if err := handleAccountRecord(r); err != nil {
					loggers.Error.Printf("handle record error:%s", err.Error())
				} else {
					self.sqsClient.Delete(r)
				}
			}
		default:
			{
				loggers.Warn.Printf("RecordCenter not follow this message type:%s", r.Type)
				self.sqsClient.Delete(r)
			}
		}
	}
}

func handleSettleRecord(r *domain.Message) error {
	var (
		record = domain.SettleRecord{
			MessageBase: r.MessageBase,
			Settle:      *r.Settle,
		}
		assets []*domain.AssetRecord
	)
	if err := adaptor.InsertSettleRecord([]*domain.SettleRecord{&record}); err != nil {
		return err
	}

	for i := range r.Assets {
		r.Assets[i].MessageBase = r.MessageBase
		assets = append(assets, &r.Assets[i])
	}

	if err := adaptor.InsertAssetRecord(assets); err != nil {
		return err
	}

	return nil
}

func handleAccountRecord(r *domain.Message) error {
	r.Account.MessageBase = r.MessageBase
	if err := adaptor.InsertAccountRecord([]*domain.AccountRecord{r.Account}); err != nil {
		return err
	}

	return nil
}
