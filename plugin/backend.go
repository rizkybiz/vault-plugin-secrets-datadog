package plugin

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

// datadogBackend defines a struct that extends the Vault backend
// and stores the datadog API Client
type datadogBackend struct {
	*framework.Backend
	lock   sync.RWMutex
	client *datadogClient
}

// backendHelp defines the helptext for the datadog backend
const backendHelp = `
The datadog secrets backend allows for the dynamic generation of 
datadog API and App keys. After mounting this backend, credentials to 
interact with the datadog API must be configured with the /config 
endpoints.
`

// Factory returns a new datadog backend as logical.Backend
func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := newBackend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// newBackend returns a new datadogBackend and sets up the paths it will handle and
// secrets it will store
func newBackend() *datadogBackend {

	var b = datadogBackend{}
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			pathRole(&b),
			[]*framework.Path{
				pathConfig(&b),
				pathConfigRotate(&b),
				pathAPIKey(&b),
				pathAppKey(&b),
			},
		),
		Secrets: []*framework.Secret{
			b.datadogAPIKey(),
			b.datadogAppKey(),
		},
		BackendType: logical.TypeLogical,
		Invalidate:  b.invalidate,
	}

	return &b
}

// reset clears the datadog client config for a new backend to be configured
func (b *datadogBackend) reset() {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.client = nil
}

// invalidate clears an existing datadog client configuration within the backend
func (b *datadogBackend) invalidate(ctx context.Context, key string) {
	if key == "config" {
		b.reset()
	}
}

// getClient locks the datadog backend as it configures and creates a new
// datadog API client
func (b *datadogBackend) getClient(ctx context.Context, s logical.Storage) (*datadogClient, error) {
	b.lock.RLock()
	unlockFunc := b.lock.RUnlock
	defer func() { unlockFunc() }()

	if b.client != nil {
		return b.client, nil
	}

	b.lock.RUnlock()
	b.lock.Lock()
	unlockFunc = b.lock.Unlock

	config, err := getConfig(ctx, s)
	if err != nil {
		return nil, err
	}

	if config == nil {
		config = new(datadogConfig)
	}

	b.client, err = NewClient(config)
	if err != nil {
		return nil, err
	}

	return b.client, nil
}
