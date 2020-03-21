package configuration

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"github.com/hasangenc0/corona/pkg/helpers"
)

const configPathTemplate = "config/.env.%s.yaml"

type Config struct {
	Db struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Timeout  int32  `yaml:"timeout"`
	}

	Server struct {
		Port    string `yaml:"port"`
		Timeout string `yaml:"timeout"`
	}

	Api struct {
		Corona string `yaml:"corona"`
	}
}

func Read(env string) *Config {
	filePath := helpers.GetPath(fmt.Sprintf(configPathTemplate, env))
	content := helpers.ReadFile(filePath)

	config := &Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		panic(err)
	}

	return config
}
