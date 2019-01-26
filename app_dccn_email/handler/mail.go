package handler

import (
	"context"
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dccnwrapper "github.com/Ankr-network/dccn-common/wrapper"
	"github.com/Ankr-network/dccn-hub/app_dccn_email/subscriber"
)

type Mail struct{}

// Send send mail
func (e *Mail) Send(ctx context.Context, req *common_proto.MailEvent, rsp *common_proto.Error) error {
	log.Println("Received Example.Call request")
	err := subscriber.SendMail(req)
	dccnwrapper.PbError(&rsp, err)
	return nil
}
