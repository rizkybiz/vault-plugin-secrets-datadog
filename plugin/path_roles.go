package plugin

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	pathRoleDef             = "roles/"
	pathRoleHelpSynopsis    = "Manages the Vault role for generating Datadog API and Application Keys."
	pathRoleHelpDescription = `
	This path allows you to read and write roles used to generate Datadog API and Application Keys.
	You can configure scopes associated with Application Keys by providing a list of scopes with the 
	input data.
	`
	pathRoleListHelpSynopsis    = "List the existing roles in datadog backend"
	pathRoleListHelpDescription = "Roles will be listed by the role name."
)

var (
	appKeyScopes = []string{
		"user_access_invite",
		"user_access_manage",
		"user_access_read",
		"usage_read",
		"incident_read",
		"incident_settings_write",
		"incident_write",
		"security_monitoring_filters_read",
		"security_monitoring_filters_write",
		"security_monitoring_rules_read",
		"security_monitoring_rules_write",
		"security_monitoring_signals_read",
		"dashboards_public_share",
		"dashboards_read",
		"dashboards_write",
		"events_read",
		"metrics_read",
		"timeseries_query",
		"monitors_downtime",
		"monitors_read",
		"monitors_write",
		"synthetics_global_variable_read",
		"synthetics_global_variable_write",
		"synthetics_private_location_read",
		"synthetics_private_location_write",
		"synthetics_read",
		"synthetics_write",
	}
)

// datadogRoleEntry defines the data associated with
// a Vault role for interoperating with the datadog
// api
type datadogRoleEntry struct {
	Name         string        `json:"name"`
	AppKeyScopes []string      `json:"app_key_scopes"`
	TTL          time.Duration `json:"ttl"`
	MaxTTL       time.Duration `json:"max_ttl"`
}

// pathRole defines the framework.Path for datadog roles
func pathRole(b *datadogBackend) []*framework.Path {

	return []*framework.Path{
		{
			Pattern: pathRoleDef + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeLowerCaseString,
					Description: "Required. Name of the role",
					Required:    true,
				},
				"app_key_scopes": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Optional. List of datadog permissions scopes to be applied to the application key.",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Optional. Default lease time for generated credentials. If not set or set to 0, system default will be used.",
				},
				"max_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Optional. Maximum lease time for role. If not set or set to 0, system default will be used.",
				},
			},
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: b.pathRolesRead,
				},
				logical.CreateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: b.pathRolesWrite,
				},
				logical.DeleteOperation: &framework.PathOperation{
					Callback: b.pathRolesDelete,
				},
			},
			HelpSynopsis:    pathRoleHelpSynopsis,
			HelpDescription: pathRoleHelpDescription,
		},
		{
			Pattern: pathRoleDef + "?$",
			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ListOperation: &framework.PathOperation{
					Callback: b.pathRolesList,
				},
			},
			HelpSynopsis:    pathRoleListHelpSynopsis,
			HelpDescription: pathRoleListHelpDescription,
		},
	}
}

// pathRolesList lists the datadog roleEntries
func (b *datadogBackend) pathRolesList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	entries, err := req.Storage.List(ctx, pathRoleDef)
	if err != nil {
		return nil, err
	}

	return logical.ListResponse(entries), nil
}

// pathRolesRead returns a specifc datadog roleEntry
func (b *datadogBackend) pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	entry, err := b.getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: entry.toResponseData(),
	}, nil
}

// pathRolesWrite creates or updates a datadog roleEntry
func (b *datadogBackend) pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	name := d.Get("name").(string)
	if name == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	roleEntry, err := b.getRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if roleEntry == nil {
		roleEntry = &datadogRoleEntry{}
	}

	roleEntry.Name = name

	createOperation := (req.Operation == logical.CreateOperation)

	if scopes, ok := d.GetOk("app_key_scopes"); ok {
		roleEntry.AppKeyScopes = scopes.([]string)
		// check validity of provided scopes
		for _, scope := range roleEntry.AppKeyScopes {
			valid := contains(appKeyScopes, scope)
			if !valid {
				return nil, fmt.Errorf("provided scope %s is not a valid datadog application key scope", scope)
			}
		}
	} else if createOperation {
		roleEntry.AppKeyScopes = d.Get("app_key_scopes").([]string)
		// check validity of provided scopes
		for _, scope := range roleEntry.AppKeyScopes {
			valid := contains(appKeyScopes, scope)
			if !valid {
				return nil, fmt.Errorf("provided scope %s is not a valid datadog application key scope", scope)
			}
		}
	}

	if ttlRaw, ok := d.GetOk("ttl"); ok {
		roleEntry.TTL = time.Duration(ttlRaw.(int)) * time.Second
	} else if createOperation {
		roleEntry.TTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	if maxTTLRaw, ok := d.GetOk("max_ttl"); ok {
		roleEntry.MaxTTL = time.Duration(maxTTLRaw.(int)) * time.Second
	} else if createOperation {
		roleEntry.MaxTTL = time.Duration(d.Get("max_ttl").(int)) * time.Second
	}

	if roleEntry.MaxTTL != 0 && roleEntry.TTL > roleEntry.MaxTTL {
		return logical.ErrorResponse("ttl cannot be greater than max_ttl"), nil
	}

	if err := setRole(ctx, req.Storage, name, roleEntry); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRolesDelete deletes a datadog roleEntry
func (b *datadogBackend) pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	err := req.Storage.Delete(ctx, pathRoleDef+d.Get("name").(string))
	if err != nil {
		return nil, fmt.Errorf("error deleting datadog role: %w", err)
	}

	return nil, nil
}

// getRole gets the role from the Vault storage API
func (b *datadogBackend) getRole(ctx context.Context, s logical.Storage, name string) (*datadogRoleEntry, error) {

	if name == "" {
		return nil, fmt.Errorf("missing role name")
	}

	entry, err := s.Get(ctx, pathRoleDef+name)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var role datadogRoleEntry
	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// setRole sets the role into the Vault storage API
func setRole(ctx context.Context, s logical.Storage, name string, roleEntry *datadogRoleEntry) error {

	entry, err := logical.StorageEntryJSON(pathRoleDef+name, roleEntry)
	if err != nil {
		return err
	}

	if entry == nil {
		return fmt.Errorf("failed to create storage entry for role")
	}

	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

func validateScopes(scopes []string) error {

	return nil
}

// toResponseData returns response data for a datadog role entry
func (r *datadogRoleEntry) toResponseData() map[string]interface{} {

	return map[string]interface{}{
		"app_key_scopes": r.AppKeyScopes,
		"ttl":            r.TTL.Seconds(),
		"max_ttl":        r.MaxTTL.Seconds(),
	}

}
