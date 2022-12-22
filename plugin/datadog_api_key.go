package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	datadogAPIKeyType = "datadog_api_key"
)

type datadogAPIKey struct {
	APIKey   string `json:"api_key"`
	APIKeyID string `json:"api_key_id"`
}

func (b *datadogBackend) datadogAPIKey() *framework.Secret {

	return &framework.Secret{
		Type: datadogAPIKeyType,
		Fields: map[string]*framework.FieldSchema{
			"api_key": {
				Type:        framework.TypeString,
				Description: "datadog API Key",
			},
		},
		Renew:  b.apiKeyRenew,
		Revoke: b.apiKeyRevoke,
	}
}

func (b *datadogBackend) apiKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	roleRaw, ok := req.Secret.InternalData["role"]
	if !ok {
		return nil, fmt.Errorf("secret is missing role internal data")
	}

	// get the role entry
	role := roleRaw.(string)
	roleEntry, err := b.getRole(ctx, req.Storage, role)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %w", err)
	}

	if roleEntry == nil {
		return nil, errors.New("error retrieving role: role is nil")
	}

	resp := &logical.Response{Secret: req.Secret}

	if roleEntry.TTL > 0 {
		resp.Secret.TTL = roleEntry.TTL
	}
	if roleEntry.MaxTTL > 0 {
		resp.Secret.MaxTTL = roleEntry.MaxTTL
	}

	return resp, nil
}

func (b *datadogBackend) apiKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	apiKeyID := ""
	// we passed the api_key_id into internalData when we created it
	apiKeyIDRaw, ok := req.Secret.InternalData["api_key_id"]
	if ok {
		apiKeyID, ok = apiKeyIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for apiKeyID in secret internal data")
		}
	}

	if err := deleteAPIKey(ctx, client, apiKeyID); err != nil {
		return nil, fmt.Errorf("error revoking API Key: %w", err)
	}
	return nil, nil
}

func createAPIKey(ctx context.Context, c *datadogClient, name string) (*datadogAPIKey, error) {

	apiKey, err := c.createAPIKey(ctx, name)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func deleteAPIKey(ctx context.Context, c *datadogClient, apiKeyID string) error {

	err := c.deleteAPIKey(ctx, apiKeyID)
	if err != nil {
		return err
	}

	return nil
}
