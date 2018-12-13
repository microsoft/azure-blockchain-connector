package providers

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/aad/deviceflow"
	"azure-blockchain-connector/proxy"
	"context"
	"net/http"
)

type OAuthDeviceFlow struct {
	*deviceflow.Config
	Token *deviceflow.Token
}

func (df *OAuthDeviceFlow) RequestAccess() (err error) {
	var ctx = context.Background()

	tok, err := aad.DeviceFlowGrant(ctx, df.Config)
	df.Token = tok
	return
}

func (df *OAuthDeviceFlow) Client(params *proxy.Params) *http.Client {
	return http.DefaultClient
}

func (df *OAuthDeviceFlow) Modify(params *proxy.Params, req *http.Request) {
	req.Header.Set("Authorization", "Bearer"+" "+df.Token.AccessToken)
}
