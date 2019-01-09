package devicecode

import (
	"abc/internal/util"
	"context"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"net/url"
	"strings"
	"time"
)

type Endpoint struct {
	DeviceAuthURL string
	TokenURL      string
}

type Config struct {
	*oauth2.Config
	Resource string `auth_device:"resource"`
	Endpoint
}

// AuthDevice returns a device auth struct which contains a device code
// and authorization information provided for users to enter on another device.
func (c *Config) AuthDevice(ctx context.Context, params ...util.StringKVP) (*DeviceAuth, error) {
	v := url.Values{
		"client_id": {c.ClientID},
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	for _, kv := range params {
		v.Set(kv.K, kv.V)
	}
	return retrieveDeviceAuth(ctx, c, v)
}

// Poll does a polling to exchange an device code for a token.
func (c *Config) Poll(ctx context.Context, da *DeviceAuth, params ...util.StringKVP) (*oauth2.Token, error) {
	v := url.Values{
		"client_id": {c.ClientID},
		// Providers may use "device_code" for short, in RFC it must be set to this value.
		// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.4
		//"grant_type":  {"device_code"},
		"grant_type":  {"urn:ietf:params:oauth:grant-type:device_code"},
		"device_code": {da.DeviceCode},
		"code":        {da.DeviceCode},
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	for _, kv := range params {
		v.Set(kv.K, kv.V)
	}

	// If no interval was provided, the client MUST use a reasonable default polling interval.
	// See https://tools.ietf.org/html/draft-ietf-oauth-device-flow-07#section-3.5
	interval := da.Interval
	if interval == 0 {
		interval = 5
	}

	for {
		time.Sleep(time.Duration(interval) * time.Second)

		tok, err := retrieveToken(ctx, c, v)
		if err == nil {
			return tok, nil
		}
		errTyp := parseError(err)
		switch errTyp {
		case errAccessDenied, errExpiredToken:
			return tok, errors.New("oauth2: " + errTyp)
		case errSlowDown:
			interval += 5
			fallthrough
		case errAuthorizationPending:
		}
	}
}

func (c *Config) Grant(ctx context.Context) (*oauth2.Token, error) {
	deviceAuth, err := c.AuthDevice(ctx, util.FieldsToStringKVPs(c, "auth_device")...)
	if err != nil {
		return nil, err
	}
	fmt.Println("Open:", deviceAuth.VerificationURI)
	fmt.Println("Enter:", deviceAuth.UserCode)
	return c.Poll(ctx, deviceAuth)
}
