package config

import (
	"business/support/config/decrypt"
	"encoding/base64"
	"fmt"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	for k, v := range Conf.DBS {
		fmt.Printf("%-20s %#v\n", k, v)
	}
	for k, v := range Conf.Services {
		fmt.Printf("%-20s %#v\n", k, v)
	}
	for k, v := range Conf.SessionConf {
		fmt.Printf("%-20s %#v\n", k, v)
	}

	fmt.Printf("%-20s %#v\n", "MQConf", Conf.MQConf)
	fmt.Printf("%-20s %#v\n", "Precision", Conf.Precision)
	fmt.Printf("%-20s %#v\n", "AdminUserId", Conf.AdminUserId)

}

func Test_AESEncrypt(t *testing.T) {
	d := decrypt.AESEncrypt([]byte("&*IUTGkq12hgdfghshg%#%@@*q754979"), []byte("ahduf*^$%$%8586768&*(&@!,m;l,m;g"))
	s := base64.StdEncoding.EncodeToString(d)

	fmt.Println(s)

	sDec, _ := base64.StdEncoding.DecodeString(s)

	db := decrypt.AESDecrypt(sDec, []byte("ahduf*^$%$%8586768&*(&@!,m;l,m;g"))

	fmt.Println(string(db[:]))
}
