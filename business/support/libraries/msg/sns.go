package msg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func newSnsClient(
	region string,
	accessKey string,
	accessSecret string,
	token string) *sns.SNS {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, token),
	})
	if err != nil {
		return nil
	}
	client := sns.New(sess)
	if nil == client {
		return nil
	}
	return client
}
