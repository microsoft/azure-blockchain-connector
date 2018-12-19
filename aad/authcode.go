package aad

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
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

	// check state token to avoid CSRF
	if values.Get("state") != state {
		err = errors.New("oauth2: state token not the same")
		return "", err
	}

	return values.Get("code"), err
}

// AuthCodeGrantWithServer prints the authorization url to stdio and the user need to click the url to perform a grant.
// This method listen to a port to receive the callback of the code and the state token from the server.
// Then, it will terminate the server and returns received token values.
// The browser window will also be closed(via window.close()) immediately after getting the information required.
func AuthCodeGrantServer(ctx context.Context, conf *oauth2.Config, svcAddr string) (*oauth2.Token, error) {

	stateToken := newStateToken()
	authUrl := conf.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)
	fmt.Println("Authorize:", "http://"+svcAddr+"/authorization")

	mux := http.NewServeMux()
	srv := &http.Server{Addr: svcAddr, Handler: mux}

	mux.Handle("/", http.RedirectHandler("/authorization", http.StatusSeeOther))
	mux.Handle("/authorization", http.RedirectHandler(authUrl, http.StatusSeeOther))

	complete := make(chan struct{})
	var tok *oauth2.Token
	var err error
	mux.HandleFunc("/authorization/callback", func(w http.ResponseWriter, r *http.Request) {
		code, err := resolveCallback(r.URL.RawQuery, stateToken)

		// exchange authorization_code for access_token
		tok, err = conf.Exchange(ctx, code)
		if err != nil {
			close(complete)
			return
		}

		_, _ = io.WriteString(w, `<script>window.close()</script>`)

		close(complete)
	})

	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	select {
	case <-complete:
		go func() {
			err := srv.Shutdown(ctx)
			if err != nil {
				log.Println(err)
			}
		}()
		return tok, err
	case err := <-srvErr:
		return nil, err
	}
}
