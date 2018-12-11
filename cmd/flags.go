package main

import (
	"azure-blockchain-connector/aad"
	"azure-blockchain-connector/proxy"
	"flag"
	"fmt"
	"os"
)

const (
	defaultLocalAddr   = "127.0.0.1:3100"
	defaultAuthSvcAddr = "127.0.0.1:3101"
	// Do not use oauth grant means using basic auth
)

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

	flag.StringVar(&params.Method, "method", proxy.MethodBasicAuth, "AAD: Grant type if use AAD OAuth, 'authcode' for authorization code grant or 'device' for device flow grant.")

	// basic auth
	flag.StringVar(&params.Local, "local", defaultLocalAddr, "Local address to bind to")
	flag.StringVar(&params.Remote, "remote", "", "Remote endpoint address")
	flag.StringVar(&params.CertPath, "cert", "", "(Optional) File path to root CA")
	flag.BoolVar(&params.Insecure, "insecure", false, "Skip certificate verifications")
	flag.StringVar(&params.Username, "username", "", "Basic auth: The username you want to login with")
	flag.StringVar(&params.Password, "password", "", "Basic auth: The password you want to login with")

	// AAD OAuth
	var clientID, clientSecret, scopes string
	flag.StringVar(&clientID, "client-id", "", "AAD: Client ID")
	flag.StringVar(&clientSecret, "client-secret", "", "AAD: Client Secret, required when grant type is 'authcode'")
	flag.StringVar(&scopes, "scopes", "", "AAD: Scope, should be a space-delimiter string")
	flag.StringVar(&params.AuthSvcAddr, "svc-addr", defaultAuthSvcAddr, "Should be consistent with AAD redirect config")

	var whenlogstr string
	var whatlogstr string
	var debugmode bool
	flag.StringVar(&whenlogstr, "whenlog", proxy.LogWhenOnError, "Configuration about in what cases logs should be prited. Alternatives: always, onNon200 and onError. Default: onError")
	flag.StringVar(&whatlogstr, "whatlog", proxy.LogWhatBasic, "Configuration about what information should be included in logs. Alternatives: basic and detailed. Default: basic")
	flag.BoolVar(&debugmode, "debugmode", false, "Open debug mode. It will set whenlog to always and whatlog to detailed, and original settings for whenlog and whatlog are covered.")

	flag.Parse()

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

	var p *proxy.Provider

	switch params.Method {
	case proxy.MethodBasicAuth:
		checkStr(params.Remote, params.Username, params.Password)
	case proxy.MethodOAuthAuthCode:
		checkStr(clientID, clientSecret)
		params.AuthCodeConf = aad.NewAuthCodeConfig(clientID, clientSecret, scopes)
	case proxy.MethodOAuthDeviceFlow:
		checkStr(clientID)
		params.DeviceFlowConf = aad.NewDeviceFlowConfig(clientID, scopes)
	}

	return &proxy.Proxy{Params: params, Provider: p}
}
