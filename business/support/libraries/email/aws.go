package email

import (
	"business/support/config"
	"business/support/domain"
	"business/support/libraries/loggers"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func SendByAws(email *domain.Email) error {

	if email.CharSet == "" {
		email.CharSet = "UTF-8"
	}

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Conf.AWSConf.Region),
		Credentials: credentials.NewStaticCredentials(config.Conf.AWSConf.AccessKey, config.Conf.AWSConf.AccessSecret, ""),
	})

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(string(email.Recipient)),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(email.CharSet),
					Data:    aws.String(email.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(email.CharSet),
					Data:    aws.String(email.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(email.CharSet),
				Data:    aws.String(email.Subject),
			},
		},
		Source: aws.String(email.Sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				loggers.Warn.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				loggers.Warn.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				loggers.Warn.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				loggers.Warn.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			loggers.Warn.Println(err.Error())
		}

		return err
	}
	loggers.Info.Printf("SendEmail from:%s to:%s id:%s", email.Sender, email.Recipient, aws.StringValue(result.MessageId))
	return nil
}
