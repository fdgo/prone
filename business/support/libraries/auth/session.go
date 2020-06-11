package auth

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/errors"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

func SessionConfAdaptor(version string) *domain.SessionConf {
	sessionConf, ok := config.Conf.SessionConf[version]
	if !ok {
		sessionConf, ok = config.Conf.SessionConf["1.0"]
		if ok {
			return &sessionConf
		} else {
			return nil
		}
	} else {
		return &sessionConf
	}
}

func SAuth(uid int64, sign, version, dev string, ts int64) (*domain.Session, error) {
	sessionConf := SessionConfAdaptor(version)
	if nil == sessionConf {
		return nil, errors.Forbidden
	}

	iv := fmt.Sprintf("%d", ts)
	token, err := CBCDecrypt([]byte(sessionConf.Key), []byte(iv), sign)
	if err != nil {
		return nil, errors.Forbidden
	}

	session, err := GetSession(uid, dev)
	if err != nil {
		return nil, errors.Forbidden
	}

	if token != session.Token {
		return nil, errors.Forbidden
	}

	if !CheckTs(ts) {
		return nil, errors.ClientTimeInvalidError
	}

	if err := SetSessionTtl(uid, version, dev); err != nil {
		return nil, errors.Forbidden
	}

	return session, nil
}

func GetSession(uid int64, dev string) (*domain.Session, error) {
	conn := redisPool.NewConn()
	key := SessionKey(uid, dev)
	val, err := conn.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	var session domain.Session
	if err := json.Unmarshal(val, &session); err != nil {
		return nil, err
	}

	session.Account, err = GetAccountByUid(uid)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func SetSessionTtl(uid int64, version, dev string) error {
	sessionConf := SessionConfAdaptor(version)
	if nil == sessionConf {
		return errors.New("Invalid version")
	}
	conn := redisPool.NewConn()
	key := SessionKey(uid, dev)
	if f, err := conn.Expire(key, sessionConf.Timeout).Result(); err != nil || !f {
		return err
	}
	return nil
}

func DeleteSession(uid int64, dev string) error {
	conn := redisPool.NewConn()
	key := SessionKey(uid, dev)
	if _, err := conn.Del(key).Result(); err != nil {
		return err
	}
	return nil
}

func NewSession(account *domain.Account, version, dev string, ts int64) (*domain.Session, error) {
	sessionConf := SessionConfAdaptor(version)
	if nil == sessionConf {
		return nil, errors.New(fmt.Sprintf("Invalid version:%s", version))
	}
	token, err := GenerateUUIDString()
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	session := domain.Session{
		Ts:        ts,
		Token:     token,
		Account:   account,
		CreatedAt: &now,
	}

	conn := redisPool.NewConn()
	key := SessionKey(account.AccountId, dev)
	val, err := json.Marshal(&session)
	if err != nil {
		return nil, err
	}

	if err := conn.Set(key, string(val), sessionConf.Timeout).Err(); err != nil {
		return nil, err
	}

	return &session, nil
}

func GenerateUUIDString() (string, error) {
	u, err := uuid.NewV1()
	if err != nil {
		return "", err
	}
	str := ByteToHex(u.Bytes())
	return str, nil
}

func ByteToHex(data []byte) string {
	buffer := new(bytes.Buffer)
	for _, b := range data {
		buffer.WriteString(fmt.Sprintf("%02x", int64(b&0xff)))
	}

	return buffer.String()
}

func SessionKey(uid int64, dev string) string {
	return fmt.Sprintf("SESSION.%d_%s", uid, dev)
}

func CBCDecrypt(key, iv []byte, sign string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(iv)%block.BlockSize() != 0 {
		return "", errors.New("Invalid iv length")
	}
	ciphertext, err := hex.DecodeString(sign)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return string(ciphertext), nil
}

func CheckTs(ts int64) bool {
	curTs := time.Now().UTC().Unix()
	if (curTs-ts/1000000) > 1800 || (curTs-ts/1000000) < -1800 {
		return false
	}
	return true
}

func SetAssetSession(token string, timeout time.Duration) error {
	key := fmt.Sprintf("SESSION.ASSET.%s", token)
	conn := redisPool.NewConn()
	if err := conn.Set(key, 1, timeout).Err(); err != nil {
		return err
	}
	return nil
}

func AssetTokenAuth(token string, timeout time.Duration) bool {
	key := fmt.Sprintf("SESSION.ASSET.%s", token)
	conn := redisPool.NewConn()
	if err := conn.Get(key).Err(); err != nil {
		return false
	}
	return true
}

func ResetAssetSession(newToken, oldToken string, timeout time.Duration) bool {
	conn := redisPool.NewConn()
	newKey := fmt.Sprintf("SESSION.ASSET.%s", newToken)
	oldKey := fmt.Sprintf("SESSION.ASSET.%s", oldToken)
	if err := conn.Rename(oldKey, newKey).Err(); err != nil {
		return false
	}
	if f, err := conn.Expire(newKey, timeout).Result(); err != nil || !f {
		return false
	}
	return true
}
