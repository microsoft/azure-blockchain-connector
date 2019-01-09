package proxyproviders

import (
	"abc/internal/aad/devicecode"
	"abc/internal/proxy"
	"context"
	"net/http"
)

type OAuthDeviceCode struct {
	*devicecode.Config
	client *http.Client
}

func (df *OAuthDeviceCode) RequestAccess() (err error) {
	var ctx = context.Background()

	tok, err := df.Config.Grant(ctx)
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
