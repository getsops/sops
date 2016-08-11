package sops

import (
	"go.mozilla.org/sops/yaml"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDecryptSimpleYAML(t *testing.T) {
	in := `foo: "ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]"`
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
	in := "foo:\n  - bar: \"ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]\""
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
	sops, err := JSONStore{}.Metadata(string(in))
	expectedVersion := "1.13"
	if err != nil {
		t.Errorf("json parsing thew error: %s", err)
	}
	if sops.Version != expectedVersion {
		t.Errorf("version should be %s, was %s", expectedVersion, sops.Version)
	}
}

func TestDecryptSimpleJSON(t *testing.T) {
	json := `{"foo": "ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]"}`
	key := strings.Repeat("f", 32)
	expected := "{\"foo\":\"foo\"}"
	decryption, err := JSONStore{}.Decrypt(json, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	if decryption != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", decryption, expected)
	}
}
