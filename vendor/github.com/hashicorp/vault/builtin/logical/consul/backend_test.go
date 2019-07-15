package consul

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/helper/testhelpers/consul"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

func TestBackend_Config_Access(t *testing.T) {
	t.Run("config_access", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "1.3.0")
		})
		t.Run("1.4.0-rc", func(t *testing.T) {
			t.Parallel()
			testBackendConfigAccess(t, "1.4.0-rc1")
		})
	})
}

func testBackendConfigAccess(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, version)
	defer cleanup()

	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address": connData["address"].(string),
		"scheme":  "http",
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_Renew_Revoke(t *testing.T) {
	t.Run("renew_revoke", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendRenewRevoke(t, "1.3.0")
		})
		t.Run("1.4.0-rc", func(t *testing.T) {
			t.Parallel()
			t.Run("legacy", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke(t, "1.4.0-rc1")
			})

			testBackendRenewRevoke14(t, "1.4.0-rc1")
		})
	})
}

func testBackendRenewRevoke(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, version)
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policy": base64.StdEncoding.EncodeToString([]byte(testPolicy)),
		"lease":  "6h",
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	generatedSecret := resp.Secret
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token string `mapstructure:"token"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated token: %s", d.Token)

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Verifying that the generated token works...")
	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Verifying that the generated token does not work...")
	_, err = client.KV().Put(&consulapi.KVPair{
		Key:   "foo",
		Value: []byte("bar"),
	}, nil)
	if err == nil {
		t.Fatal("expected error")
	}

}

func testBackendRenewRevoke14(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, version)
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policies": []string{"test"},
		"lease":    "6h",
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	generatedSecret := resp.Secret
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated token: %s with accessor %s", d.Token, d.Accessor)

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Verifying that the generated token works...")
	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)

	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
	}

	t.Log("Verifying that the generated token does not exist...")
	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackend_LocalToken(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, "1.4.0-rc1")
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policies": []string{"test"},
		"ttl":      "6h",
		"local":    false,
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test_local"
	req.Data = map[string]interface{}{
		"policies": []string{"test"},
		"ttl":      "6h",
		"local":    true,
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
		Local    bool   `mapstructure:"local"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated token: %s with accessor %s", d.Token, d.Accessor)

	if d.Local {
		t.Fatalf("requested global token, got local one")
	}

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Verifying that the generated token works...")
	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test_local"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("Generated token: %s with accessor %s", d.Token, d.Accessor)

	if !d.Local {
		t.Fatalf("requested local token, got global one")
	}

	// Build a client and verify that the credentials work
	consulapiConfig = consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err = consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Verifying that the generated token works...")
	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}
}

func TestBackend_Management(t *testing.T) {
	t.Run("management", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendManagement(t, "1.3.0")
		})
		t.Run("1.4.0-rc", func(t *testing.T) {
			t.Parallel()
			testBackendManagement(t, "1.4.0-rc1")
		})
	})
}

func testBackendManagement(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, version)
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWriteManagementPolicy(t, "test", ""),
			testAccStepReadManagementToken(t, "test", connData),
		},
	})
}

func TestBackend_Basic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		t.Run("pre-1.4.0", func(t *testing.T) {
			t.Parallel()
			testBackendBasic(t, "1.3.0")
		})
		t.Run("1.4.0-rc", func(t *testing.T) {
			t.Parallel()
			t.Run("legacy", func(t *testing.T) {
				t.Parallel()
				testBackendRenewRevoke(t, "1.4.0-rc1")
			})

			testBackendBasic(t, "1.4.0-rc1")
		})
	})
}

func testBackendBasic(t *testing.T, version string) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := consul.PrepareTestContainer(t, version)
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepConfig(t, connData),
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepReadToken(t, "test", connData),
		},
	})
}

func TestBackend_crud(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, ""),
			testAccStepWritePolicy(t, "test2", testPolicy, ""),
			testAccStepWritePolicy(t, "test3", testPolicy, ""),
			testAccStepReadPolicy(t, "test", testPolicy, 0),
			testAccStepListPolicy(t, []string{"test", "test2", "test3"}),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func TestBackend_role_lease(t *testing.T) {
	b, _ := Factory(context.Background(), logical.TestBackendConfig())
	logicaltest.Test(t, logicaltest.TestCase{
		LogicalBackend: b,
		Steps: []logicaltest.TestStep{
			testAccStepWritePolicy(t, "test", testPolicy, "6h"),
			testAccStepReadPolicy(t, "test", testPolicy, 6*time.Hour),
			testAccStepDeletePolicy(t, "test"),
		},
	})
}

func testAccStepConfig(
	t *testing.T, config map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      config,
	}
}

func testAccStepReadToken(
	t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Token string `mapstructure:"token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated token: %s", d.Token)

			// Build a client and verify that the credentials work
			config := consulapi.DefaultConfig()
			config.Address = conf["address"].(string)
			config.Token = d.Token
			client, err := consulapi.NewClient(config)
			if err != nil {
				return err
			}

			log.Printf("[WARN] Verifying that the generated token works...")
			_, err = client.KV().Put(&consulapi.KVPair{
				Key:   "foo",
				Value: []byte("bar"),
			}, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func testAccStepReadManagementToken(
	t *testing.T, name string, conf map[string]interface{}) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "creds/" + name,
		Check: func(resp *logical.Response) error {
			var d struct {
				Token string `mapstructure:"token"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				return err
			}
			log.Printf("[WARN] Generated token: %s", d.Token)

			// Build a client and verify that the credentials work
			config := consulapi.DefaultConfig()
			config.Address = conf["address"].(string)
			config.Token = d.Token
			client, err := consulapi.NewClient(config)
			if err != nil {
				return err
			}

			log.Printf("[WARN] Verifying that the generated token works...")
			_, _, err = client.ACL().Create(&consulapi.ACLEntry{
				Type: "management",
				Name: "test2",
			}, nil)
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func testAccStepWritePolicy(t *testing.T, name string, policy string, lease string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"policy": base64.StdEncoding.EncodeToString([]byte(policy)),
			"lease":  lease,
		},
	}
}

func testAccStepWriteManagementPolicy(t *testing.T, name string, lease string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.UpdateOperation,
		Path:      "roles/" + name,
		Data: map[string]interface{}{
			"token_type": "management",
			"lease":      lease,
		},
	}
}

func testAccStepReadPolicy(t *testing.T, name string, policy string, lease time.Duration) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ReadOperation,
		Path:      "roles/" + name,
		Check: func(resp *logical.Response) error {
			policyRaw := resp.Data["policy"].(string)
			out, err := base64.StdEncoding.DecodeString(policyRaw)
			if err != nil {
				return err
			}
			if string(out) != policy {
				return fmt.Errorf("mismatch: %s %s", out, policy)
			}

			l := resp.Data["lease"].(int64)
			if lease != time.Second*time.Duration(l) {
				return fmt.Errorf("mismatch: %v %v", l, lease)
			}
			return nil
		},
	}
}

func testAccStepListPolicy(t *testing.T, names []string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.ListOperation,
		Path:      "roles/",
		Check: func(resp *logical.Response) error {
			respKeys := resp.Data["keys"].([]string)
			if !reflect.DeepEqual(respKeys, names) {
				return fmt.Errorf("mismatch: %#v %#v", respKeys, names)
			}
			return nil
		},
	}
}

func testAccStepDeletePolicy(t *testing.T, name string) logicaltest.TestStep {
	return logicaltest.TestStep{
		Operation: logical.DeleteOperation,
		Path:      "roles/" + name,
	}
}

const testPolicy = `
key "" {
	policy = "write"
}
`
