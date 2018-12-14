package deviceflow

import (
	"context"
	"encoding/json"
	"errors"
	"golang.org/x/net/context/ctxhttp"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Endpoint contains OAuth2 device flow's device authorization url
// and device access token url.
type Endpoint struct {
	DeviceCodeURL string
	TokenURL      string
}

// Config describes the information required by the device flow.
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// Scope specifies optional requested permissions.
	Scopes []string

	Endpoint Endpoint
}

type DeviceAuth struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURL string `json:"verification_url"`
	//VerificationURLComplete string `json:"verification_uri_complete"`
	ExpiresIn int `json:"expires_in,string"`
	Interval  int `json:"interval,string"`
	// Some providers has a message field for display.
	Message string `json:"message"`
}

// AuthDevice returns a URL and a code for user to authenticate
// using another device. It also provides device code for Poll().

// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.3
func (c *Config) AuthDevice(ctx context.Context) (*DeviceAuth, error) {
	v := url.Values{
		"client_id": {c.ClientID},
		"scope":     {strings.Join(c.Scopes, " ")},
	}

	req, err := http.NewRequest("POST", c.Endpoint.DeviceCodeURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, err
	}

	if code := res.StatusCode; code < 200 || code > 299 {
		return nil, errors.New("oauth2: deviceflow.AuthDevice: " + res.Status)
	}

	var da = &DeviceAuth{}
	err = json.NewDecoder(res.Body).Decode(&da)
	if err != nil {
		return nil, err
	}
	return da, nil
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Error        string `json:"error"`
}

const (
	ErrAuthorizationPending = "authorization_pending"
	ErrSlowDown             = "slow_down"
	ErrAccessDenied         = "access_denied"
	ErrExpiredToken         = "expired_token"
)

// Poll does a polling to get token.
func (c *Config) Poll(ctx context.Context, da *DeviceAuth) (*Token, error) {
	v := url.Values{
		"client_id": {c.ClientID},
		// Providers may use "device_code" for short, in RFC it must be set to this value.
		// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.4
		//"grant_type":  {"device_code"},
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"scope":       {strings.Join(c.Scopes, " ")},
		"device_code": {da.DeviceCode},
	}

	// If no interval was provided, the client MUST use a reasonable default polling interval.
	// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.5
	interval := da.Interval
	if interval == 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		req, err := http.NewRequest("POST", c.Endpoint.TokenURL, strings.NewReader(v.Encode()))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res, err := ctxhttp.Do(ctx, nil, req)
		if err != nil {
			return nil, err
		}

		var tok = &Token{}
		err = json.NewDecoder(res.Body).Decode(tok)

		if res.StatusCode == http.StatusOK {
			return tok, err
		}

		switch tok.Error {
		case ErrAccessDenied, ErrExpiredToken:
			return tok, errors.New("oauth2: " + tok.Error)
		case ErrSlowDown:
			interval += 5
			fallthrough
		case ErrAuthorizationPending:
		}

	}
}
