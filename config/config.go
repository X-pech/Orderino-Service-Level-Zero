package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type PostgresConfig struct {
	URL        string `yaml:"url"`
	Port       int    `yaml:"port"`
	Db         string `yaml:"db"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	TableName  string `yaml:"tableName"`
	IDColumn   string `yaml:"idColumn"`
	DataColumn string `yaml:"dataColumn"`
}

type StanConfig struct {
	URL           string `yaml:"url"`
	ListenerName  string `yaml:"listenerName"`
	ClusterName   string `yaml:"clusterName"`
	TopicName     string `yaml:"topicName"`
	PublisherName string `yaml:"testPublisherName"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Stan     StanConfig     `yaml:"stan"`
	App      AppConfig      `yaml:"app"`
}

func New(filename string) (Config, error) {
	file, err := ioutil.ReadFile(filename)

	if err != nil {
		return Config{}, err
	}

	result := new(Config)
	err = yaml.Unmarshal(file, result)

	return *result, err
}
