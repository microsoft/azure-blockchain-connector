package aad

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
)

var stateToken = uuid.New().String()

func AuthCodeGrant(ctx context.Context, conf *oauth2.Config, svcAddr string) (*oauth2.Token, error) {

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
		queries := r.URL.Query()
		code := queries.Get("code")

		// check state token to avoid CSRF
		state := queries.Get("state")
		if state != stateToken {
			return
		}

		// exchange authorization_code for access_token
		tok, err = conf.Exchange(ctx, code)
		if err != nil {
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
