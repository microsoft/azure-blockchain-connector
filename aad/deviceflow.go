package aad

import (
	"azure-blockchain-connector/aad/deviceflow"
	"context"
	"fmt"
)

func DeviceFlowGrant(conf *deviceflow.Config) {
	var ctx = context.Background()

	deviceAuth, err := conf.AuthDevice(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Open:", deviceAuth.VerificationURL)
	fmt.Println("Enter:", deviceAuth.UserCode)

	tok, err := conf.Poll(ctx, deviceAuth)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Token:", tok.AccessToken)
	fmt.Println("Expires in:", tok.ExpiresIn)
}
