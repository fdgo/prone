package adaptor

import (
	"fmt"
	"io"

	"business/support/libraries/loggers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func Upload(reader io.Reader, key, contentType, bucket, acl string) error {
	input := s3manager.UploadInput{
		ACL:         aws.String(acl),
		Bucket:      aws.String(bucket),
		ContentType: aws.String(contentType),
		Key:         aws.String(key),
		Body:        reader,
	}

	result, err := uploader.Upload(&input)
	if err != nil {
		fmt.Println(err)
		return err
	}

	loggers.Debug.Printf("Upload OK %s %s %s", key, bucket, result.Location)
	return nil
}
