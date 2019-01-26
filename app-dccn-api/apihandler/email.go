package apihandler

import (
	"context"

	"github.com/micro/go-micro/client"

	ankr_default "github.com/Ankr-network/dccn-common/protos"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

type ApiEmail struct {
	api mail.MailService
}

func (p *ApiEmail) Send(ctx context.Context, req *common_proto.MailEvent, rsp *common_proto.Error) error {
	out, _ := p.api.Send(ctx, req)
	*rsp = *out
	return nil
}

func NewApiEmail(c client.Client) *ApiEmail {
	return &ApiEmail{
		api: mail.NewMailService(ankr_default.EmailRegistryServerName, c),
	}
}
