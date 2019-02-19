package subscriber

import (
	"log"
	"testing"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
)

func TestSendEmail(t *testing.T) {
	sender := NewSender(&mail.MailEvent{
		Type: mail.EmailType_CONFIRM_REGISTRATION,
		From: ankr_default.NoReplyEmailAddress,
		To:   []string{"xuexiacm@163.com"},
		OpMail: &mail.MailEvent_ConfirmRegistration{ConfirmRegistration: &mail.ConfirmRegistration{
			UserName: "test",
			UserId:   "00001",
			Code:     "sanghaifa",
		}},
	})

	if err := sender.SendEmail(); err != nil {
		log.Fatal(err.Error())
	}
}
