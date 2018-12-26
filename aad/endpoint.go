package aad

import (
	"azure-blockchain-connector/aad/deviceflow"
	"golang.org/x/oauth2"
	"strings"
)

type EndpointBase string

const (
	//EndpointAuthorize            = "https://login.windows-ppe.net/<tenant id>/oauth2/authorize"
	//EndpointToken                = "https://login.windows-ppe.net/<tenant id>/oauth2/token"
	//EndpointDeviceCode           = "https://login.windows-ppe.net/<tenant id>/oauth2/devicecode"
	//EndpointRedirectNativeClient = "https://login.windows-ppe.net/common/oauth2/nativeclient"
	EndpointAuthorize            = "https://login.microsoftonline.com/<tenant id>/oauth2/authorize"
	EndpointToken                = "https://login.microsoftonline.com/<tenant id>/oauth2/token"
	EndpointDeviceCode           = "https://login.microsoftonline.com/<tenant id>/oauth2/devicecode"
	EndpointRedirectNativeClient = "https://login.microsoftonline.com/common/oauth2/nativeclient"
	TenantCommon                 = "common"
	TenantOrganizations          = "organizations"
)

func Endpoint(base EndpointBase, tenantID string) string {
	return strings.Replace(string(base), "<tenant id>", tenantID, 1)
}

func AuthCodeEndpoint(tenantID string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:  Endpoint(EndpointAuthorize, tenantID),
		TokenURL: Endpoint(EndpointToken, tenantID),
	}
}

func DeviceFlowEndpoint(tenantID string) deviceflow.Endpoint {
	return deviceflow.Endpoint{
		DeviceCodeURL: Endpoint(EndpointDeviceCode, tenantID),
		TokenURL:      Endpoint(EndpointToken, tenantID),
	}
}
