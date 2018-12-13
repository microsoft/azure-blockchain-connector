package aad

import (
	"azure-blockchain-connector/aad/deviceflow"
	"golang.org/x/oauth2"
	"strings"
)

func NewAuthCodeConfig(clientID, tenantID, clientSecret, scopes string) *oauth2.Config {
	return &oauth2.Config{
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/authorize",
			TokenURL: "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token",
		},
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       strings.Split(scopes, " "),
	}
}

func NewDeviceFlowConfig(clientID, tenantID, scopes string) *deviceflow.Config {
	return &deviceflow.Config{
		Endpoint: deviceflow.Endpoint{
			DeviceCodeURL: "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/devicecode",
			TokenURL:      "https://login.microsoftonline.com/" + tenantID + "/oauth2/v2.0/token",
		},
		ClientID: clientID,
		Scopes:   strings.Split(scopes, " "),
	}
}
