package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathConfigDef             = "config"
	configStoragePath         = "config"
	pathConfigHelpSynopsis    = "Configure the datadog backend"
	pathConfigHelpDescription = `
	The Datadog secret backend requires credentials for managing
	API and App keys.

	You must provide an API and App key scoped at 
	least with the ability to create an API and 
	App key before using this secrets backend.
	`
)

type datadogConfig struct {
	APIKey string `json:"api_key"`
	AppKey string `json:"app_key"`
}

func pathConfig(b *datadogBackend) *framework.Path {

	return &framework.Path{
		Pattern: pathConfigDef,
		Fields: map[string]*framework.FieldSchema{
			"api_key": {
				Type:        framework.TypeString,
				Description: "The API Key for accessing datadog's API",
				Required:    true,
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "API Key",
					Sensitive: true,
				},
			},
			"app_key": {
				Type:        framework.TypeString,
				Description: "The Application Key scoped to admin level priveleges",
				DisplayAttrs: &framework.DisplayAttributes{
					Name:      "Application Key",
					Sensitive: true,
				},
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathConfigRead,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathConfigWrite,
			},
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathConfigDelete,
			},
		},
		ExistenceCheck:  b.PathConfigExistenceCheck,
		HelpSynopsis:    pathConfigHelpSynopsis,
		HelpDescription: pathConfigHelpDescription,
	}
}

func (b *datadogBackend) pathConfigRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	_, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Data: map[string]interface{}{},
	}, nil
}

func (b *datadogBackend) pathConfigWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	config, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	createOperation := (req.Operation == logical.CreateOperation)

	if config == nil {
		if !createOperation {
			return nil, errors.New("config not found during update operation")
		}
		config = new(datadogConfig)
	}

	if apiKey, ok := data.GetOk("api_key"); ok {
		config.APIKey = apiKey.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing API Key in configuration")
	}

	if appKey, ok := data.GetOk("app_key"); ok {
		config.AppKey = appKey.(string)
	} else if !ok && createOperation {
		return nil, fmt.Errorf("missing Application Key in configuration")
	}

	entry, err := logical.StorageEntryJSON(configStoragePath, config)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	b.reset()

	return nil, nil
}

func (b *datadogBackend) pathConfigDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	err := req.Storage.Delete(ctx, configStoragePath)

	if err == nil {
		b.reset()
	}

	return nil, err
}

func (b *datadogBackend) PathConfigExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {

	out, err := req.Storage.Get(ctx, configStoragePath)
	if err != nil {
		return false, fmt.Errorf("existence check failed: %w", err)
	}
	return out != nil, nil
}

func getConfig(ctx context.Context, s logical.Storage) (*datadogConfig, error) {
	entry, err := s.Get(ctx, configStoragePath)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	config := new(datadogConfig)
	if err := entry.DecodeJSON(&config); err != nil {
		return nil, fmt.Errorf("error reading root configuration: %w", err)
	}

	return config, nil
}
