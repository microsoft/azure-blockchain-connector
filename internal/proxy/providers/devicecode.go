package providers

import (
	"abc/internal/aad"
	"abc/internal/aad/devicecode"
	"abc/internal/proxy"
	"context"
	"net/http"
)

type OAuthDeviceCode struct {
	*devicecode.Config
	token  *devicecode.Token
	client *http.Client
}

func (df *OAuthDeviceCode) RequestAccess() (err error) {
	var ctx = context.Background()

	tok, err := aad.DeviceFlowGrant(ctx, df.Config)
	df.token = tok
	printToken(tok)

	df.client = &http.Client{}
	return
}

func (df *OAuthDeviceCode) Client() *http.Client {
	return df.client
}

func (df *OAuthDeviceCode) Modify(params *proxy.Params, req *http.Request) {
	req.Header.Set("Authorization", "Bearer"+" "+df.token.AccessToken)
}
