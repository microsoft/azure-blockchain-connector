package devicecode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/oauth2"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in,string"`
	raw          map[string]interface{}
}

func (tr *tokenResp) expiry() (t time.Time) {
	if v := tr.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

func retrieveToken(ctx context.Context, c *Config, v url.Values) (*oauth2.Token, error) {

	req, err := http.NewRequest("POST", c.Endpoint.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := res.StatusCode; code < 200 || code > 299 {
		return nil, &oauth2.RetrieveError{
			Response: res,
			Body:     body,
		}
	}

	var tr tokenResp
	if err = json.Unmarshal(body, &tr); err != nil {
		return nil, err
	}

	// for extra fields
	tr.raw = make(map[string]interface{})
	_ = json.Unmarshal(body, &tr.raw)

	token := &oauth2.Token{
		TokenType:    tr.TokenType,
		AccessToken:  tr.AccessToken,
		RefreshToken: tr.RefreshToken,
		Expiry:       tr.expiry(),
	}

	// if token is not provided, use previous value
	if token.RefreshToken == "" {
		token.RefreshToken = v.Get("refresh_token")
	}
	if token.AccessToken == "" {
		return token, errors.New("oauth2: server response missing access_token")
	}

	return token, nil
}
