package devicecode

import (
	"abc/internal/oauth2dc"
	"context"
	"fmt"
)

type Config struct {
	*oauth2dc.Config
	Resource string `key:"resource"`
}

func (c *Config) Grant(ctx context.Context) (*oauth2dc.Token, error) {
	deviceAuth, err := c.Config.AuthDevice(ctx, oauth2dc.SetAuthURLParam("resource", c.Resource))
	if err != nil {
		return nil, err
	}
	fmt.Println("Open:", deviceAuth.VerificationURI)
	fmt.Println("Enter:", deviceAuth.UserCode)
	return c.Config.Poll(ctx, deviceAuth)
}
