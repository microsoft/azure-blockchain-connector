package main

import (
	"azure-blockchain-connector/proxy"
	"flag"
	"fmt"
	"os"
)

func NewProxyFromFlags() *proxy.Proxy {
	var params = &proxy.Params{}

	flag.StringVar(&params.Local, "local", defaultLocalAddr, "Local address to bind to")
	flag.StringVar(&params.Remote, "remote", "", "Remote endpoint address")

	flag.StringVar(&params.CertPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.Insecure, "insecure", false, "Skip certificate verifications")

	// basic auth
	flag.StringVar(&params.Username, "username", "", "Basic auth: The username you want to login with")
	flag.StringVar(&params.Password, "password", "", "Basic auth: The password you want to login with")

	// AAD OAuth
	var oauthType, clientID, clientSecret, scopes string
	flag.StringVar(&oauthType, "aad", proxy.OAuthGrantTypeNone, "AAD: Grant type if use AAD OAuth, 'authcode' for authorization code grant or 'device' for device flow grant.")
	flag.StringVar(&clientID, "client-id", "", "AAD: Client ID")
	flag.StringVar(&clientSecret, "client-secret", "", "AAD: Client Secret, required when grant type is 'authcode'")
	flag.StringVar(&scopes, "scopes", "", "AAD: Scope")
	flag.StringVar(&params.AuthSvcAddr, "oauth service address", defaultAuthSvcAddr, "Should be consistent with AAD redirect config")

	var whenlogstr string
	var whatlogstr string
	var debugmode bool
	flag.StringVar(&whenlogstr, "whenlog", proxy.LogWhenOnError, "Configuration about in what cases logs should be prited. Alternatives: always, onNon200 and onError. Default: onError")
	flag.StringVar(&whatlogstr, "whatlog", proxy.LogWhatBasic, "Configuration about what information should be included in logs. Alternatives: basic and detailed. Default: basic")
	flag.BoolVar(&debugmode, "debugmode", false, "Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.")

	flag.Parse()

	if params.Remote == "" || params.Username == "" || params.Password == "" {
		flag.Usage()
		os.Exit(-1)
	}

	params.SetOAuthConfig(oauthType, clientID, clientSecret, scopes)

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

	return &proxy.Proxy{Params: params}
}
