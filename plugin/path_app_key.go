package plugin

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	appKeyPath        = "appkey/"
	pathAppKeyHelpSyn = `
	Generate a datadog Application Key from a role.
	`
	pathAppKeyHelpDesc = `
	This path generates a datadog Application Key based on a particular 
	role.
	`
)

func pathAppKey(b *datadogBackend) *framework.Path {
	return &framework.Path{
		Pattern: appKeyPath + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the role",
				Required:    true,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation:   b.pathAppKeyRead,
			logical.UpdateOperation: b.pathAppKeyRead,
		},
		HelpSynopsis:    pathAppKeyHelpSyn,
		HelpDescription: pathAppKeyHelpDesc,
	}
}

func (b *datadogBackend) pathAppKeyRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

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
		return nil, fmt.Errorf("error generating UUID for App key name: %w", err)
	}
	keyName := roleName + "-" + uuid

	appKey, err := createAppKey(ctx, client, keyName, roleEntry.AppKeyScopes)
	if err != nil {
		return nil, fmt.Errorf("error creating datadog application key: %w", err)
	}

	resp := b.Secret(datadogAppKeyType).Response(map[string]interface{}{
		"app_key": appKey.AppKey,
	}, map[string]interface{}{
		"app_key_id": appKey.AppKeyID,
		"role":       roleEntry.Name,
	})

	if roleEntry.TTL > 0 {
		resp.Secret.TTL = roleEntry.TTL
	}

	if roleEntry.MaxTTL > 0 {
		resp.Secret.MaxTTL = roleEntry.MaxTTL
	}

	return resp, nil
}
