package aad

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zserge/webview"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

var stateToken = uuid.New().String()

// AuthCodeGrant prints the authorization url to stdio and the user need to click the url to perform a grant.
// This method listen to a port to receive the callback of the code and the state token from the server.
// Then, it will terminate the server and returns received token values.
// The browser window will also be closed(via window.close()) immediately after getting the information required.
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
			err = errors.New("oauth2: state token not the same")
			close(complete)
			return
		}

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

// todo: modify the condition for getting the code and state, current method is dangerous
// todo: window flashing past for an logged-in users
// todo: timeout settings
// AuthCodeWebview opens a webview window with authURL, and detects the redirect URL to get auth code and state.
func AuthCodeWebview(authURL string) {
	complete := false
	href := stringChange("")
	w := webview.New(webview.Settings{
		Title:     "Authorization",
		Width:     800,
		Height:    600,
		URL:       authURL,
		Resizable: true,
		ExternalInvokeCallback: func(w webview.WebView, data string) {
			if href.Changed(data) {
				complete = true
				fmt.Println(data)
				w.Terminate()
			}
		},
	})

	defer w.Exit()

	go func(fin *bool) {
		for !complete {
			w.Dispatch(func() {
				_ = w.Eval(`window.external.invoke(window.location.href)`)
			})
			time.Sleep(time.Second)
		}
	}(&complete)
	w.Run()
}

// todo: name need to be changed
// todo: move to proxy package
// frame will exec command to fetch the auth code from given authURL and state.
func frame(authURL string, stateToken string, flagName string) (string, error) {
	pth, err := os.Executable()
	if err != nil {
		return "", err
	}
	arg := strings.Join([]string{"-", flagName, "=", authURL}, "")
	out, err := exec.Command(pth, arg).Output()

	callbackURL := string(out)
	u, err := url.Parse(callbackURL)
	if err != nil {
		return "", err
	}
	values := u.Query()

	code := values.Get("code")

	// check state token to avoid CSRF
	state := values.Get("state")
	if state != stateToken {
		err = errors.New("oauth2: state token not the same")
		return "", err
	}

	return code, err
}

// todo: move to proxy package
// AuthCodeGrantWithFrame use a webview window to help users to complete authorization.
func AuthCodeGrantWithFrame(ctx context.Context, conf *oauth2.Config, flagName string) (*oauth2.Token, error) {
	authUrl := conf.AuthCodeURL(stateToken, oauth2.AccessTypeOffline)

	code, err := frame(authUrl, stateToken, flagName)
	if err != nil {
		return nil, err
	}

	// exchange authorization_code for access_token
	tok, err := conf.Exchange(ctx, code)
	return tok, err
}
