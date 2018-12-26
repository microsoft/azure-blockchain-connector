package providers

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/aad/oauth2/devicecode"
	"azure-blockchain-connector/proxy"
	"context"
	"fmt"
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
	fmt.Println(tok.AccessToken)
	return
}

func (df *OAuthDeviceCode) Client() *http.Client {
	return df.client
}

func (df *OAuthDeviceCode) Modify(params *proxy.Params, req *http.Request) {
	req.Header.Set("Authorization", "Bearer"+" "+df.token.AccessToken)
}
