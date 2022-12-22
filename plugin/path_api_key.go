package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	apiKeyPath        = "apikey/"
	pathAPIKeyHelpSyn = `
	Generate a datadog API Token from a role.
	`
	pathAPIKeyHelpDesc = `
	This path generates a datadog API Key based on a particular 
	role.
	`
)

func pathAPIKey(b *datadogBackend) *framework.Path {
	return &framework.Path{
		Pattern: apiKeyPath + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathAPIKeyRead,
			logical.UpdateOperation: b.pathAPIKeyRead,
		},
		HelpSynopsis:    pathAPIKeyHelpSyn,
		HelpDescription: pathAPIKeyHelpDesc,
	}
}

func (b *datadogBackend) pathAPIKeyRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	roleName := d.Get("name").(string)

	roleEntry, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("error generating UUID for API key name: %w", err)
	}
	keyName := roleName + "-" + uuid

	apiKey, err := createAPIKey(ctx, client, keyName)
	if err != nil {
		return nil, fmt.Errorf("error creating datadog API key: %w", err)
	}

	resp := b.Secret(datadogAPIKeyType).Response(map[string]interface{}{
		"api_key": apiKey.APIKey,
	}, map[string]interface{}{
		"api_key_id": apiKey.APIKeyID,
	})

	if roleEntry.TTL > 0 {
		resp.Secret.TTL = roleEntry.TTL
	}

	if roleEntry.MaxTTL > 0 {
		resp.Secret.MaxTTL = roleEntry.MaxTTL
	}

	return resp, nil

}
