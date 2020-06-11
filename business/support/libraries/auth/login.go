package auth

import (
	"business/support/domain"
	"business/support/libraries/redis"
	"fmt"
	"time"
)

const MAX_LOGIN_RETRY = 5

// 登陆重试检查,超过5次,账号禁止登陆2小时
func LoginRetry(account int64) int64 {
	conn := redisPool.NewConn()
	retryKey := fmt.Sprintf("%s.%d", domain.KLoginRetryTimes.Key, account)
	count, _ := conn.Incr(retryKey).Result()
	if count >= 5 {
		disableKey := fmt.Sprintf("%s.%d", domain.KDisableLogin.Key, account)
		conn.Set(disableKey, time.Now().Unix(), domain.KDisableLogin.Timeout)
	}
	if count == 1 {
		conn.Expire(retryKey, domain.KLoginRetryTimes.Timeout)
	}

	return count
}

// 重置密码成功后,去除登陆限制
func EnableLogin(account int64) {
	conn := redisPool.NewConn()
	disableKey := fmt.Sprintf("%s.%d", domain.KDisableLogin.Key, account)
	conn.Del(disableKey)

}

// 账号是否被临时禁止登陆
func IsDisabledAccountLogin(account int64) bool {
	conn := redisPool.NewConn()
	disableKey := fmt.Sprintf("%s.%d", domain.KDisableLogin.Key, account)
	if err := conn.Get(disableKey).Err(); err != nil {
		return false
	}

	return true
}

// GetAccountDisabledTime 获得账号禁用时的时间戳
func GetAccountDisabledTime(account int64) int64 {
	conn := redisPool.NewConn()
	disableKey := fmt.Sprintf("%s.%d", domain.KDisableLogin.Key, account)
	nonce, err := conn.Get(disableKey).Int64()
	if err == redis.Nil {
		return 0
	}

	return nonce
}
