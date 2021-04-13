package configs

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// LoadConfigYaml 从yaml配置文件中读取配置信息
func LoadConfigYaml(yamlFiles ...string) (*Config, error) {
	conf := &Config{}
	for _, yamlFile := range yamlFiles {
		yamlByte, err := ioutil.ReadFile(yamlFile)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("read yaml(%s)", yamlFile))
		}
		err = yaml.Unmarshal(yamlByte, conf)
		if err != nil {
			return nil, errors.Wrap(err, "config yaml unmarshal")
		}
	}

	return conf, nil
}
