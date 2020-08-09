//
// Reader will return the configuration structure
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func LoadConfiguration(path string) (*Configuration, error) {
	var config Configuration

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
