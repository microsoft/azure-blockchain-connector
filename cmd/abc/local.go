package main

import (
	"abc/internal/aad/authcode"
	"abc/internal/proxy/providers"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"runtime"
	"strings"
)

//noinspection ALL
const (
	EnableProxy = false
)

//noinspection ALL
func init() {
	providers.EnablePrintToken = true
	//EnableProxy = true

	if EnableProxy {
		u, _ := url.Parse("http://localhost:8888")
		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(u)}
	}
}

// Account: auxtm434@live.com
// Password: #Bugsfor$1
func testSampleAuxtm434() {
	if runtime.GOOS != "windows" {
		return
	}

	conf := &providers.OAuthAuthCode{
		Config: &authcode.Config{
			Config: &oauth2.Config{
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://login.windows-ppe.net/auxteststageauto.ccsctp.net/oauth2/authorize",
					TokenURL: "https://login.windows-ppe.net/auxteststageauto.ccsctp.net/oauth2/token",
				},
				ClientID:    "a8196997-9cc1-4d8a-8966-ed763e15c7e1",
				RedirectURL: "https://login.windows-ppe.net/common/oauth2/nativeclient",
			},
			Resource: "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca",
			Prompt:   authcode.PromptSelectAccount,
		},
		UseWebview: true,
		ArgName:    authcode.DefaultWebviewFlag,
	}

	err := conf.RequestAccess()
	if err != nil {
		log.Fatalln(err)
	}

	c := conf.Client()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		fmt.Println(a)

		if !strings.Contains(a, "Bearer") {
			log.Fatalln("bearer token not found")
		}
	}))
	defer ts.Close()

	_, err = c.Get(ts.URL)
	if err != nil {
		log.Fatalln(err)
	}
}
