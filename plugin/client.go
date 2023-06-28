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

	// as of v2.14 of the DD API, a call to datadogv2.ApplicationKeyCreateRequest
	// requires the Scopes attributes to be a datadog.NullableList[string]
	ns := datadog.NewNullableList[string](&scopes)

	body := datadogV2.ApplicationKeyCreateRequest{
		Data: datadogV2.ApplicationKeyCreateData{
			Attributes: datadogV2.ApplicationKeyCreateAttributes{
				Name:   name,
				Scopes: *ns,
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
