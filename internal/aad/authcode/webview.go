// +build linux darwin

package authcode

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
	"io"
)

// now webview is not supported in macOS and Linux

const DefaultWebviewFlag = "authcode-webview"

// AuthCodeWebview opens a window for grant operations. Once authorized, it writes the code to the writer.
func Webview(authURL string, out io.Writer) {
}

// GrantWebview use a webview window to help users to complete authorization.
func GrantWebview(ctx context.Context, conf *oauth2.Config, extraParamsSrc interface{}, flagName string) (*oauth2.Token, error) {
	return nil, errors.New("webview: not supported in macOS and Linux, please use aaddevice or aadclient method")
}
