package authcode

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"net/url"
	"reflect"
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

func options(v interface{}, keys []string) []oauth2.AuthCodeOption {
	var opts []oauth2.AuthCodeOption
	elm := reflect.ValueOf(v).Elem()
	typ := elm.Type()

	for i := 0; i < elm.NumField(); i++ {
		f := elm.Field(i)
		if f.Type() != reflect.TypeOf("") {
			continue
		}
		k, v := typ.Field(i).Tag.Get("key"), elm.Field(i).String()
		if !stringSliceContains(keys, k) || v == "" {
			continue
		}
		opts = append(opts, oauth2.SetAuthURLParam(k, v))
	}
	return opts
}

type OptionsSource interface {
	OptionsNameList() ([]string, []string)
}

func authCodeOptions(src OptionsSource) ([]oauth2.AuthCodeOption, []oauth2.AuthCodeOption) {
	l, l2 := src.OptionsNameList()
	// always requesting a refresh token is not bad now
	return append(options(src, l), oauth2.AccessTypeOffline), options(src, l2)
}

// authCodeGrant returns an oauth2 token type with customizable code fetching process.
func authCodeGrant(ctx context.Context, conf *oauth2.Config, src OptionsSource, fn fnRequestAuthCode) (*oauth2.Token, error) {
	stateToken := newStateToken()
	opts, opts2 := authCodeOptions(src)

	authURL := conf.AuthCodeURL(stateToken, opts...)

	code, err := fn(authURL, stateToken)
	if err != nil {
		return nil, err
	}

	opts2 = append(opts2, oauth2.SetAuthURLParam("client_id", conf.ClientID))
	tok, err := conf.Exchange(ctx, code, opts2...)

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
	ResponseMode string `key:"response_mode"`
	Resource     string `key:"resource"`
	Prompt       string `key:"prompt"`
	LoginHint    string `key:"login_hint"`
	DomainHint   string `key:"domain_hint"`
	//CodeChallengeMethod string
	//CodeChallenge       string
	//CodeVerifier        string
}

func (c *Config) OptionsNameList() ([]string, []string) {
	return []string{"response_mode", "resource", "prompt", "login_hint", "domain_hint"}, []string{"resource"}
}

func (c *Config) Server(ctx context.Context, svcAddr string) (*oauth2.Token, error) {
	return GrantServer(ctx, c.Config, c, svcAddr)
}

func (c *Config) Webview(ctx context.Context, flagName string) (*oauth2.Token, error) {
	return GrantWebview(ctx, c.Config, c, flagName)
}
