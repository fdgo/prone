package config

import (
	"business/support/config/decrypt"
	"business/support/domain"
	"business/support/libraries/loggers"
	"business/support/localconfig"
	"business/support/rpc"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/shopspring/decimal"
)

type Config struct {
	DBS         map[string]domain.DBConf
	Services    map[string]domain.ServerConf
	SessionConf map[string]domain.SessionConf
	//ContractAddress    *domain.ETHADDRESS
	MQConf             *domain.MQConf
	ArchiveS3          *domain.S3Conf
	AWSConf            *domain.AWSConf
	SNSConf            *domain.SNSConf
	SQSConf            *domain.SQSConf
	ActionSNSConf      *domain.SNSConf
	RecordSQSConf      *domain.SQSConf
	NotifySQSConf      *domain.SQSConf
	CloudNotifySQSConf *domain.SQSConf
	ActivitySQSConf    *domain.SQSConf
	BackStageSQSConf   *domain.SQSConf
	SMSSQSConf         *domain.SQSConf
	EmailSQSConf       *domain.SQSConf
	ContractDCSQSConf  *domain.SQSConf
	SpotDCSQSConf      *domain.SQSConf
	SendCloudEmailConf *domain.SendCloudConf
	SendCloudSMSConf   *domain.SendCloudConf
	OSSInfo            *domain.OSSInfo
	CaptchConf         *domain.CaptchConf
	NECConf            *domain.NECConf
	Consuls            []string
	LogLevel           string
	Precision          int64 // 精度
	AdminUserId        int64
	IPDatPath          string `default:"/data/etc/business/qqwry.dat"`
}

var (
	Conf    Config
	IsDebug bool
)

var (
	configRPC     *rpc.ConfigRPC
	configFileDir = "/data/etc/business"
	key           = "ahduf*^$%$%8586768&*(&@!,m;l,m;g"
	rpckey        = "s01RX9MMjHpLUrI0wc7OavBtu7J+HP5SvueGuMmowJZkI1wNNvX6903zwSo7tR1kazKsoFbmth7GEcqymPgbT715lq/1unKrCfYkll67zjI="
	//datakey = &*IUTGkq12hgdfghshg%#%@@*q754979
	datakey = "6M4RAYqf5iAkgAKFhbJSTQkAoXuLKzsbkFC20+zaPKdEIV8dQUEBPItlCMp3Ix5W"
)

func init() {
	decimal.DivisionPrecision = 19
	IsDebug = localconfig.IsDebug
	if localconfig.IsDebug {
		if err := configor.Load(&Conf, localconfig.ConfigFilePath); err != nil {
			fmt.Printf("config Read config debug file:%s error:%s", localconfig.ConfigFilePath, err.Error())
		}

		fmt.Printf("Debug Config Load %s \n", localconfig.ConfigFilePath)
		return
	}

	configRPC = rpc.NewConfigRPC()
	var rd domain.ConfigRequestData
	if len(localconfig.Conf.Ip) > 0 {
		rd.Ip = localconfig.Conf.Ip
	}

	d, err := json.Marshal(rd)
	if err != nil {
		panic(errors.New("config Marshal error"))
	}

	var req domain.ConfigRequest
	req.Data = d

	sdec, _ := base64.StdEncoding.DecodeString(rpckey)
	sk := decrypt.AESDecrypt(sdec, []byte(key))
	if len(sk) == 0 {
		panic(errors.New("config rpckey error"))
	}

	sign, e := decrypt.KSign(req.Data, string(sk[:]))
	if e != nil {
		panic(errors.New("config sign data error"))
	}

	req.ActionSign = sign

	var config domain.ConfigResp

	e = configRPC.GetGlobalConfig(context.Background(), &req, &config)
	if e != nil || len(config.Data) == 0 {
		panic(errors.New("config data nil"))
	}

	bdk, _ := base64.StdEncoding.DecodeString(datakey)
	dk := decrypt.AESDecrypt(bdk, []byte(key))
	if len(dk) == 0 {
		panic(errors.New("config key error"))
	}

	data := decrypt.AESDecrypt(config.Data, []byte(dk))
	if data == nil || len(data) == 0 {
		panic(errors.New("config data destroyed"))
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	configFilePath := fmt.Sprintf("%s/%d.yml", configFileDir, r.Uint32())
	e = ioutil.WriteFile(configFilePath, data, 0770)
	if e != nil {
		loggers.Error.Printf("config Write temp configfile fail %v", e)
		panic(errors.New("config Write temp file fail"))
	}

	if err := configor.Load(&Conf, configFilePath); err != nil {
		loggers.Error.Printf("config Read config file:%s error:%s", configFilePath, err.Error())
		panic(errors.New("config load config error"))
	}

	e = os.Remove(configFilePath)
	loggers.Debug.Println("config init config success", e)
}
