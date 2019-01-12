package wrapper

import (
	"context"
	"errors"
	"log"
	"os"

	go_micro_srv_usermgr "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"
)

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		// Really shouldn't be using a global here, find a better way
		// of doing this, since you can't pass it into a wrapper.
		userMgrService := go_micro_srv_usermgr.NewUserMgrService("go.micro.srv.usermgr", client.DefaultClient)
		_, err := userMgrService.VerifyToken(context.Background(), &go_micro_srv_usermgr.Token{Token: token})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
