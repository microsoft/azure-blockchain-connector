package proxyproviders

import (
	"abc/internal/oauth2dc"
	"abc/internal/proxy"
	"context"
	"fmt"
	"net/http"
)

type OAuthDeviceCode struct {
	*oauth2dc.Config
	client *http.Client
}

func (df *OAuthDeviceCode) RequestAccess() (err error) {
	var ctx = context.Background()

	deviceAuth, err := df.Config.AuthDevice(ctx)
	if err != nil {
		return err
	}

	fmt.Println("Open:", deviceAuth.VerificationURI)
	fmt.Println("Enter:", deviceAuth.UserCode)

	tok, err := df.Config.Poll(ctx, deviceAuth)
	if err != nil {
		return
	}
	printToken(tok)

	df.client = df.Config.Client(ctx, tok)
	return
}

func (df *OAuthDeviceCode) Client() *http.Client {
	return df.client
}

func (df *OAuthDeviceCode) Modify(params *proxy.Params, req *http.Request) {
}
