package handler

import (
	"context"
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	"github.com/Ankr-network/dccn-hub/app-dccn-email/subscriber"
)

type Mail struct{}

// Send send mail
func (e *Mail) Send(ctx context.Context, req *common_proto.MailEvent, rsp *common_proto.Error) error {

	log.Println("Received Example.Call request")
	return subscriber.SendEmail(req)
}
