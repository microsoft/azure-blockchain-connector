package authcode

import (
	"bufio"
	"context"
	"fmt"
	"github.com/zserge/webview"
	"golang.org/x/oauth2"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	styleSuccPage = `body {font-family: 'Segoe UI', SegoeUI, 'Segoe WP', Tahoma, Arial, sans-serif; font-size: 16px; font-weight: 400;color: white; background-color: #007FFF; user-select: none; }
							section {position: fixed; top: 0; right: 0; bottom: 0; left: 0; display: flex; flex-flow: column nowrap; align-items: center; justify-content: center;}
							h2 {display: block; margin-top: -7vh; font-size: 24px; font-weight: 600; line-height: 1.2;}`
	scriptShowSuccPage = `document.body.innerHTML = '<section><h2>Authorized Success!</h2><div>This window will close in <span id="cnt">5</span> seconds...</div></section>'`
	scriptCountdown    = `var cnt=document.querySelector("#cnt");if(cnt)var timer=setInterval(function(){var a=parseInt(cnt.innerHTML);a-=1,cnt.innerHTML=""+a,0>=a&&clearInterval(timer)},1e3);`
)

// AuthCodeWebview opens a window for grant operations. Once authorized, it writes the code to the writer.
func Webview(authURL string, out io.Writer) {
	complete := false

	href := stringChange("")
	w := webview.New(webview.Settings{
		Title:     "Request Access",
		Width:     800,
		Height:    600,
		URL:       authURL,
		Resizable: true,
		ExternalInvokeCallback: func(w webview.WebView, data string) {
			if !href.Changed(data) {
				return
			}
			u, _ := url.Parse(data)

			// detect condition: if the querystring includes a code field
			if u.Query().Get("code") == "" && u.Query().Get("error") == "" {
				return
			}

			complete = true
			_, _ = fmt.Fprintln(out, u.RawQuery)

			w.InjectCSS(styleSuccPage)
			_ = w.Eval(scriptShowSuccPage)
			_ = w.Eval(scriptCountdown)

			go func() {
				time.Sleep(6 * time.Second)
				w.Dispatch(func() {
					w.Terminate()
				})
			}()
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

// requestAuthCodeWebview executes a command to fetch the auth code from given authURL and state.
func requestAuthCodeWebview(authURL string, state string, flagName string) (code string, err error) {
	pth, err := os.Executable()
	if err != nil {
		return
	}
	arg := strings.Join([]string{"-", flagName, "=", authURL}, "")

	cmd := exec.Command(pth, arg)
	out, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	err = cmd.Start()
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(out)
	scanner.Scan()
	query, err := scanner.Text(), scanner.Err()

	return resolveCallback(query, state)
}

// GrantWebview use a webview window to help users to complete authorization.
func GrantWebview(ctx context.Context, conf *oauth2.Config, src OptionsSource, flagName string) (*oauth2.Token, error) {

	return authCodeGrant(ctx, conf, src, func(authURL, stateToken string) (s string, e error) {
		return requestAuthCodeWebview(authURL, stateToken, flagName)
	})
}
