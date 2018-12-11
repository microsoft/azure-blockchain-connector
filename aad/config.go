package aad

import (
	"azure-blockchain-connector/aad/deviceflow"
	"golang.org/x/oauth2"
	"strings"
)

// Authorization Code Grant endpoint and config

var AuthCodeEndpoint = oauth2.Endpoint{
	AuthURL:  "https://login.microsoftonline.com/organizations/oauth2/v2.0/authorize",
	TokenURL: "https://login.microsoftonline.com/organizations/oauth2/v2.0/token",
}

func NewAuthCodeConfig(clientID, clientSecret, scopes string) *oauth2.Config {
	return &oauth2.Config{
		Endpoint:     AuthCodeEndpoint,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(scopes, " "),
	}
}

// Device Flow Grant endpoint and config

var DeviceFlowEndpoint = deviceflow.Endpoint{
	DeviceCodeURL: "https://login.microsoftonline.com/organizations/oauth2/v2.0/devicecode",
	TokenURL:      "https://login.microsoftonline.com/organizations/oauth2/v2.0/token",
}

func NewDeviceFlowConfig(clientID, scopes string) *deviceflow.Config {
	return &deviceflow.Config{
		Endpoint: DeviceFlowEndpoint,
		ClientID: clientID,
		Scopes:   strings.Split(scopes, " "),
	}
}
