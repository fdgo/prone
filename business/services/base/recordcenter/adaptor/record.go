package adaptor

import (
	"business/support/domain"
	"business/support/libraries/loggers"
)

func InsertOrderRecord(records []*domain.OrderRecord) error {
	trans := dbPool.NewConn().Begin()
	for _, ins := range records {
		if ins == nil {
			continue
		}
		if err := dbPool.NewConn().Create(ins).Error; err != nil {
			loggers.Error.Printf("Insert order record error:%s", err.Error())
			trans.Rollback()
			return err
		}
	}
	trans.Commit()
	return nil
}

func InsertSettleRecord(records []*domain.SettleRecord) error {
	trans := dbPool.NewConn().Begin()
	for _, ins := range records {
		if ins == nil {
			continue
		}
		ins.ID = 0
		if err := dbPool.NewConn().Create(ins).Error; err != nil {
			loggers.Error.Printf("Insert settle record error:%s", err.Error())
			trans.Rollback()
			return err
		}
	}
	trans.Commit()
	return nil
}

func InsertAssetRecord(records []*domain.AssetRecord) error {
	trans := dbPool.NewConn().Begin()
	for _, ins := range records {
		if ins == nil {
			continue
		}
		ins.ID = 0
		if err := dbPool.NewConn().Create(ins).Error; err != nil {
			loggers.Error.Printf("Insert asset record error:%s", err.Error())
			trans.Rollback()
			return err
		}
	}
	trans.Commit()
	return nil
}

func InsertAccountRecord(records []*domain.AccountRecord) error {
	trans := dbPool.NewConn().Begin()
	for _, ins := range records {
		if ins == nil {
			continue
		}
		ins.ID = 0
		if err := dbPool.NewConn().Create(ins).Error; err != nil {
			loggers.Error.Printf("Insert account record error:%s", err.Error())
			trans.Rollback()
			return err
		}
	}
	trans.Commit()
	return nil
}
