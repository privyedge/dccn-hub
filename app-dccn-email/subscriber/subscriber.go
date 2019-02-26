package subscriber

import (
	"context"
	"log"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
)

type Subscriber struct {
}

func (p *Subscriber) Handler(ctx context.Context, e *mail.MailEvent) error {
	log.Println("Function Received message: ", e)

	sender := Sender{MailEvent: e}
	return sender.SendEmail()
}
