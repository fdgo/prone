package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	cfg "business/support/model/config"
	"business/support/rpc"
	"context"
	"encoding/base32"
	"fmt"
	"math/rand"
	"time"
)

func GetAccountByUid(uid int64) (*domain.Account, error) {
	if uid == 0 {
		return nil, errors.InternalServerError
	}
	db := dbPool.NewConn()
	var account domain.Account
	dbResult := db.Where("account_id = ?", uid).First(&account)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &account, nil
}

func GetAccountByEmail(email domain.EMAIL) (*domain.Account, error) {
	db := dbPool.NewConn()
	var account domain.Account
	dbResult := db.Where("email = ?", email).First(&account)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &account, nil
}

func GetLoginRecord(accountId, limit int64) ([]domain.LoginRecord, error) {
	var (
		accountRecords []domain.AccountRecord
		loginRecords   []domain.LoginRecord
	)
	if err := dbRecordsPool.NewConn().
		Where("account_id = ?", accountId).
		Where("base_action = ?", "LOGIN").
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

func GetAccountByPhone(phone domain.PHONE) (*domain.Account, error) {
	db := dbPool.NewConn()
	var account domain.Account
	dbResult := db.Where("phone = ?", phone).First(&account)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return &account, nil
}

func GetUserAssets(uid int64) (*[]domain.UserAssets, error) {
	db := dbPool.NewConn()
	var assets []domain.UserAssets

	dbResult := db.Where("account_id = ?", uid).Find(&assets)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	return &assets, nil
}

func CreateAccount(account *domain.Account, inviterId int64) (*domain.Account, error) {
	var reply rpc.CreateAccountReply
	req := domain.Account{
		Email:         account.Email,
		Phone:         account.Phone,
		Password:      account.Password,
		AccountName:   account.AccountName,
		AccountStatus: account.AccountStatus,
		AccountType:   account.AccountType,
		RegisterIp:    account.RegisterIp,
		QD:            account.QD,
		LatestLoginAt: account.LatestLoginAt,
		LatestLoginIp: account.RegisterIp,
	}
	if err := accountRPC.Create(context.Background(), &req, &reply); err != nil {
		return nil, errors.InternalServerError
	}
	reply.Account.Password = ""
	if nil != reply.Error {
		return nil, reply.Error
	} else if nil == reply.Account {
		return nil, errors.SysError
	}
	loggers.Debug.Printf("CreateAccount succeed")
	return reply.Account, nil
}

func ActiveAccount(account *domain.Account) error {
	var reply rpc.ActiveAccountReply
	req := domain.Account{
		AccountId: account.AccountId,
	}
	if err := accountRPC.ActiveAccount(context.Background(), &req, &reply); err != nil {
		return errors.InternalServerError
	}
	if nil != reply.Error {
		return reply.Error
	}
	return nil
}

func ResetPassword(account *domain.Account) error {
	var reply rpc.ResetPasswordReply
	req := domain.Account{
		AccountId: account.AccountId,
		Password:  account.Password,
	}
	if err := accountRPC.ResetPassword(context.Background(), &req, &reply); err != nil {
		return errors.InternalServerError
	}
	if nil != reply.Error {
		return reply.Error
	}
	return nil
}

func BindEmail(account *domain.Account) error {
	var reply rpc.BindEmailReply
	if err := accountRPC.BindEmail(context.Background(), account, &reply); err != nil {
		return errors.InternalServerError
	}
	if nil != reply.Error {
		loggers.Debug.Println(reply.Error)
		return reply.Error
	}
	return nil
}

func BindPhone(account *domain.Account) error {
	var reply rpc.BindPhoneReply
	if err := accountRPC.BindPhone(context.Background(), account, &reply); err != nil {
		return errors.InternalServerError
	}
	if nil != reply.Error {
		return reply.Error
	}
	return nil
}

func PublishAccountRecord(record *domain.Message) {
	snsClient.PublishAsync(record)
}

func AddAssetPassword(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.AddAssetPassword(context.Background(), account, &reply); err != nil {
		if errors.Equal(err, errors.AssetPasswordHadExisted) {
			return errors.AssetPasswordHadExisted
		}
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func ResetAssetPassword(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.ResetAssetPassword(context.Background(), account, &reply); err != nil {
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func ResetAssetPasswordEffectiveTime(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.ResetAssetPasswordEffectiveTime(context.Background(), account, &reply); err != nil {
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func SetAntiFishingText(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.SetAntiFishingText(context.Background(), account, &reply); err != nil {
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func SetGAKey(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.SetGAKey(context.Background(), account, &reply); err != nil {
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func GenGAKey(account *domain.Account) (string, error) {
	key := fmt.Sprintf("%s.%s", domain.KGAKey.Key, account)
	source := []rune("234567ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for {
		rand.Seed(time.Now().UnixNano())
		gaKey := make([]rune, 16)
		for i := range gaKey {
			gaKey[i] = source[rand.Intn(16)]
		}
		if err := redisPool.NewConn().Set(key, string(gaKey), domain.KGAKey.Timeout).Err(); err != nil {
			return "", err
		}
		if _, err := base32.StdEncoding.DecodeString(string(gaKey)); err != nil {
			continue
		}

		return string(gaKey), nil
	}
}

func GetCacheGAKey(account *domain.Account) (string, error) {
	key := fmt.Sprintf("%s.%s", domain.KGAKey.Key, account)
	gaKey, err := redisPool.NewConn().Get(key).Result()
	if err != nil {
		return "", err
	}
	return gaKey, err
}

func KYCUpdate(account *domain.Account) error {
	var reply domain.Account
	if err := accountRPC.KYCUpdate(context.Background(), account, &reply); err != nil {
		if errors.InvalidKYCStatus.Equal(err) {
			return errors.BadRequest
		}
		if errors.Equal(err, errors.AccountNotFound) {
			return errors.AccountNotFound
		}
		return errors.InternalServerError
	}
	return nil
}

func SetAccountName(account *domain.Account) (*domain.Account, error) {
	var reply domain.Account
	if err := accountRPC.SetAccountName(context.Background(), account, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func UpdateAccountAvatar(account *domain.Account) (*domain.Account, error) {
	var reply domain.Account
	if err := accountRPC.UpdateAccountAvatar(context.Background(), account, &reply); err != nil {
		return nil, err
	}
	return &reply, nil
}

func CheckKYCIDNo(IDNo string) (bool, error) {
	account := new(domain.Account)
	dbResult := dbPool.NewConn().Where("passed_id_no = ?", IDNo).Where("kyc_status = ?", domain.KYC_STATUS_APPROVED).First(account)
	if dbResult.RecordNotFound() {
		return true, nil
	}
	if err := dbResult.Error; err != nil {
		return false, err
	}
	return false, nil

}

func GetAccountByRegisterIp(ip string, registerTime *time.Time) int {
	var count int
	dbPool.NewConn().Model(&domain.Account{}).Where("register_ip = ?", ip).Where("created_at > ?", registerTime).Count(&count)
	return count
}

type RegisterIPLimitConfig struct {
	LimitTime int `json:"limit_time"`
	Count     int `json:"count"`
}

func GetRegisterIpLimitConfig() *RegisterIPLimitConfig {
	var (
		key  = "RegisterIPLimit"
		conf RegisterIPLimitConfig
	)
	if err := cfg.GetAndCacheGlobal(key).UnmarshalJSON(&conf); err != nil {
		return nil
	}
	return &conf
}
