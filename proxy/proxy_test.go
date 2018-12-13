package proxy

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var insecureClient = func() *http.Client {
	return &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},},}
}

type emptyProvider struct {
	c *http.Client
}

func (p *emptyProvider) RequestAccess() error {
	p.c = insecureClient()
	return nil
}

func (p *emptyProvider) Client(params *Params) *http.Client {
	return p.c
}

func (p *emptyProvider) Modify(params *Params, req *http.Request) {
}

func BenchmarkDirect(b *testing.B) {
	ns := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ns.Close()
	c := insecureClient()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = c.Get(ns.URL)
	}
}

func BenchmarkProxy(b *testing.B) {
	proxy := &Proxy{
		Params: &Params{
			Whenlog: LogWhenOnError,
		},
		Provider: &emptyProvider{},
	}
	_ = proxy.Provider.RequestAccess()

	ns := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer ns.Close()

	ps := httptest.NewServer(proxy)
	defer ps.Close()

	u, _ := url.Parse(ns.URL)
	proxy.Params.Remote = u.Host
	proxy.Params.Local = ps.URL

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = http.Get(ps.URL)
	}
}
