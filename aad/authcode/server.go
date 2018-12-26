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
	htmlSuccPage         = `<!DOCTYPE html><html lang="en"><head> <meta charset="UTF-8"> <title></title> <style>body{font-family: 'Segoe UI', SegoeUI, 'Segoe WP', Tahoma, Arial, sans-serif; font-size: 16px; font-weight: 400; color: white; background-color: #007FFF; user-select: none;}section{position: fixed; top: 0; right: 0; bottom: 0; left: 0; display: flex; align-items: center; justify-content: center; flex-flow: column nowrap;}h2{margin-top: -10vh; font-size: 24px; font-weight: 600; line-height: 1.2;}</style></head><body><section> <h2>Authorization succeeded</h2> <div>You may now close this window.</div></section></body></html>`
)

func prefixHTTP(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	return addr
}

func CallbackPath(addr string) string {
	return prefixHTTP(addr) + PathAuthCodeCallback
}

// GrantServer prints the authorization url to stdio and the user need to click the url to perform a grant.
// This method listen to a port to receive the callback of the code and the state token from the server.
// Then, it will terminate the server and returns received token values.
// The browser window will also be closed(via window.close()) immediately after getting the information required.
func GrantServer(ctx context.Context, conf *oauth2.Config, src OptionsSource, svcAddr string) (*oauth2.Token, error) {

	return authCodeGrant(ctx, conf, src, func(authURL, stateToken string) (code string, err error) {
		mux := http.NewServeMux()
		srv := &http.Server{Addr: svcAddr, Handler: mux}
		mux.Handle(PathAuthCodeAuth, http.RedirectHandler(authURL, http.StatusSeeOther))

		complete := make(chan struct{})
		mux.HandleFunc(PathAuthCodeCallback, func(w http.ResponseWriter, r *http.Request) {

			code, err = resolveCallback(r.URL.RawQuery, stateToken)
			if err != nil {
				close(complete)
				return
			}
			_, _ = io.WriteString(w, htmlSuccPage)

			close(complete)
		})

		srvErr := make(chan error, 1)
		go func() {
			srvErr <- srv.ListenAndServe()
		}()
		fmt.Println("Authorize:", prefixHTTP(svcAddr)+PathAuthCodeAuth)

		select {
		case <-complete:
			go func() {
				err := srv.Shutdown(ctx)
				if err != nil {
					log.Println(err)
				}
			}()
		case err = <-srvErr:
		}
		return
	})
}
