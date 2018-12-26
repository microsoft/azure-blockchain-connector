package authcode

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	PathAuthCodeAuth     = "/_click_to_auth"
	PathAuthCodeCallback = "/_callback"
)

func CallbackPath(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	return addr + PathAuthCodeCallback
}

// GrantServer prints the authorization url to stdio and the user need to click the url to perform a grant.
// This method listen to a port to receive the callback of the code and the state token from the server.
// Then, it will terminate the server and returns received token values.
// The browser window will also be closed(via window.close()) immediately after getting the information required.
func GrantServer(ctx context.Context, conf *oauth2.Config, src OptionsSource, svcAddr string) (*oauth2.Token, error) {
	stateToken := newStateToken()
	authURL := conf.AuthCodeURL(stateToken)
	fmt.Println("Authorize:", "http://"+svcAddr+PathAuthCodeAuth)

	fmt.Println(authURL)

	mux := http.NewServeMux()
	srv := &http.Server{Addr: svcAddr, Handler: mux}

	mux.Handle(PathAuthCodeAuth, http.RedirectHandler(authURL, http.StatusSeeOther))

	complete := make(chan struct{})
	var tok *oauth2.Token
	var err error
	mux.HandleFunc(PathAuthCodeCallback, func(w http.ResponseWriter, r *http.Request) {
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
