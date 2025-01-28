package op

import (
	"context"
	"errors"

	"github.com/1password/onepassword-sdk-go"
	"github.com/charmbracelet/log"
	"github.com/gbh-tech/envi/pkg/utils"
)

type Client struct {
	Ctx    context.Context
	Client *onepassword.Client
}

func NewClient(token string) *Client {
	if token == "" {
		log.Fatal("Unauthorized. OnePassword token not present.")
	}

	ctx := context.TODO()
	client, err := onepassword.NewClient(
		ctx,
		onepassword.WithServiceAccountToken(token),
		onepassword.WithIntegrationInfo(
			"My 1Password Integration",
			"v1.0.0",
		),
	)

	if err != nil {
		log.Fatalf("Failed to create OnePassword client: %v", err)
	}

	return &Client{
		Ctx:    ctx,
		Client: client,
	}
}

func (client *Client) GetVaultName(options Options) string {
	var VaultName string

	vaults, err := client.Client.Vaults.ListAll(context.Background())
	if err != nil {
		log.Fatalf("Error fetching vaults: %s", err)
	}

	for {
		vault, err := vaults.Next()

		if vault != nil && vault.ID == options.Vault {
			VaultName = vault.Title
			break
		}

		if errors.Is(err, onepassword.ErrorIteratorDone) {
			log.Fatalf("Vault with ID '%s' not found", options.Vault)
		} else if err != nil {
			log.Fatalf("Error iterating through vaults: %s", err)
		}
	}

	return VaultName
}

func (client *Client) GenerateEnvFile(options Options) {
	vaultName := client.GetVaultName(options)

	for _, item := range options.Items {
		vaultItem, err := client.Client.Items.Get(
			client.Ctx,
			options.Vault,
			item,
		)
		if err != nil {

			log.Fatalf("failed to get vault item: %v", err)
		}

		envData := make(utils.EnvVarObject)
		for _, field := range vaultItem.Fields {
			envData[field.Title] = field.Value
		}

		for _, path := range options.Path {
			if err := utils.GenerateEnvFile(envData, path); err != nil {
				log.Fatalf("Failed to generate env file at %s: %v", path, err)
			}
			log.Infof("File generated in %s from vault: %s and item: %s using 1Password!\n", path, vaultName, vaultItem.Title)
		}
	}
}
