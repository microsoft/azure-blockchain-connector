package main

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/aad/deviceflow"
	"azure-blockchain-connector/proxy"
	"azure-blockchain-connector/proxy/providers"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"os"
	"runtime"
)

const (
	defaultLocalAddr   = "127.0.0.1:3100"
	defaultAuthSvcAddr = "127.0.0.1:3101"
	// Do not use oauth grant means using basic auth
)

// cli argument for getting auth code with webview, used internally
const flagAuthCodeWebview = "authcode-webview"

// checkStr checks if the str is "", then print flag.Usage to ask the user.
// Keep the same exit code -1 with the former implementation.
func checkStr(ss ...string) {
	for _, s := range ss {
		if s == "" {
			flag.Usage()
			os.Exit(-1)
		}
	}
}

func newProxyFromFlags() *proxy.Proxy {
	var params = &proxy.Params{}

	flag.StringVar(&params.Method, "method", proxy.MethodBasicAuth, "'basic' for basic auth, 'authcode' for authorization code grant, 'client' for client credentials grant, 'device' for device flow grant. Default: basic")
	flag.StringVar(&params.Local, "local", defaultLocalAddr, "Local address to bind to")
	flag.StringVar(&params.Remote, "remote", "", "Remote endpoint address")

	flag.StringVar(&params.CertPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.Insecure, "insecure", false, "Skip certificate verifications")

	// basic auth
	var username, password string
	flag.StringVar(&username, "username", "", "Basic auth: The username you want to login with")
	flag.StringVar(&password, "password", "", "Basic auth: The password you want to login with")

	// AAD OAuth
	var clientID, tenantID, clientSecret, authSvcAddr string
	//var scopesString string
	flag.StringVar(&clientID, "client-id", "", "AAD: Application (client) ID")
	flag.StringVar(&tenantID, "tenant-id", "", "AAD: Directory (tenant) ID")
	flag.StringVar(&clientSecret, "client-secret", "", "AAD: Client Secret, required when grant type is 'authcode'")
	//flag.StringVar(&scopesString, "scopes", "", "AAD: Scope, should be a space-delimiter string")
	flag.StringVar(&authSvcAddr, "svc-addr", defaultAuthSvcAddr, "Should be consistent with AAD redirect config")

	var whenlogstr string
	var whatlogstr string
	var debugmode bool
	flag.StringVar(&whenlogstr, "whenlog", proxy.LogWhenOnError, "Configuration about in what cases logs should be prited. Alternatives: always, onNon200 and onError. Default: onError")
	flag.StringVar(&whatlogstr, "whatlog", proxy.LogWhatBasic, "Configuration about what information should be included in logs. Alternatives: basic and detailed. Default: basic")
	flag.BoolVar(&debugmode, "debugmode", false, "Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.")

	flag.Parse()

	switch params.Method {
	case proxy.MethodBasicAuth, proxy.MethodOAuthAuthCode, proxy.MethodOAuthDeviceFlow, proxy.MethodOAuthClientCredentials:
	default:
		fmt.Println("Unexpected method value. Expected: basic, authcode, device")
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
	}

	// hard code scopes
	// Azure: one scope must be supplied
	// "offline_access" is used to request a refresh token
	var scopes = []string{"offline_access", "api://285286f5-b97b-4b45-ba35-92a74f35756a/basic"}
	if params.Method == proxy.MethodOAuthClientCredentials {
		// See https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-oauth2-client-creds-grant-flow
		// this method should not provide a refresh token
		scopes = []string{"https://graph.microsoft.com/.default"}
	}

	// hard code redirect URL settings for different OS type
	var redirectURL string
	switch runtime.GOOS {
	case "windows":
		// mshtml will throw an error for "urn:ietf:wg:oauth:2.0:oob"
		redirectURL = "https://login.microsoftonline.com/common/oauth2/nativeclient"
	case "darwin":
		// macOS webview may start a download for "https://login.microsoftonline.com/common/oauth2/nativeclient"
		// which cannot be fixed now.
		// todo: macOS redirectURL: check if works
		redirectURL = "urn:ietf:wg:oauth:2.0:oob"
	case "linux":
		fallthrough
	default:
	}

	checkStr(params.Local, params.Remote)

	p := (func() proxy.Provider {
		switch params.Method {
		case proxy.MethodOAuthAuthCode:
			checkStr(clientID, tenantID)
			return &providers.OAuthAuthCode{
				Config: &oauth2.Config{
					Endpoint:     aad.AuthCodeEndpoint(tenantID),
					ClientID:     clientID,
					ClientSecret: clientSecret,
					Scopes:       scopes,
					RedirectURL:  redirectURL,
				},
				SvcAddr: authSvcAddr,
				ArgName: flagAuthCodeWebview,
			}
		case proxy.MethodOAuthClientCredentials:
			checkStr(clientID, clientSecret)
			return &providers.OAuthClientCredentials{
				Config: &clientcredentials.Config{
					ClientID:       clientID,
					ClientSecret:   clientSecret,
					TokenURL:       aad.Endpoint(aad.EndpointToken, tenantID),
					Scopes:         scopes,
					EndpointParams: nil,
				},
			}
		case proxy.MethodOAuthDeviceFlow:
			checkStr(clientID, tenantID)
			return &providers.OAuthDeviceFlow{
				Config: &deviceflow.Config{
					Endpoint: aad.DeviceFlowEndpoint(tenantID),
					ClientID: clientID,
					Scopes:   scopes,
				},
			}
		case proxy.MethodBasicAuth:
			fallthrough
		default:
			checkStr(params.Remote, username, password)
			return &providers.BasicAuth{
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
