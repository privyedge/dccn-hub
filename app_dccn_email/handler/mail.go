package handler

import (
	"context"

	"github.com/Ankr-network/dccn-hub/app_dccn_email/subscriber"

	mail "github.com/Ankr-network/dccn-hub/app_dccn_email/proto/v1"
	log "github.com/micro/go-log"
)

type Mail struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Mail) Send(ctx context.Context, req *mail.MailEvent, rsp *mail.Response) error {
	log.Log("Received Example.Call request")
	return subscriber.SendMail(req)
}
