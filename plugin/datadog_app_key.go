package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	datadogAppKeyType = "datadog_app_key"
)

type datadogAppKey struct {
	AppKey   string `json:"app_key"`
	AppKeyID string `json:"app_key_id"`
}

func (b *datadogBackend) datadogAppKey() *framework.Secret {
	return &framework.Secret{
		Type: datadogAppKeyType,
		Fields: map[string]*framework.FieldSchema{
			"app_key": {
				Type:        framework.TypeString,
				Description: "datadog Application Key",
			},
		},
		Renew:  b.appKeyRenew,
		Revoke: b.appKeyRevoke,
	}
}

func (b *datadogBackend) appKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
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

func (b *datadogBackend) appKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	client, err := b.getClient(ctx, req.Storage)
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	appKeyID := ""

	appKeyIDRaw, ok := req.Secret.InternalData["app_key_id"]
	if ok {
		appKeyID, ok = appKeyIDRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid value for appKeyID in secret internal data")
		}
	}

	if err := deleteAppKey(ctx, client, appKeyID); err != nil {
		return nil, fmt.Errorf("error revoking Application Key: %w", err)
	}
	return nil, nil
}

func createAppKey(ctx context.Context, c *datadogClient, name string, scopes []string) (*datadogAppKey, error) {

	appKey, err := c.createAppKey(ctx, name, scopes)
	if err != nil {
		return nil, err
	}

	return appKey, nil
}

func deleteAppKey(ctx context.Context, c *datadogClient, appKeyID string) error {

	err := c.deleteAppKey(ctx, appKeyID)
	if err != nil {
		return err
	}

	return nil
}
