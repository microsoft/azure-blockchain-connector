package aad

import (
	"azure-blockchain-connector/aad/oauth2/devicecode"
	"context"
	"fmt"
)

func DeviceFlowGrant(ctx context.Context, conf *devicecode.Config) (*devicecode.Token, error) {

	deviceAuth, err := conf.AuthDevice(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Open:", deviceAuth.VerificationURL)
	fmt.Println("Enter:", deviceAuth.UserCode)

	return conf.Poll(ctx, deviceAuth)
}
