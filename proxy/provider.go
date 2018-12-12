package proxy

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/aad/deviceflow"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

type Provider interface {
	RequestAccess() error
	Client(params *Params) *http.Client
	Modify(params *Params, req *http.Request)
}

type BasicAuth struct {
	Remote             string
	CertPath           string
	Insecure           bool
	Username, Password string
	pool               *x509.CertPool
	client             *http.Client
}

func mustReadPem(path string) []byte {
	if path == "" {
		return []byte("")
	}
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("ReadFile err:", err)
	}
	return pem
}

func (ba *BasicAuth) init() {
	ba.pool = x509.NewCertPool()
	ba.pool.AppendCertsFromPEM(mustReadPem(ba.CertPath))
	ba.client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            ba.pool,
			InsecureSkipVerify: ba.Insecure,
		},
		MaxIdleConnsPerHost: 1024,
	}}
}

func (ba *BasicAuth) test(remote string) error {
	req, err := http.NewRequest(http.MethodGet, "https://"+remote, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept-Encoding", "identity")
	req.SetBasicAuth(ba.Username, ba.Password)

	res, err := ba.client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 401 {
		return errors.New("unable to pass the authentication on the remote server")
	}
	return nil
}

func (ba *BasicAuth) RequestAccess() error {
	ba.init()
	err := ba.test(ba.Remote)
	return err
}

func (ba *BasicAuth) Client(params *Params) *http.Client {
	return ba.client
}

func (ba *BasicAuth) Modify(params *Params, req *http.Request) {
	req.SetBasicAuth(ba.Username, ba.Password)
}

type OAuthAuthCode struct {
	*oauth2.Config
	SvcAddr string
	Token   *oauth2.Token
	client  *http.Client
}

func (ac *OAuthAuthCode) RequestAccess() error {
	ctx := context.Background()
	tok, err := aad.AuthCodeGrant(ctx, ac.Config, ac.SvcAddr)
	if err != nil {
		return err
	}
	ac.client = ac.Config.Client(ctx, tok)
	ac.Token = tok
	return nil
}

func (ac *OAuthAuthCode) Client(params *Params) *http.Client {
	return ac.client
}

func (ac *OAuthAuthCode) Modify(params *Params, req *http.Request) {
}

type OAuthDeviceFlow struct {
	*deviceflow.Config
	Token *deviceflow.Token
}

func (df *OAuthDeviceFlow) RequestAccess() (err error) {
	var ctx = context.Background()

	tok, err := aad.DeviceFlowGrant(ctx, df.Config)
	if tok != nil {
		fmt.Println("Token:", tok.AccessToken)
		fmt.Println("Expires in:", tok.ExpiresIn)
	}
	df.Token = tok
	return
}

func (df *OAuthDeviceFlow) Client(params *Params) *http.Client {
	return http.DefaultClient
}

func (df *OAuthDeviceFlow) Modify(params *Params, req *http.Request) {
	req.Header.Set("Authorization", "Bearer"+" "+df.Token.AccessToken)
}
