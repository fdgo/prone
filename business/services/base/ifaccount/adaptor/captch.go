package adaptor

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
)

func IncrVerifyCodeCaptchLimit(action, ip string) bool {
	var (
		conf = config.Conf.CaptchConf
		key  = fmt.Sprintf("%s.%s.%s", domain.KCaptchLimit.Key, action, ip)
		conn = redisPool.NewConn()
	)
	count := conn.Incr(key).Val()
	if count == 1 {
		conn.Expire(key, config.Conf.CaptchConf.LimitTime)
	}
	return count > conf.LimitVerifyCodeReq
}

func CheckCaptchLimit(action, ip string) bool {
	// 防止恶意撞库,强制要求验证
	return true
	switch action {
	case "login":
		return CheckLoginCaptchLimit(ip)
	default:
		return checkVerifyCodeReqLimit(action, ip)
	}
}

func CheckLoginCaptchLimit(ip string) bool {
	start, _ := time.ParseInLocation("2006-01-02", time.Now().In(localTime).Format("2006-01-02"), localTime)
	end := start.AddDate(0, 0, 1)
	loginRecords, err := GetLoginRecordsByIP(ip, &start, &end)
	if err != nil {
		return false
	}
	accounts := make(map[int64]bool)
	for i := range loginRecords {
		accounts[loginRecords[i].AccountId] = true
	}
	return int64(len(accounts)) > config.Conf.CaptchConf.LimitLoginAccount
}

func checkVerifyCodeReqLimit(action, ip string) bool {
	var (
		conf = config.Conf.CaptchConf
		key  = fmt.Sprintf("%s.%s.%s", domain.KCaptchLimit.Key, action, ip)
		conn = redisPool.NewConn()
	)
	count, err := conn.Get(key).Int64()
	if err != nil {
		if !redisPool.IsNil(err) {
			loggers.Error.Printf("CheckCaptchLimit error:%s", err.Error())
		}
		return false
	}

	return count >= conf.LimitVerifyCodeReq
}

func GetLoginRecordsByIP(ip string, start, end *time.Time) ([]domain.AccountRecord, error) {
	var accountRecords []domain.AccountRecord
	if err := dbRecordsPool.NewConn().
		Where("client_ip = ?", ip).
		Where("log_time >= ?", start).
		Where("log_time < ?", end).
		Find(&accountRecords).Error; err != nil {
		return nil, err
	}

	return accountRecords, nil
}

func genNECSignature(secretKey string, form map[string]string) string {
	var keys []string
	for key, _ := range form {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buf := bytes.NewBufferString("")
	for _, key := range keys {
		buf.WriteString(key + form[key])
	}
	buf.WriteString(secretKey)
	has := md5.Sum(buf.Bytes())
	return fmt.Sprintf("%x", has)
}

func CheckCaptchValidate(validate string) bool {
	trans := &http.Transport{
		MaxIdleConns:       4,
		IdleConnTimeout:    time.Second * 30,
		DisableCompression: true,
	}
	client := http.Client{
		Transport: trans,
		Timeout:   time.Second * 10,
	}

	form := url.Values{
		"captchaId": {config.Conf.NECConf.CaptchID},
		"validate":  {validate},
		"user":      {""},
		"secretId":  {config.Conf.NECConf.SecretID},
		"version":   {"v2"},
		"timestamp": {strconv.FormatInt(time.Now().Unix()*1000, 10)},
		"nonce":     {strconv.FormatInt(100000+int64(rand.Int31n(899999)), 10)},
	}
	form.Add("signature", genSignature(config.Conf.NECConf.SecretKey, form))

	req, err := http.NewRequest(http.MethodPost, config.Conf.NECConf.CaptchAPI, bytes.NewReader([]byte(form.Encode())))
	if err != nil {
		loggers.Error.Printf("Request NECCaptch api error:%s", err.Error())
		return true
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		loggers.Error.Printf("Request NECCaptch api error:%s", err.Error())
		return true
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return true
	}
	var result struct {
		Msg    string `json:"msg"`
		Error  int    `json:"error"`
		Result bool   `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		loggers.Error.Printf("Request NECCaptch api error:%s", err.Error())
		return true
	}

	if result.Error == 0 && result.Result {
		return true
	}
	loggers.Warn.Printf("CheckCaptchValidate error:%d result:%v msg:%s", result.Error, result.Result, result.Msg)
	return false

}

func genSignature(secretKey string, form url.Values) string {
	var keys []string
	for key := range form {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	buf := bytes.NewBufferString("")
	for _, key := range keys {
		buf.WriteString(key + form.Get(key))
	}
	buf.WriteString(secretKey)
	has := md5.Sum(buf.Bytes())
	return fmt.Sprintf("%x", has)
}
