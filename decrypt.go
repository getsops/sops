package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"ulfr.io/sops/sops"
)

func DecryptFile(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	encYamlMap := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(fileBytes, encYamlMap); err != nil {
		return err
	}

	sopsBytes, err := yaml.Marshal(encYamlMap["sops"])
	if err != nil {
		return err
	}

	sopsData, err := sops.NewData(sopsBytes)
	if err != nil {
		return err
	}

	orderedMap := make(yaml.MapSlice, 0)
	err = yaml.Unmarshal(fileBytes, &orderedMap)

	decOrderedMap := sopsData.DecryptMapSlice(orderedMap, "")
	out, err := yaml.Marshal(decOrderedMap)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}
