package yaml

import (
	"fmt"
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
	store := YAMLStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	if store.Data[0].Value != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", store.Data[0].Value, expected)
	}
}

func TestDecryptNestedYaml(t *testing.T) {
	in := "foo:\n  - bar: \"ENC[AES256_GCM,data:tFAp,iv:5G6/F6LaoYgIV5je0TBvqLoyyo6IZK7rCjA7T5nuy1k=,tag:4M/t1ZN5ZRWBiMrYdRpBhg==,type:str]\""
	key := strings.Repeat("f", 32)
	expected := "foo"
	store := YAMLStore{}
	err := store.Load(in, key)
	if err != nil {
		t.Errorf("Decryption failed: %s", err)
	}
	foo := store.Data[0]
	if foo.Key != "foo" {
		t.Errorf("Key does not match: %s != %s", foo.Key, "foo")
	}
	bar := foo.Value.([]interface{})[0].(goyaml.MapSlice)[0]
	if bar.Value != expected {
		t.Errorf("Decryption does not match expected result: %q != %q", bar.Value, expected)
	}
}

func TestYamlRoundtrip(t *testing.T) {
	key := strings.Repeat("f", 32)
	f := func(tree goyaml.MapSlice) bool {
		store := YAMLStore{}
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
	m := make(goyaml.MapSlice, 0)
	m1 := make(goyaml.MapSlice, 0)
	m1 = append(m1, goyaml.MapItem{Key: "bar", Value: "baz"})
	m = append(m, goyaml.MapItem{Key: "foo", Value: m1})
	m = append(m, goyaml.MapItem{Key: "baz", Value: 2})
	m = append(m, goyaml.MapItem{Key: "bar", Value: []interface{}{1, 2, 3}})
	m = append(m, goyaml.MapItem{Key: "foobar", Value: false})
	if !f(m) {
		t.Error(nil)
	}
}

func TestYamlMetadata(t *testing.T) {
	file, err := os.Open("test_resources/example.yaml")
	if err != nil {
		t.Errorf("Could not open test resource file: %s", err)
	}
	defer file.Close()
	in, err := ioutil.ReadAll(file)
	store := YAMLStore{}
	err = store.GetMetadata(string(in))
	expectedVersion := "1.13"
	if err != nil {
		t.Errorf("yaml parsing threw error: %s", err)
	}
	if store.Metadata.Version != expectedVersion {
		t.Errorf("version should be %s, was %s", expectedVersion, store.Metadata.Version)
	}
}
