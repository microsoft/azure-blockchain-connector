package main

import (
	"abc/internal/aad"
	"abc/internal/aad/authcode"
	"abc/internal/aad/devicecode"
	"abc/internal/oauth2dc"
	"abc/internal/proxy"
	"abc/internal/proxyproviders"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"net/url"
	"os"
	"runtime"
	"strings"
)

// checkStr checks if the str is "", then print flag.Usage to ask the user.
// Keep the same exit code -1 with the former implementation.
func checkStr(namesStr string, ss ...string) {
	names := strings.Split(namesStr, " ")
	for i, s := range ss {
		if s == "" {
			if i < len(names) {
				fmt.Printf("Error: param '%s' is required for the current method.\n", names[i])
			}
			flag.Usage()
			os.Exit(-1)
		}
	}
}

func newProxyFromFlags() *proxy.Proxy {
	var params = &proxy.Params{}

	flag.StringVar(&params.Method, "method", methodBasicAuth, "Authentication method. Basic auth (basic), authorization code (aadauthcode), client credentials (aadclient) and device flow(aaddevice)")
	flag.StringVar(&params.Local, "local", defaultLocalAddr, "Local address to bind to")
	flag.StringVar(&params.Remote, "remote", "", "Remote endpoint address")

	flag.StringVar(&params.CertPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.Insecure, "insecure", false, "(Optional) Skip certificate verifications")

	// basic auth
	var username, password string
	flag.StringVar(&username, "username", "", "Basic auth: username")
	flag.StringVar(&password, "password", "", "Basic auth: password")

	// AAD OAuth
	var (
		clientID, tenantID, clientSecret string
		useWebview                       bool
		authSvcAddr                      string
	)
	flag.StringVar(&clientID, "client-id", "", "OAuth: application (client) ID")
	flag.StringVar(&tenantID, "tenant-id", "", "OAuth: directory (tenant) ID")
	flag.StringVar(&clientSecret, "client-secret", "", "OAuth: client secret")
	flag.BoolVar(&useWebview, "webview", true, "OAuth: open a webview o to receive callbacks, applicable for Windows/macOS")
	flag.StringVar(&authSvcAddr, "authcode-addr", defaultLocalAddr, "OAuth: local address to receive callbacks")

	var whenlogstr string
	var whatlogstr string
	var debugmode bool
	flag.StringVar(&whenlogstr, "whenlog", proxy.LogWhenOnError, "Configuration about in what cases logs should be prited. Alternatives: always, onNon200 and onError")
	flag.StringVar(&whatlogstr, "whatlog", proxy.LogWhatBasic, "Configuration about what information should be included in logs. Alternatives: basic and detailed")
	flag.BoolVar(&debugmode, "debugmode", false, "Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.")

	flag.Parse()

	switch params.Method {
	case methodBasicAuth, methodOAuthAuthCode, methodOAuthClientCredentials, methodOAuthDeviceFlow:
	default:
		fmt.Println("Unexpected method value. Expected: basic, aadauthcode, aadclient, aaddevice")
		os.Exit(-1)
	}

	switch whenlogstr {
	case proxy.LogWhenOnError, proxy.LogWhenOnNon200, proxy.LogWhenAlways:
	default:
		fmt.Println("Unexpected whenlog value. Expected: always, onNon200 or onError")
		os.Exit(-1)
	}

	switch whatlogstr {
	case proxy.LogWhatBasic, proxy.LogWhatDetailed:
	default:
		fmt.Println("Unexpected whatlog value. Expected: basic or detailed")
		os.Exit(-1)
	}

	if debugmode {
		params.Whenlog = proxy.LogWhenAlways
		params.Whatlog = proxy.LogWhatDetailed
		proxyproviders.EnablePrintToken = true
	}

	// hard code scopes
	var scopes = []string{""}
	// In Azure AD v1, the scope field is ignored
	//scopes = []string{"offline_access", "scope value here"}
	if params.Method == methodOAuthClientCredentials {
		// See https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow
		// this method should not provide a refresh token
		scopes = []string{"https://graph.microsoft.com/.default"}
	}

	var redirectURL = authcode.CallbackPath(authSvcAddr)
	// hard code redirect URL settings for different OS webviews
	// "urn:ietf:wg:oauth:2.0:oob": webviews do not support the protocol
	// webkit(macOS): visit "nativeclient" start a download automatically
	if useWebview {
		switch runtime.GOOS {
		case "windows":
			redirectURL = aad.EndpointRedirectNativeClient
		case "darwin", "linux":
			fallthrough
		default:
			useWebview = false
		}
	}

	checkStr("local remote", params.Local, params.Remote)

	p := (func() proxy.Provider {
		switch params.Method {
		case methodOAuthAuthCode:
			checkStr("tenant-id", tenantID)
			return &proxyproviders.OAuthAuthCode{
				Config: &authcode.Config{
					Config: &oauth2.Config{
						Endpoint:     aad.AuthCodeEndpoint(tenantID),
						ClientID:     hcAuthcodeClientId,
						ClientSecret: clientSecret,
						Scopes:       scopes,
						RedirectURL:  redirectURL,
					},
					Resource: hcResource,
					Prompt:   authcode.PromptSelectAccount,
				},
				UseWebview: useWebview,
				SvcAddr:    authSvcAddr,
				ArgName:    flagAuthCodeWebview,
			}
		case methodOAuthDeviceFlow:
			checkStr("tenant-id", tenantID)
			return &proxyproviders.OAuthDeviceCode{
				Config: &devicecode.Config{
					Config: &oauth2dc.Config{
						Endpoint: aad.DeviceCodeEndpoint(tenantID),
						ClientID: hcAuthcodeClientId,
						Scopes:   scopes,
					},
					Resource: hcResource,
				},
			}
		case methodOAuthClientCredentials:
			checkStr("client-id client-secret", clientID, clientSecret)
			return &proxyproviders.OAuthClientCredentials{
				Config: &clientcredentials.Config{
					ClientID:     clientID,
					ClientSecret: clientSecret,
					TokenURL:     aad.Endpoint(aad.EndpointToken, tenantID),
					Scopes:       scopes,
					EndpointParams: url.Values{
						"resource": {hcResource},
					},
				},
			}
		case methodBasicAuth:
			fallthrough
		default:
			checkStr("username password", username, password)
			return &proxyproviders.BasicAuth{
				Remote:   params.Remote,
				Username: username,
				Password: password,
			}
		}
	})()

	return &proxy.Proxy{
		Params:   params,
		Provider: p,
	}
}
