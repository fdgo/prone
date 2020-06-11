package adaptor

import (
	"business/support/domain"
	emailSender "business/support/libraries/email"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/shopspring/decimal"
)

func NextSendEmail(timeout time.Duration) (*domain.Email, error) {
	conn := redisPool.NewConn()
	key := domain.KWaitSendEmails
	v, err := conn.BLPop(timeout, key.Key).Result()

	if err != nil {
		return nil, err
	}
	var email domain.Email
	if err := json.Unmarshal([]byte(v[1]), &email); err != nil {
		return nil, err
	}

	return &email, err
}

func GetEmailTemplate(name, language string) (*domain.EmailTemplate, error) {
	conn := redisPool.NewConn()
	var (
		err      error
		template domain.EmailTemplate
	)
	key := fmt.Sprintf("%s.%s_%s", domain.KEmailTemplate.Key, name, language)
	v, err := conn.Get(key).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &template); err != nil {
			return nil, err
		}
		return &template, nil
	}
	pTemplate, err := GetEmailTemplateFromDB(name, language)
	if err != nil {
		return nil, err

	}
	v, err = json.Marshal(pTemplate)
	if err != nil {
		return nil, err
	}
	if err := conn.Set(key, v, domain.KEmailTemplate.Timeout).Err(); err != nil {
		loggers.Warn.Printf("Cache %s sms template error:%s", name, err.Error())
	}
	return pTemplate, nil
}

func GetEmailTemplateFromDB(name, language string) (*domain.EmailTemplate, error) {
	db := confDB.NewConn()
	var template domain.EmailTemplate
	dbResult := db.Where("name = ?", name).First(&template)
	if dbResult.RecordNotFound() {
		loggers.Warn.Printf("Email template:%s not exist", name)
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		loggers.Warn.Printf("Get Email template:%s error:%s", name, dbResult.Error.Error())
		return nil, dbResult.Error
	}
	return &template, nil
}

func GetEmailService(email string) string {
	emailDomain := EmailDomain(email)
	services, err := GetEmailServices()
	if err != nil {
		return ""
	}
	if service, ok := (*services)[emailDomain]; ok {
		return service
	}
	if service, ok := (*services)[domain.EMAIL_SERVICE_DEFAULT]; ok {
		return service
	}
	return ""
}

func GetEmailServices() (*map[string]string, error) {
	conn := redisPool.NewConn()
	var services map[string]string
	v, err := conn.Get(domain.KEmailServies.Key).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &services); err == nil {
			return &services, nil
		}
		loggers.Error.Printf("Unmarshal email services from cache error %s", err.Error())
	}
	pServices, err := GetEmailServicesFromDB()
	if err != nil {
		loggers.Error.Printf("Get email servies config from db error %s", err.Error())
		return nil, err
	}
	value, err := json.Marshal(*pServices)
	if err != nil {
		loggers.Error.Printf("Marshal email servies config error %s", err.Error())
	}
	if err := conn.Set(domain.KEmailServies.Key, value, domain.KEmailServies.Timeout).Err(); err != nil {
		loggers.Error.Printf("Cache email servies config error %s", err.Error())
	}
	return pServices, nil
}

func GetEmailServicesFromDB() (*map[string]string, error) {
	db := confDB.NewConn()
	var services []domain.EmailService
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
		servicesMap[ins.Domain] = ins.Service
	}
	return &servicesMap, nil
}

func EmailDomain(email string) string {
	if len(email) < 1 {
		return email
	}
	arr := strings.Split(email, "@")
	if len(arr) != 2 {
		return email
	}
	return strings.Trim(arr[1], " ")
}

func SendEmail(email *domain.EmailMessage) error {
	if email == nil {
		return errors.New("email param is nil")
	}

	// Render
	tmpl, err := GetEmailTemplate(email.TemplateName, email.Language)
	if err != nil {
		return err
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"decimal": func(vol decimal.Decimal, places int32) string {
			return strings.TrimRight(vol.StringFixed(places), "0")
		},
	}).Parse(tmpl.Body)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, email.Source); err != nil {
		return err
	}

	// Send
	emailCtx := domain.Email{
		Sender:    email.Sender,
		Service:   GetEmailService(string(email.Recipient)),
		Recipient: email.Recipient,
		Subject:   email.Subject,
		HtmlBody:  buffer.String(),
		CharSet:   "UTF-8",
	}
	if emailCtx.Subject == "" {
		emailCtx.Subject = tmpl.Subject
	}
	if emailCtx.Sender == "" {
		emailCtx.Sender = "no-reply@bifund.com"
	}

	if err := emailSender.SendEmail(&emailCtx); err != nil {
		return err
	}

	return nil
}
