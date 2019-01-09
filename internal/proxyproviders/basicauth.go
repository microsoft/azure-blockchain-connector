package proxyproviders

import (
	"abc/internal/proxy"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

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

func (ba *BasicAuth) Client() *http.Client {
	return ba.client
}

func (ba *BasicAuth) Modify(params *proxy.Params, req *http.Request) {
	req.SetBasicAuth(ba.Username, ba.Password)
}
