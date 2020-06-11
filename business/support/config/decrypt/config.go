package decrypt

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

var (
	configFileDir = "/data/etc/business"
)

func EncryptConfig(config, k string) (bool, error) {
	src, e := ioutil.ReadFile(config)
	if e != nil {
		return false, e
	}

	data := AESEncrypt(src, []byte(k))
	if len(data) == 0 {
		return false, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sr := fmt.Sprintf("%s/%d.yml", configFileDir, r.Uint32())
	e = ioutil.WriteFile(sr, data, 0770)
	if e != nil {
		return false, e
	}

	e = os.Rename(sr, config)
	if e != nil {
		return false, e
	}

	return true, nil
}

func DecryptConfig(config, k string) (bool, error) {
	src, e := ioutil.ReadFile(config)
	if e != nil {
		return false, e
	}

	data := AESDecrypt(src, []byte(k))
	if len(data) == 0 {
		return false, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sr := fmt.Sprintf("%s/%d.dat", configFileDir, r.Uint32())
	e = ioutil.WriteFile(sr, data, 0770)
	if e != nil {
		return false, e
	}

	e = os.Rename(sr, config)
	if e != nil {
		return false, e
	}

	return true, nil
}
