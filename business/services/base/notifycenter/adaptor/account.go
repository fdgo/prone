package adaptor

import (
	"business/support/domain"
	"time"
)

func GetAccountByID(accountID int64) (*domain.Account, error) {
	var account domain.Account
	if err := exDB.NewConn().Where("account_id = ?", accountID).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func GetLatestRecord(accountID int64, logTime *time.Time) (*domain.LoginRecord, error) {
	var accountRecord domain.AccountRecord

	if err := recordDB.NewConn().
		Where("account_id = ?", accountID).
		Where("base_action = ?", "LOGIN").
		Where("log_time < ?", logTime).
		Order("log_time DESC").
		First(&accountRecord).Error; err != nil {
		return nil, err
	}
	return &domain.LoginRecord{
		LoginTime: accountRecord.LogTime,
		LoginIP:   accountRecord.ClientIP,
		Device:    accountRecord.Device,
	}, nil
}

func GetLoginRecord(accountId, limit int64) ([]domain.LoginRecord, error) {
	var (
		accountRecords []domain.AccountRecord
		loginRecords   []domain.LoginRecord
	)
	if err := recordDB.NewConn().
		Where("account_id = ?", accountId).
		Where("base_action = ?", "LOGIN").
		Order("log_time DESC").
		Limit(limit).
		Find(&accountRecords).Error; err != nil {
		return nil, err
	}
	for i := range accountRecords {
		loginRecords = append(loginRecords, domain.LoginRecord{
			LoginTime: accountRecords[i].LogTime,
			LoginIP:   accountRecords[i].ClientIP,
			Device:    accountRecords[i].Device,
		})
	}
	return loginRecords, nil
}
