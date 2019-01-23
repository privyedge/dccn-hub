package handler

import (
	"context"
	"log"

	"github.com/Ankr-network/dccn-hub/app_dccn_email/subscriber"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
)

type Mail struct{}

// Call is a single request apihandler called via client.Call or the generated client code
func (e *Mail) Send(ctx context.Context, req *mail.MailEvent, rsp *mail.Response) error {
	log.Println("Received Example.Call request")
	return subscriber.SendMail(req)
}
