package providers

import (
	"abc/internal/proxy"
	"context"
	"fmt"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
)

type OAuthClientCredentials struct {
	*clientcredentials.Config
	client *http.Client
}

func (ac *OAuthClientCredentials) RequestAccess() error {
	ctx := context.Background()

	tok, err := ac.Config.Token(ctx)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Access", tok.AccessToken)
	}
	ac.client = ac.Config.Client(ctx)
	return nil
}

func (ac *OAuthClientCredentials) Client() *http.Client {
	return ac.client
}

func (ac *OAuthClientCredentials) Modify(params *proxy.Params, req *http.Request) {
}
