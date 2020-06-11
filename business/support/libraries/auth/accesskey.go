package auth

import (
	"business/support/domain"
)

//func AccessKeyAuth(accessKey, signatureStr, remoteIp string, body []byte, ts int64) (*domain.Session, error) {
//	if !CheckTs(ts) {
//		return nil, errors.ClientTimeInvalidError
//	}
//	apiKey, err := GetApiKey(accessKey)
//	if err != nil {
//		return nil, err
//	}
//	if apiKey.BindIP != nil && len(*apiKey.BindIP) != 0 {
//		pass := false
//		for _, ip := range *apiKey.BindIP {
//			if ip == remoteIp {
//				pass = true
//				break
//			}
//		}
//		if !pass {
//			return nil, errors.New("client ip no bind ip")
//		}
//	}
//	if apiKey.ExpireAt != nil {
//		if time.Now().UTC().After(*apiKey.ExpireAt) {
//			return nil, errors.New("accessKey expired")
//		}
//	}
//
//	md5Ctx := md5.New()
//	var ctx bytes.Buffer
//	ctx.Write(body)
//	ctx.WriteString(apiKey.AccessSecret)
//	ctx.WriteString(fmt.Sprintf("%d", ts))
//	md5Ctx.Write(ctx.Bytes())
//	cipher := md5Ctx.Sum(nil)
//	if strings.ToUpper(signatureStr) != strings.ToUpper(hex.EncodeToString(cipher)) {
//		loggers.Debug.Printf("body:%s", ctx.String())
//		loggers.Debug.Printf("key:%s", accessKey)
//		loggers.Debug.Printf("sec:%s", apiKey.AccessSecret)
//		loggers.Debug.Printf("com:%s", strings.ToUpper(hex.EncodeToString(cipher)))
//		loggers.Debug.Printf("rec:%s", strings.ToUpper(signatureStr))
//		return nil, errors.New("no match")
//	}
//	account, err := GetAccountByUid(apiKey.AccountId)
//	if err != nil {
//		return nil, errors.AccountNotFound
//	}
//
//	if account.AccountStatus == domain.ACCOUNT_STATUS_DISABLE {
//		return nil, errors.AccountHadDisabled
//	}
//
//	session := domain.Session{
//		Account: account,
//	}
//	return &session, nil
//}

func GetAccountByUid(accountId int64) (*domain.Account, error) {
	var account domain.Account
	if err := dbPool.NewConn().Where("account_id = ?", accountId).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
