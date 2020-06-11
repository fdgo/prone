package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"encoding/json"
	"fmt"
	"math/rand"
)

func GetSMSTemplate(name, language string) (*domain.SMSTemplate, error) {
	conn := redisPool.NewConn()
	var (
		err      error
		template domain.SMSTemplate
	)
	key := fmt.Sprintf("%s.%s_%s", domain.KSMSTemplate.Key, name, language)
	v, err := conn.Get(key).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &template); err != nil {
			return nil, err
		}
		return &template, nil
	}
	pTemplate, err := GetSMSTemplateFromDB(name, language)
	if err != nil {
		return nil, err

	}
	v, err = json.Marshal(pTemplate)
	if err != nil {
		return nil, err
	}
	if err := conn.Set(key, v, domain.KSMSTemplate.Timeout).Err(); err != nil {
		loggers.Warn.Printf("Cache %s sms template error:%s", name, err.Error())
	}
	return pTemplate, nil
}

func GetSMSTemplateFromDB(name, language string) (*domain.SMSTemplate, error) {
	db := confDBPool.NewConn()

	var template domain.SMSTemplate
	dbResult := db.Where("name = ?", name).Where("language = ?", language).First(&template)
	if dbResult.RecordNotFound() {
		loggers.Warn.Printf("SMS template:%s not exist", name)
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		loggers.Warn.Printf("Get SMS template:%s error:%s", name, dbResult.Error.Error())
		return nil, dbResult.Error
	}
	return &template, nil
}

func VerifyEmailCode(email domain.EMAIL, code, typ string) error {
	conn := redisPool.NewConn()
	key := fmt.Sprintf("%s.%s.EMAIL.%s", domain.KVerifyCode.Key, typ, email)
	v, err := conn.Get(key).Result()
	if err != nil {
		if redisPool.IsNil(err) {
			loggers.Debug.Printf("%s is nil", key)
			return errors.VerifyCodeError
		}
		return err
	}
	if v != code {
		limitKey := fmt.Sprintf("%s.LIMIT_ERR", key)
		if conn.Incr(limitKey).Val() >= 10 {
			conn.Del(key)
		}
		loggers.Debug.Printf("%s %s:%s", key, v, code)
		return errors.VerifyCodeError
	}
	conn.Del(key)
	return nil
}

func SendEmailVerifyCode(to domain.EMAIL, typ, language string) error {
	conn := redisPool.NewConn()

	retryKey := fmt.Sprintf("%s.%s", domain.KVerifyCodeTimes.Key, to)
	if err := conn.Get(retryKey).Err(); err == nil {
		loggers.Info.Printf("SendEmailVerifyCode %s request to offen", to)
		return nil
	} else {
		conn.Set(retryKey, "1", domain.KVerifyCodeTimes.Timeout)
	}

	code := rand.Intn(899999) + 100000
	key := fmt.Sprintf("%s.%s.EMAIL.%s", domain.KVerifyCode.Key, typ, to)
	if err := conn.Set(key, code, domain.KVerifyCode.Timeout).Err(); err != nil {
		conn.Del(retryKey)
		return err
	}

	source := struct{ Code int }{code}
	emailMessage := domain.EmailMessage{
		Sender:       "no-reply@bifund.com",
		Recipient:    to,
		TemplateName: typ,
		Language:     language,
		Source:       &source,
	}
	SendEmail(&emailMessage)

	return nil
}

func SendSMSVerifyCode(to domain.PHONE, typ, language string) error {
	conn := redisPool.NewConn()
	retryKey := fmt.Sprintf("%s.%s", domain.KVerifyCodeTimes.Key, to)
	if err := conn.Get(retryKey).Err(); err == nil {
		loggers.Info.Printf("SendSMSVerifyCode %s request to offen", to)
		return nil
	} else {
		conn.Set(retryKey, "1", domain.KVerifyCodeTimes.Timeout)
	}
	code := rand.Intn(899999) + 100000
	key := fmt.Sprintf("%s.%s.PHONE.%s", domain.KVerifyCode.Key, typ, to)
	if err := conn.Set(key, code, domain.KVerifyCode.Timeout).Err(); err != nil {
		conn.Del(retryKey)
		return err
	}

	source := map[string]interface{}{
		"code": code,
	}
	msg := domain.SMSMessage{
		PhoneNumer:   to,
		Source:       &source,
		TemplateName: typ,
		Language:     language,
	}
	SendSMS(&msg)

	loggers.Debug.Printf("SendSMSVerifyCode type:%s phone:%s code:%d", typ, to, code)
	return nil
}

func VerifySMSCode(phone domain.PHONE, code, typ string) error {
	conn := redisPool.NewConn()
	key := fmt.Sprintf("%s.%s.PHONE.%s", domain.KVerifyCode.Key, typ, phone)
	v, err := conn.Get(key).Result()
	if err != nil {
		if redisPool.IsNil(err) {
			loggers.Debug.Printf("%s is nil", key)
			return errors.VerifyCodeError
		}
		return err
	}
	if v != code {
		limitKey := fmt.Sprintf("%s.LIMIT_ERR", key)
		if conn.Incr(limitKey).Val() >= 10 {
			conn.Del(key)
		}
		loggers.Debug.Printf("%s %s:%s", key, v, code)
		return errors.VerifyCodeError
	}
	conn.Del(key)
	return nil

}

func SendEmail(email *domain.EmailMessage) {
	message := domain.Message{
		MessageBase: domain.MessageBase{
			Action:   domain.MESSAGE_ACTION_CREATED,
			Type:     domain.MESSAGE_TYPE_EMAIL,
			From:     "IFACCOUNT",
			FromType: domain.FROM_TYPE_SYSTEM,
		},
		EmailMessage: email,
	}

	sqsClient.SendAsync(&message)
}

func SendSMS(sms *domain.SMSMessage) {
	message := domain.Message{
		MessageBase: domain.MessageBase{
			Action:   domain.MESSAGE_ACTION_CREATED,
			Type:     domain.MESSAGE_TYPE_SMS,
			From:     "IFACCOUNT",
			FromType: domain.FROM_TYPE_SYSTEM,
		},
		SMSMessage: sms,
	}

	sqsClient.SendAsync(&message)
}
