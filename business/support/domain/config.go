package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type DBConf struct {
	Host         string `require:"true"`
	Port         int    `require:"true"`
	Username     string
	Password     string
	DBName       string
	MaxIdleConns int `default:"4"`
	MaxOpenConns int `default:"8"`
}

type ServerConf struct {
	HostName        string `env:"HOSTNAME"`
	Name            string `env:"SERVICENAME"`
	ListenAddress   string `default:":8081"`
	RPCAddress      string `default:":8082"`
	EthereumNetwork string
	EtherscanApiKey string
	ContractAddress string
	PubKey          string
	PriKey          string
	RPCGroup        int64 // ToDo Glen Contract 正式上线时需要记得改下线上环境的配置文件
	Whitelist       []string
	OmniNode        []*RpcConf
	EOSNode         *RpcConf
	EOSAccount      string
	KYCRewardMax    decimal.Decimal
	NormalRewardMax decimal.Decimal
	EthWallet       string //资金归集目的地
}

type MQConf struct {
	Host    string
	Port    int
	Network string
}

type SessionConf struct {
	Key     string
	Timeout time.Duration
}

type AWSConf struct {
	Region       string
	AccessKey    string
	AccessSecret string
}

type SNSConf struct {
	Region              string
	AccessKey           string
	AccessSecret        string
	TopicArn            string
	ActionTopicArn      string
	AssetRecordTopicArn string
}

type SQSConf struct {
	Region              string
	AccessKey           string
	AccessSecret        string
	QueueUrl            string
	BaseUrl             string
	SMSNotifyUrl        string
	EmailNotifyUrl      string
	NotifyDLQUrl        string
	ActivityQueueUrl    string
	AssetRecordQueueUrl string
}

type SendCloudConf struct {
	ApiUrl  string
	ApiUser string
	ApiKey  string
}

type RpcConf struct {
	Host string
	Port int
	User string
	Pass string
	Net  string
}

type S3Conf struct {
	Region       string
	Bucket       string
	AccessKey    string
	AccessSecret string
}

type CaptchConf struct {
	LimitTime          time.Duration
	LimitLoginAccount  int64
	LimitVerifyCodeReq int64
}

type NECConf struct {
	CaptchID  string
	SecretKey string
	SecretID  string
	CaptchAPI string
	SMSAPI    string
	SMSID     string
}

var DevMap = map[string]bool{
	"ios":     true,
	"android": true,
	"web":     true,
	"api":     true,
}

type ConfigRequest struct {
	Data       []byte
	ActionSign string
}

type ConfigRequestData struct {
	Ip string
}

type ConfigResp struct {
	Data []byte
}
