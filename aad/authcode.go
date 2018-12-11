package aad

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

var stateToken = uuid.New().String()

func AuthCodeGrant(conf *oauth2.Config, svcAddr string) {
	ctx := context.Background()

	authUrl := conf.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)
	fmt.Println("Authorize:", "http://"+svcAddr+"/authorization")

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/authorization", http.StatusSeeOther))
	mux.Handle("/authorization", http.RedirectHandler(authUrl, http.StatusSeeOther))
	mux.HandleFunc("/authorization/callback", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		code := queries.Get("code")

		// check state token to avoid CSRF
		state := queries.Get("state")
		if state != stateToken {
			return
		}

		// exchange authorization_code for access_token
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}

		tokJson, err := json.Marshal(tok)
		fmt.Println(string(tokJson))

		//client := conf.Client(ctx, tok)
	})

	log.Fatal(http.ListenAndServe(svcAddr, mux))
}
