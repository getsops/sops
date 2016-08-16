package json

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestJsonMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.json")
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	store := JSONStore{}
	sops, err := store.Metadata(string(in))
	expectedVersion := "1.13"
	if err != nil {
		t.Errorf("json parsing threw error: %s", err)
	}
	if sops.Version != expectedVersion {
		t.Errorf("version should be %s, was %s", expectedVersion, sops.Version)
	}
}

func TestDecryptSimpleJSON(t *testing.T) {
	in := `{"foo": "ENC[AES256_GCM,data:3mDI,iv:2OL363jDglmOa+k6qCIh5RGUm+isJNrqP4umqOEb+1s=,tag:bJiC+QHsSQPJFXCX94sRGQ==,type:str]"}`
	key := strings.Repeat("f", 32)
	expected := "foo"
	store := JSONStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	if store.Data["foo"] != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", store.Data["foo"], expected)
	}
}

func TestJsonRoundtrip(t *testing.T) {
	key := strings.Repeat("f", 32)
	f := func(tree map[string]interface{}) bool {
		store := JSONStore{}
		in, err := json.Marshal(tree)
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
		out, err := json.Marshal(store.Data)
		if err != nil {
			t.Error(err)
		}
		if string(in) != string(out) {
			fmt.Printf("Expected %q, got %q\n", string(in), string(out))
			return false
		}
		return true
	}
	m := make(map[string]interface{})
	m1 := make(map[string]interface{})
	m1["bar"] = "baz"
	m["foo"] = m1
	m["baz"] = 2
	m["bar"] = []interface{}{1, 2, 3}
	m["foobar"] = false
	if !f(m) {
		t.Error(nil)
	}
}
