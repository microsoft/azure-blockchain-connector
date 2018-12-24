package aad

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/url"
)

// todo: timeout settings
// AuthCodeWebview opens a webview window with authURL, and detects the redirect URL to get auth code and state.

func newStateToken() string {
	return uuid.New().String()
}

// resolveCallback returns the code field of a query string.
func resolveCallback(query string, state string) (string, error) {

	values, err := url.ParseQuery(query)
	if err != nil {
		return "", err
	}

	errName := values.Get("error")
	if errName != "" {
		return "", errors.New("oauth2: server: " + errName)
	}

	// check state token to avoid CSRF
	if values.Get("state") != state {
		err = errors.New("oauth2: state token not the same")
		return "", err
	}

	return values.Get("code"), err
}

// fnRequestAuthCode represents the process of visiting a URL to get authorization code.
type fnRequestAuthCode func(authURL, stateToken string) (code string, err error)

// authCodeGrant returns an oauth2 token type with customizable code fetching process.
func authCodeGrant(ctx context.Context, conf *oauth2.Config, fn fnRequestAuthCode) (*oauth2.Token, error) {
	stateToken := newStateToken()
	authUrl := conf.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)

	code, err := fn(authUrl, stateToken)
	if err != nil {
		return nil, err
	}

	tok, err := conf.Exchange(ctx, code)
	return tok, err
}
