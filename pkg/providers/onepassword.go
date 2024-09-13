package providers

import (
	"github.com/1Password/connect-sdk-go/connect"
)

func NewOnePasswordProvider() connect.Client {
	client, err := connect.NewClientFromEnvironment()

	if err != nil {
		return nil
	}

	return client
}
