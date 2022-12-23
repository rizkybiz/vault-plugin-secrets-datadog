package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

type datadogClient struct {
	*datadog.APIClient
}

func NewClient(config *datadogConfig) (*datadogClient, error) {

	if config == nil {
		return nil, errors.New("client configuration was nil")
	}

	// create datadog APIClient
	conf := datadog.NewConfiguration()
	if config.APIKey == "" {
		return nil, errors.New("datadog API key was not provided")
	}
	conf.AddDefaultHeader("DD-API-KEY", config.APIKey)
	if config.AppKey == "" {
		return nil, errors.New("datadog aaplication key was not provided")
	}
	conf.AddDefaultHeader("DD-APPLICATION-KEY", config.AppKey)
	c := datadog.NewAPIClient(conf)

	return &datadogClient{c}, nil
}

func (c *datadogClient) listAPIKeys(ctx context.Context) ([]*datadogAPIKey, error) {

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	resp, _, err := api.ListAPIKeys(ctx)
	if err != nil {
		return nil, err
	}

	partialKeys := resp.Data
	var apiKeys []*datadogAPIKey
	for _, partKey := range partialKeys {
		resp, _, err := api.GetAPIKey(ctx, *partKey.Id)
		if err != nil {
			return nil, err
		}
		apiKeys = append(apiKeys, &datadogAPIKey{
			APIKey:   *resp.Data.Attributes.Key,
			APIKeyID: *partKey.Id,
		})
	}
	return apiKeys, nil
}

func (c *datadogClient) listAppKeys(ctx context.Context) ([]*datadogAppKey, error) {

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	resp, _, err := api.ListApplicationKeys(ctx)
	if err != nil {
		return nil, err
	}

	partialKeys := resp.Data
	var appKeys []*datadogAppKey
	for _, partKey := range partialKeys {
		resp, _, err := api.GetApplicationKey(ctx, *partKey.Id)
		if err != nil {
			return nil, err
		}
		appKeys = append(appKeys, &datadogAppKey{
			AppKey:   *resp.Data.Attributes.Key,
			AppKeyID: *partKey.Id,
		})
	}
	return appKeys, nil
}

func (c *datadogClient) createAPIKey(ctx context.Context, apiKeyName string) (*datadogAPIKey, error) {

	body := datadogV2.APIKeyCreateRequest{
		Data: datadogV2.APIKeyCreateData{
			Attributes: *datadogV2.NewAPIKeyCreateAttributes(apiKeyName),
			Type:       datadogV2.APIKEYSTYPE_API_KEYS,
		},
	}

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	ddresp, _, err := api.CreateAPIKey(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("error creating datadog API key; %w", err)
	}
	respData := ddresp.GetData()

	return &datadogAPIKey{
		APIKeyID: *respData.Id,
		APIKey:   *respData.Attributes.Key,
	}, err
}

func (c *datadogClient) deleteAPIKey(ctx context.Context, apiKeyID string) error {

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	_, err := api.DeleteAPIKey(ctx, apiKeyID)
	if err != nil {
		return fmt.Errorf("error deleting datadog API key: %w", err)
	}
	return nil
}

func (c *datadogClient) createAppKey(ctx context.Context, name string, scopes []string) (*datadogAppKey, error) {

	body := datadogV2.ApplicationKeyCreateRequest{
		Data: datadogV2.ApplicationKeyCreateData{
			Attributes: datadogV2.ApplicationKeyCreateAttributes{
				Name:   name,
				Scopes: scopes,
			},
			Type: datadogV2.APPLICATIONKEYSTYPE_APPLICATION_KEYS,
		},
	}

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	ddresp, _, err := api.CreateCurrentUserApplicationKey(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("error creating datadog application key: %w", err)
	}

	respData := ddresp.GetData()

	return &datadogAppKey{
		AppKeyID: *respData.Id,
		AppKey:   *respData.Attributes.Key,
	}, nil
}

func (c *datadogClient) deleteAppKey(ctx context.Context, appKeyID string) error {

	api := datadogV2.NewKeyManagementApi(c.APIClient)

	_, err := api.DeleteApplicationKey(ctx, appKeyID)
	if err != nil {
		return fmt.Errorf("error deleting datadog application key: %w", err)
	}

	return nil
}
