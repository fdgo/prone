package message

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"testing"
)

func Test_Publish(t *testing.T) {
	sns, err := NewSNSClient(&config.Conf.ActionSNSConf)
	if err != nil {
		t.Error(err)
	}
	order := domain.Order{
		OrderId:   10000001,
		StockCode: "ETH/BTC",
	}
	record := domain.Record{
		Order: &order,
	}
	record.Action = "PUBLISH_CREATED"
	record.RecordType = domain.RECORD_TYPE_ASSET
	record.UserType = domain.USER_TYPE_NORMAL
	record.User = "guyun.hy@gmail.com"
	if err = sns.Publish(&record); err != nil {
		t.Error(err)
	}
}

func Test_Send(t *testing.T) {
	sqs, err := NewSQSClient(&config.Conf.RecordSQSConf)
	if err != nil {
		t.Error(err)
	}
	order := domain.Order{
		OrderId:   10000001,
		StockCode: "ETH/BTC",
	}
	record := domain.Record{
		Order: &order,
	}
	record.Action = "SEND_CREATED"
	record.RecordType = domain.RECORD_TYPE_ASSET
	record.UserType = domain.USER_TYPE_NORMAL
	record.User = "guyun.hy@gmail.com"
	if err = sqs.Send(&record); err != nil {
		t.Error(err)
	}
}

func Test_Recevie(t *testing.T) {
	sqs, err := NewSQSClient(&config.Conf.RecordSQSConf)
	if err != nil {
		t.Error(err)
	}
	records, err := sqs.Recevie(3)
	loggers.Debug.Println(records)
	if err != nil {
		t.Error(err)
	}
	sqs.Delete(records)
}
