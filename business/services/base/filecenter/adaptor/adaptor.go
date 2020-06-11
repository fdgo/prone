package adaptor

import (
	"business/support/config"
	"business/support/domain"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	uploader *s3manager.Uploader
	bucket   *oss.Bucket
)

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Conf.ArchiveS3.Region),
		Credentials: credentials.NewStaticCredentials(config.Conf.ArchiveS3.AccessKey, config.Conf.ArchiveS3.AccessSecret, ""),
	})
	if err != nil {
		panic(err)
	}
	uploader = s3manager.NewUploader(sess)

	bucket, err = newBucket(config.Conf.OSSInfo)
	if err != nil {
		panic(err)
	}
}

func newBucket(info *domain.OSSInfo) (*oss.Bucket, error) {
	client, err := oss.New(info.EndPoint, info.AccessID, info.AccessKey)
	if err != nil {
		return nil, err
	}

	//if err := client.CreateBucket(info.Bucket); err != nil {
	//	return nil, err
	//}

	bucket, err := client.Bucket(info.Bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

func PutObject(objectID string, reader io.Reader, options ...oss.Option) error {
	return bucket.PutObject(objectID, reader, options...)
}
