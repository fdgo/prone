package auth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"time"
)

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func CurGACode(secret string) (uint32, error) {
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return 0, errors.New("Invalid GAKey")
	}
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(toBytes(time.Now().UTC().Unix() / 30))
	hash := hmacSha1.Sum(nil)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := toUint32(hashParts)
	pwd := number % 1000000
	return pwd, nil
}
