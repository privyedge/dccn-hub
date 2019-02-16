package subscriber

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "no-reply@ankr.network"

	// Specify a configuration set. If you do not want to use a configuration
	// set, comment out the following constant and the
	// ConfigurationSetName: aws.String(ConfigurationSet) argument below
	ConfigurationSet = "ConfigSet"

	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	AwsRegion = "us-west-2"

	// The subject line for the email.
	Subject = "Amazon SES Test (适用于 Go 的 AWS 开发工具包)"

	// The HTML body for the email.
	HtmlBody = "<h1> Activate User </h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'></a>.</p>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the 适用于 Go 的 AWS 开发工具包."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func SendEmail(e *common_proto.MailEvent) error {

	// Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion),
	},
	)

	// Create an SES client in the session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(e.ToAddress),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(Sender),
		// Comment or remove the following line if you are not using a configuration set
		ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println("Email Sent!")
	fmt.Println(result)
	return nil
}

func Handler(ctx context.Context, e *common_proto.MailEvent) error {
	log.Println("Function Received message: ", e)

	return SendEmail(e)
}
