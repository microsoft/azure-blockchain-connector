package main

import "abc/internal/aad/authcode"

const (
	// hard code settings
	hcAuthcodeClientId = "a8196997-9cc1-4d8a-8966-ed763e15c7e1"
	hcResource         = "5838b1ed-6c81-4c2f-8ca1-693600b4e6ca"

	defaultLocalAddr = "localhost:3100"

	methodBasicAuth              = "basic"
	methodOAuthAuthCode          = "aadauthcode"
	methodOAuthClientCredentials = "aadclient"
	methodOAuthDeviceFlow        = "aaddevice"

	// flagAuthCodeWebview is used to see if the executable is launching as a webview
	flagAuthCodeWebview = authcode.DefaultWebviewFlag
)
