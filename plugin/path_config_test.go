package plugin

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	APIKey   = "c2c72f3f3e831c89f841f1f98b8977f2"
	AppKey   = "b8fbb773f987b9b06cbced638d7dfc68cb3c7940"
	APIKeyID = "1e962ce6-b12a-4a87-bbb2-07fe5986334c"
	AppKeyID = "1e962ce6-b12a-4a87-bbb2-07fe5986334c"
)

func TestConfig(t *testing.T) {
	b, reqStorage := getTestBackend(t)

	t.Run("Test Configuration", func(t *testing.T) {
		// test the config create functionality
		err := testConfigCreate(t, b, reqStorage, map[string]interface{}{
			"api_key":    APIKey,
			"api_key_id": APIKeyID,
			"app_key":    AppKey,
			"app_key_id": AppKeyID,
		})
		assert.NoError(t, err)

		// test the config read functionality
		err = testConfigRead(t, b, reqStorage, map[string]interface{}{
			"api_key_id": "1e962ce6-b12a-4a87-bbb2-07fe5986334c",
			"app_key_id": "1e962ce6-b12a-4a87-bbb2-07fe5986334c",
		})
		assert.NoError(t, err)

		// test the config update functionality
		err = testConfigUpdate(t, b, reqStorage, map[string]interface{}{
			"api_key": APIKey,
			"app_key": "r8fbb773f987b9b06cbced638d7dfc68cb3c7940",
		})
		assert.NoError(t, err)

		// test the config deletion functionality
		err = testConfigDelete(t, b, reqStorage)
		assert.NoError(t, err)
	})
}

func testConfigCreate(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.CreateOperation,
		Path:      pathConfigDef,
		Data:      d,
		Storage:   s,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.IsError() {
		return resp.Error()
	}
	return nil
}

func testConfigRead(t *testing.T, b logical.Backend, s logical.Storage, expected map[string]interface{}) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.ReadOperation,
		Path:      pathConfigDef,
		Storage:   s,
	})
	if err != nil {
		return err
	}
	if resp == nil && expected == nil {
		return nil
	}
	if len(expected) != len(resp.Data) {
		return fmt.Errorf("read data mismatch (expected %d values, got %d)", len(expected), len(resp.Data))
	}
	return nil
}

func testConfigUpdate(t *testing.T, b logical.Backend, s logical.Storage, d map[string]interface{}) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      pathConfigDef,
		Data:      d,
		Storage:   s,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.IsError() {
		return resp.Error()
	}
	return nil
}

func testConfigDelete(t *testing.T, b logical.Backend, s logical.Storage) error {
	resp, err := b.HandleRequest(context.Background(), &logical.Request{
		Operation: logical.DeleteOperation,
		Path:      pathConfigDef,
		Storage:   s,
	})
	if err != nil {
		return err
	}
	if resp != nil && resp.IsError() {
		return resp.Error()
	}
	return nil
}
