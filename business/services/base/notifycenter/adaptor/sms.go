package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	smssender "business/support/libraries/sms"
	"encoding/json"
	"fmt"
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
	db := confDB.NewConn()
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

func SendSMS(sms *domain.SMSMessage) error {
	if sms.PhoneNumer == "" {
		loggers.Error.Printf("SendSMS phone number is null")
		return nil
	}
	smsTemplate, err := GetSMSTemplate(sms.TemplateName, sms.Language)
	if err != nil {
		return err
	}
	areaService := GetSMSService(sms.PhoneNumer.AreaCode())
	if len(smsTemplate.Service) > 0 {
		if smsTemplate.Service == domain.SMS_SERVICE_NEC && sms.PhoneNumer.AreaCode() != "+86" {
			sms.Service = areaService
		} else {
			sms.Service = smsTemplate.Service
		}
	} else {
		sms.Service = areaService
	}
	return smssender.SendSMS(sms, smsTemplate)
}

func GetSMSService(areaCode string) string {
	services, err := GetSMSServices()
	if err != nil {
		return ""
	}
	if service, ok := (*services)[areaCode]; ok {
		return service
	}
	if service, ok := (*services)[domain.SMS_SERVICE_DEFAULT]; ok {
		return service
	}
	return ""
}

func GetSMSServices() (*map[string]string, error) {
	conn := redisPool.NewConn()
	var services map[string]string
	v, err := conn.Get(domain.KSMSServies.Key).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &services); err == nil {
			return &services, nil
		}
		loggers.Error.Printf("Unmarshal sms services from cache error %s", err.Error())
	}
	pServices, err := GetSMSServicesFromDB()
	if err != nil {
		loggers.Error.Printf("Get sms servies config from db error %s", err.Error())
		return nil, err
	}
	value, err := json.Marshal(*pServices)
	if err != nil {
		loggers.Error.Printf("Marshal sms servies config error %s", err.Error())
	}
	if err := conn.Set(domain.KSMSServies.Key, value, domain.KEmailServies.Timeout).Err(); err != nil {
		loggers.Error.Printf("Cache email servies config error %s", err.Error())
	}
	return pServices, nil
}

func GetSMSServicesFromDB() (*map[string]string, error) {
	db := confDB.NewConn()
	var services []domain.SMSService
	dbResult := db.Find(&services)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	if len(services) == 0 {
		return nil, errors.NotFound
	}
	servicesMap := make(map[string]string)
	for _, ins := range services {
		servicesMap[ins.AreaCode] = ins.Service
	}
	return &servicesMap, nil
}
