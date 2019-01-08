package authcode

import (
	"abc/internal/util"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/oauth2"
	"net/url"
)

// todo: timeout settings
// AuthCodeWebview opens a webview window with authURL, and detects the redirect URL to get auth code and state.

func newStateToken() string {
	data := make([]byte, 64)
	if _, err := rand.Read(data); err != nil {
		return "_state"
	}
	return base64.StdEncoding.EncodeToString(data)
}

// resolveCallback returns the code field of a query string.
func resolveCallback(query string, state string) (string, error) {
	if query == "" {
		return "", errors.New("oauth2: authorization code not appear")
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return "", err
	}

	// check state token to avoid CSRF
	if values.Get("state") != state {
		err = errors.New("oauth2: state token not the same")
		return "", err
	}

	errName := values.Get("error")
	if errName != "" {
		return "", errors.New("oauth2: server: " + errName + "\n" + values.Get("error_description"))
	}

	return values.Get("code"), err
}

// fnRequestAuthCode represents the process of visiting a URL to get authorization code.
type fnRequestAuthCode func(authURL, stateToken string) (code string, err error)

// authCodeGrant returns an oauth2 token type with customizable code fetching process.
func authCodeGrant(ctx context.Context, conf *oauth2.Config, extraParamsSrc interface{}, fn fnRequestAuthCode) (*oauth2.Token, error) {
	stateToken := newStateToken()

	authURL := conf.AuthCodeURL(stateToken, util.FieldsToOAuthParams(extraParamsSrc, "auth_code")...)
	code, err := fn(authURL, stateToken)
	if err != nil {
		return nil, err
	}

	opts := append(util.FieldsToOAuthParams(extraParamsSrc, "exchange"), oauth2.SetAuthURLParam("client_id", conf.ClientID))
	tok, err := conf.Exchange(ctx, code, opts...)

	return tok, err
}

const (
	PromptLogin          = "login"
	PromptSelectAccount  = "select_account"
	PromptConsent        = "consent"
	PromptAdminConsent   = "admin_consent"
	ResponseModeQuery    = "query"
	ResponseModeFragment = "fragment"
	ResponseModeFormPost = "form_post"
)

type Config struct {
	*oauth2.Config
	ResponseMode string `auth_code:"response_mode"`
	Resource     string `auth_code:"resource" exchange:"resource"`
	Prompt       string `auth_code:"prompt"`
	LoginHint    string `auth_code:"login_hint"`
	DomainHint   string `auth_code:"domain_hint"`
	//CodeChallengeMethod string
	//CodeChallenge       string
	//CodeVerifier        string
}

func (c *Config) Server(ctx context.Context, svcAddr string) (*oauth2.Token, error) {
	return GrantServer(ctx, c.Config, c, svcAddr)
}

func (c *Config) Webview(ctx context.Context, flagName string) (*oauth2.Token, error) {
	return GrantWebview(ctx, c.Config, c, flagName)
}
