package handler

import (
	"context"
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"

	"github.com/Ankr-network/dccn-hub/app-dccn-email/subscriber"
)

type MailHandler struct {
}

func New() *MailHandler {
	return &MailHandler{}
}

// Send send mail
func (p *MailHandler) Send(ctx context.Context, e *mail.MailEvent, rsp *common_proto.Empty) error {

	log.Println("Debug send email")
	return subscriber.NewSender(e).SendEmail()
}
