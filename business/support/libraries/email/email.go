package email

import (
	"business/support/domain"
	"business/support/libraries/redis"
	"encoding/json"
	"strings"
)

var (
	pool *redis.RedisPool
)

func Init(redisPool *redis.RedisPool) {
	pool = redisPool
}

func SendMailAsync(email *domain.Email) error {
	v, err := json.Marshal(email)
	if err != nil {
		return err
	}
	conn := pool.NewConn()
	key := domain.KWaitSendEmails
	if err := conn.RPush(key.Key, v).Err(); err != nil {
		return err
	}
	return nil
}

func SendEmail(email *domain.Email) error {
	if email.Service == domain.EMAIL_SERVIEC_SENDCLOUD {
		return SendBySendCloud(email, nil)
	}
	if email.Service == domain.EMAIL_SERVICE_AWSSES {
		return SendByAws(email)
	}
	if isForeignEmail(string(email.Recipient)) {
		return SendByAws(email)
	} else {
		return SendBySendCloud(email, nil)
	}
}

//判断是否海外邮箱
func isForeignEmail(email string) bool {
	if strings.HasSuffix(email, "@gmail.com") ||
		strings.HasSuffix(email, "@live.com") ||
		strings.HasSuffix(email, "@msn.com") ||
		strings.HasSuffix(email, "@hotmail.com") ||
		strings.HasSuffix(email, "@outlook.com") ||
		strings.HasSuffix(email, "@yahoo.com") ||
		strings.HasSuffix(email, "@mail.com") ||
		strings.HasSuffix(email, "@aim.com") ||
		strings.HasSuffix(email, "@outlook.com") ||
		strings.HasSuffix(email, "@aol.com") {

		return true
	}
	return false
}
