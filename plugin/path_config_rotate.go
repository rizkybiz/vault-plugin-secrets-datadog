package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathConfigRotateHelpSyn = `
	Rotate the datadog API and App keys.
	`
	pathConfigRotateHelpDesc = `
	This will rotate the datadog API and App keys that are 
	used to interact with the datadog platform.
	`
)

func pathConfigRotate(b *datadogBackend) *framework.Path {

	return &framework.Path{
		Pattern: pathConfigDef + "/rotate",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRotateRead,
				Summary:  "Rotate datadog API and App Keys",
			},
		},
		HelpSynopsis:    pathConfigRotateHelpSyn,
		HelpDescription: pathConfigRotateHelpDesc,
	}
}

func (b *datadogBackend) pathConfigRotateRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting config: %w", err)
	}

	if config == nil {
		return logical.ErrorResponse("configuration not set"), nil
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	oldAPIKeyID := config.APIKeyID
	oldAppKeyID := config.AppKeyID

	uuid, _ := uuid.GenerateUUID()
	newAPIKey, err := createAPIKey(ctx, client, "vault-config-"+uuid)
	if err != nil {
		return nil, fmt.Errorf("error rotating API key: %w", err)
	}
	newAppKey, err := createAppKey(ctx, client, "vault-config-"+uuid, []string{})
	if err != nil {
		return nil, fmt.Errorf("error rotating App key: %w", err)
	}
	config.APIKey = newAPIKey.APIKey
	config.AppKey = newAppKey.AppKey
	config.APIKeyID = newAPIKey.APIKeyID
	config.AppKeyID = newAppKey.AppKeyID
	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	b.reset()

	client, err = b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	err = deleteAPIKey(ctx, client, oldAPIKeyID)
	if err != nil {
		return nil, err
	}
	err = deleteAppKey(ctx, client, oldAppKeyID)
	if err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"api_key_id": config.APIKeyID,
			"app_key_id": config.AppKeyID,
		},
	}, nil
}
