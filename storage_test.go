package sops

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDecryptSimpleYAML(t *testing.T) {
	yaml := `foo: "ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]"`
	key := strings.Repeat("f", 32)
	expected := "foo: foo\n"
	decryption, err := YAMLStore{}.Decrypt(yaml, key)
	if err != nil {
		t.Error("Decryption failed: %s", err)
	}
	if decryption != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", decryption, expected)
	}
}

func TestDecryptNestedYaml(t *testing.T) {
	yaml := "foo:\n  - bar: \"ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]\""
	key := strings.Repeat("f", 32)
	expected := "foo:\n- bar: foo\n"
	decryption, err := YAMLStore{}.Decrypt(yaml, key)
	if err != nil {
		t.Error("Decryption failed: %s", err)
	}
	if decryption != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", decryption, expected)
	}
}

func TestYamlMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.yaml")
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	sops, err := YAMLStore{}.Metadata(string(in))
	expected_version := "1.13"
	if err != nil {
		t.Error("yaml parsing threw error: %s", err)
	}
	if sops.Version != expected_version {
		t.Error("version should be %s, was %s", expected_version, sops.Version)
	}
}

func TestJsonMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.json")
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	sops, err := JSONStore{}.Metadata(string(in))
	expected_version := "1.13"
	if err != nil {
		t.Error("json parsing thew error: %s", err)
	}
	if sops.Version != expected_version {
		t.Error("version should be %s, was $s", expected_version, sops.Version)
	}
}

func TestDecryptSimpleJSON(t *testing.T) {
	json := `{"foo": "ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]"}`
	key := strings.Repeat("f", 32)
	expected := "{\"foo\":\"foo\"}"
	decryption, err := JSONStore{}.Decrypt(json, key)
	if err != nil {
		t.Error("Decryption failed: %s", err)
	}
	if decryption != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", decryption, expected)
	}
}
