package msg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func newSqsClient(
	region string,
	accessKey string,
	accessSecret string,
	token string) *sqs.SQS {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, token),
	})
	if err != nil {
		return nil
	}
	client := sqs.New(sess)
	if nil == client {
		return nil
	}
	return client
}
