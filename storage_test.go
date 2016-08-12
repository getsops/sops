package sops

import (
	"fmt"
	"go.mozilla.org/sops/json"
	"go.mozilla.org/sops/yaml"
	goyaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDecryptSimpleYAML(t *testing.T) {
	in := `foo: "ENC[AES256_GCM,data:3mDI,iv:2OL363jDglmOa+k6qCIh5RGUm+isJNrqP4umqOEb+1s=,tag:bJiC+QHsSQPJFXCX94sRGQ==,type:str]"`
	key := strings.Repeat("f", 32)
	expected := "foo"
	store := yaml.YAMLStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	if store.Data["foo"] != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", store.Data["foo"], expected)
	}
}

func TestDecryptNestedYaml(t *testing.T) {
	in := "foo:\n  - bar: \"ENC[AES256_GCM,data:tFAp,iv:5G6/F6LaoYgIV5je0TBvqLoyyo6IZK7rCjA7T5nuy1k=,tag:4M/t1ZN5ZRWBiMrYdRpBhg==,type:str]\""
	key := strings.Repeat("f", 32)
	expected := "foo"
	store := yaml.YAMLStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	foo := store.Data["foo"].([]interface{})[0].(map[interface{}]interface{})
	if foo["bar"] != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", foo["bar"], expected)
	}
}

func TestYamlRoundtrip(t *testing.T) {
	key := strings.Repeat("f", 32)
	f := func(tree map[interface{}]interface{}) bool {
		store := yaml.YAMLStore{}
		in, err := goyaml.Marshal(tree)
		if err != nil {
			t.Error(err)
		}
		store.Data = tree
		enc, err := store.Dump(key)
		if err != nil {
			t.Errorf("Error dumping: %s", err)
		}
		err = store.Load(enc, key)
		if err != nil {
			t.Error(err)
		}
		out, err := goyaml.Marshal(store.Data)
		if err != nil {
			t.Error(err)
		}
		if string(in) != string(out) {
			fmt.Printf("Expected %q, got %q\n", string(in), string(out))
			return false
		}
		return true
	}
	m := make(map[interface{}]interface{})
	m1 := make(map[interface{}]interface{})
	m1["bar"] = "baz"
	m["foo"] = m1
	m["baz"] = 2
	m["bar"] = []interface{}{1, 2, 3}
	m["foobar"] = false
	if !f(m) {
		t.Error(nil)
	}
}

func TestYamlMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.yaml")
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	store := yaml.YAMLStore{}
	sops, err := store.Metadata(string(in))
	expectedVersion := "1.13"
	if err != nil {
		t.Errorf("yaml parsing threw error: %s", err)
	}
	if sops.Version != expectedVersion {
		t.Errorf("version should be %s, was %s", expectedVersion, sops.Version)
	}
}

func TestJsonMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.json")
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	store := json.JSONStore{}
	sops, err := store.Metadata(string(in))
	expectedVersion := "1.13"
	if err != nil {
		t.Errorf("json parsing thew error: %s", err)
	}
	if sops.Version != expectedVersion {
		t.Errorf("version should be %s, was %s", expectedVersion, sops.Version)
	}
}

func TestDecryptSimpleJSON(t *testing.T) {
	in := `{"foo": "ENC[AES256_GCM,data:3mDI,iv:2OL363jDglmOa+k6qCIh5RGUm+isJNrqP4umqOEb+1s=,tag:bJiC+QHsSQPJFXCX94sRGQ==,type:str]"}`
	key := strings.Repeat("f", 32)
	expected := "foo"
	store := json.JSONStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	if store.Data["foo"] != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", store.Data["foo"], expected)
	}
}
