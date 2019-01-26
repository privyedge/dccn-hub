package wrapper

import (
	"context"
	"errors"
	"log"
	"os"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"

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
		userMgrService := usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, client.DefaultClient)
		_, err := userMgrService.VerifyToken(context.Background(), &usermgr.Token{Token: token})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}
