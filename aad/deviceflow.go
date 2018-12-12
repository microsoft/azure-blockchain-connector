package aad

import (
	"azure-blockchain-connector/aad/deviceflow"
	"context"
	"fmt"
)

func DeviceFlowGrant(ctx context.Context, conf *deviceflow.Config) (*deviceflow.Token, error) {

	deviceAuth, err := conf.AuthDevice(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Open:", deviceAuth.VerificationURL)
	fmt.Println("Enter:", deviceAuth.UserCode)

	return conf.Poll(ctx, deviceAuth)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("Token:", tok.AccessToken)
	//fmt.Println("Expires in:", tok.ExpiresIn)
}
