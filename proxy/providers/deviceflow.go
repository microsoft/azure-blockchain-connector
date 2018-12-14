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
	token  *deviceflow.Token
	client *http.Client
}

func (df *OAuthDeviceFlow) RequestAccess() (err error) {
	var ctx = context.Background()

	tok, err := aad.DeviceFlowGrant(ctx, df.Config)
	df.token = tok
	return
}

func (df *OAuthDeviceFlow) Client() *http.Client {
	return df.client
}

func (df *OAuthDeviceFlow) Modify(params *proxy.Params, req *http.Request) {
	req.Header.Set("Authorization", "Bearer"+" "+df.token.AccessToken)
}
