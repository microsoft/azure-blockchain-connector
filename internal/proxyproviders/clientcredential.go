package proxyproviders

import (
	"abc/internal/proxy"
	"context"
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
		return err
	}
	printToken(tok)

	ac.client = ac.Config.Client(ctx)
	return nil
}

func (ac *OAuthClientCredentials) Client() *http.Client {
	return ac.client
}

func (ac *OAuthClientCredentials) Modify(params *proxy.Params, req *http.Request) {
}
