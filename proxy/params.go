package proxy

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/aad/deviceflow"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
)

const (
	OAuthGrantTypeNone       = ""
	OAuthGrantTypeAuthCode   = "authcode"
	OauthGrantTypeDeviceFlow = "device"

	LogWhenOnError  = "onError"  // print log only for those who raise exceptions
	LogWhenOnNon200 = "onNon200" // print log for those who have a non-200 response, or those who raise exceptions
	LogWhenAlways   = "always"   // print log for every request

	LogWhatBasic    = "basic"    // print the request's method and URI and the response status code (and the exception message, if exception raised) in the log
	LogWhatDetailed = "detailed" // print the request's method, URI and body, and the response status code and body (and the exception message, if exception raised) in the log
	//LogAll          = "all"      // to be supported later. Compared to whatlog_detail, all Headers are printed in whatlog_all
)

type Params struct {
	Local, Remote string

	CertPath string
	pool     *x509.CertPool
	Insecure bool

	Username, Password string
	AuthCodeConf       *oauth2.Config
	AuthSvcAddr        string
	DeviceFlowConf     *deviceflow.Config

	Client *http.Client

	Whenlog string
	Whatlog string
}

func (params *Params) SetOAuthConfig(typ, clientID, clientSecret, scopes string) {
	switch typ {
	case OAuthGrantTypeAuthCode:
		params.AuthCodeConf = aad.NewAuthCodeConfig(clientID, clientSecret, scopes)
	case OauthGrantTypeDeviceFlow:
		params.DeviceFlowConf = aad.NewDeviceFlowConfig(clientID, scopes)
	}
}

func (params *Params) initCACert() {
	var caCertPath = params.CertPath
	if caCertPath != "" {
		caCrt, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			fmt.Println("ReadFile err:", err)
			return
		}
		params.pool.AppendCertsFromPEM(caCrt)
	}
}

func (params *Params) initHttpClient() {
	params.pool = x509.NewCertPool()
	params.initCACert()
	params.Client = &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            params.pool,
			InsecureSkipVerify: params.Insecure,
		},
		MaxIdleConnsPerHost: 1024,
	}}
}

func (params *Params) Init() {
	params.initHttpClient()
}

func (params *Params) TestConn() error {
	req, err := http.NewRequest(http.MethodGet, "https://"+params.Remote, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept-Encoding", "identity")
	req.SetBasicAuth(params.Username, params.Password)

	res, err := params.Client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 401 {
		return errors.New("unable to pass the authentication on the remote server")
	}
	return nil
}
