package configuration

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
	"path"
)

const configFilePath = "config/.env.yaml"

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
}

func getPath() string {
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(workingDir, configFilePath)
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content
}

func Read() *Config {
	filePath := getPath()
	content := readFile(filePath)

	config := &Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		panic(err)
	}

	return config
}
