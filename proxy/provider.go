package proxy

import (
	"azure-blockchain-connector/aad/deviceflow"
	"crypto/x509"
	"golang.org/x/oauth2"
	"net/http"
)

type Provider interface {
	Client(params *Params) *http.Client
	Modify(params *Params, req *http.Request)
}

type BasicAuth struct {
	CertPath           string
	Insecure           bool
	Username, Password string
	pool               *x509.CertPool
}

func (p *BasicAuth) Client(params *Params) *http.Client {
	return http.DefaultClient
}

func (p *BasicAuth) Modify(params *Params, req *http.Request) {

}

type OAuthAuthCode struct {
	AuthCodeConf *oauth2.Config
	AuthSvcAddr  string
}

func (p *OAuthAuthCode) Client(params *Params) *http.Client {
	return http.DefaultClient
}

func (p *OAuthAuthCode) Modify(params *Params, req *http.Request) {

}

type OAuthDeviceFlow struct {
	DeviceFlowConf *deviceflow.Config
}

func (p *OAuthDeviceFlow) Client(params *Params) *http.Client {
	return http.DefaultClient
}

func (p *OAuthDeviceFlow) Modify(params *Params, req *http.Request) {

}
