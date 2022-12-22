package plugin

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	APIKey = "c2c72f3f3e831c89f841f1f98b8977f2"
	AppKey = "b8fbb773f987b9b06cbced638d7dfc68cb3c7940"
)

func TestConfig(t *testing.T) {
	b, reqStorage := getTestBackend(t)

	t.Run("Test Configuration", func(t *testing.T) {
		// test the config create functionality
		err := testConfigCreate(t, b, reqStorage, map[string]interface{}{
			"api_key": APIKey,
			"app_key": AppKey,
		})
		assert.NoError(t, err)

		// test the config read functionality
		err = testConfigRead(t, b, reqStorage, map[string]interface{}{})
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
		Path:      configStoragePath,
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
		Path:      configStoragePath,
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
		Path:      configStoragePath,
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
		Path:      configStoragePath,
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
