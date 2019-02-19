package subscriber

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	email_templates "github.com/Ankr-network/dccn-hub/app-dccn-email/templates"
)

const (

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Sender struct {
	*mail.MailEvent
}

func NewSender(e *mail.MailEvent) *Sender {
	return &Sender{MailEvent: e}
}

func (p *Sender) subject() string {
	return "Welcome"
}

func (p *Sender) textBody() string {
	return "Welcome"
}

var emailTemplates = map[string]*template.Template{
	"registeration": template.Must(template.New("registration").Parse(email_templates.RegistrationTemplate)),
}

func (p *Sender) htmlBody() string {
	// The HTML body for the email.
	var tpl bytes.Buffer
	var html string
	switch p.Type {
	case mail.EmailType_CONFIRM_REGISTRATION:
<<<<<<< HEAD
		t := emailTemplates["registeration.html"]
=======
		t := emailTemplates["registeration"]
>>>>>>> f2970ea8e2887d3244c6fad220a19379eca8dc05
		data := struct {
			Code string
			ID   string
		}{p.GetConfirmRegistration().Code, p.GetConfirmRegistration().UserId}
		t.Execute(&tpl, data)
		html = tpl.String()
		// log.Print(data)
		// log.Printf("html: %s", html)
	case mail.EmailType_FORGET_PASSWORD:
		code := p.GetForgetPassword().Code
		email := p.GetForgetPassword().Email
		html = fmt.Sprintf("<h1>FORGET_PASSWORD %s(Validate Code)</h1><p>url<a href='https://domain.com/verify/code=%s?email=%s'></a>", code, code, email)
		log.Printf("user: %s, code: %s", email, code)
	case mail.EmailType_CHANGE_PASSWORD:
		id := p.GetChangePassword().UserId
		code := p.GetChangePassword().Code
		html = fmt.Sprintf("<h1>CHANGE_PASSWORD %s(Validate Code)</h1><p>url<a href='https://domain.com/verify/code=%s?email=%s'></a>", code, code, id)
		log.Printf("user: %s, code: %s", id, code)
	case mail.EmailType_CONFIRM_EMAIL:
		id := p.GetChangeEmail().UserId
		code := p.GetChangeEmail().Code
		html = fmt.Sprintf("<h1>CONFIRM_EMAIL %s(Validate Code)</h1><p>url<a href='https://domain.com/verify/code=%s?email=%s'></a>", code, code, id)
		log.Printf("user: %s, code: %s", id, code)
	}

	return html
}

func (p *Sender) input() *ses.SendEmailInput {

	// Assemble the email.
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(p.To),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(p.htmlBody()),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(p.textBody()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(p.subject()),
			},
		},
		Source: aws.String(p.From),
	}
}

func (p *Sender) SendEmail() error {

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Attempt to send the email.
	result, err := svc.SendEmail(p.input())

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

		return err
	}

	fmt.Println("Email Sent to address: " + strings.Join(p.To, " "))
	fmt.Println(result)
	return nil
}
