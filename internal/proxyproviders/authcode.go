package proxyproviders

import (
	"abc/internal/aad/authcode"
	"abc/internal/proxy"
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

type OAuthAuthCode struct {
	*authcode.Config
	UseWebview bool
	SvcAddr    string
	ArgName    string
	client     *http.Client
}

func (ac *OAuthAuthCode) RequestAccess() error {
	ctx := context.Background()

	var tok *oauth2.Token
	var err error

	if ac.UseWebview {
		tok, err = ac.Config.Webview(ctx, ac.ArgName)
	} else {
		tok, err = ac.Config.Server(ctx, ac.SvcAddr)
	}

	if err != nil {
		return err
	}
	printToken(tok)

	ac.client = ac.Config.Client(ctx, tok)

	return nil
}

func (ac *OAuthAuthCode) Client() *http.Client {
	return ac.client
}

func (ac *OAuthAuthCode) Modify(params *proxy.Params, req *http.Request) {
}
