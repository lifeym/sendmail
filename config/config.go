package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AccountConfig struct {
	Name        string
	SmtpRef     string `yaml:"smtpRef"`
	LoginUser   string `yaml:"loginUser"`
	Password    string
	DefaultFrom string `yaml:"defaultFrom"`
}

type SmtpConfig struct {
	Name     string
	Host     string
	Port     int
	StartTLS bool `yaml:"starttls"`
}

type AppConfig struct {
	Smtp     []SmtpConfig
	Accounts []AccountConfig

	smtpMap    map[string]*SmtpConfig
	accountMap map[string]*AccountConfig
}

func LoadConfigFile(filename string) (*AppConfig, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{}
	err = yaml.Unmarshal(bs, &appConfig)
	if err != nil {
		return nil, err
	}

	appConfig.smtpMap = make(map[string]*SmtpConfig)
	for _, ss := range appConfig.Smtp {
		appConfig.smtpMap[ss.Name] = &ss
	}

	appConfig.accountMap = make(map[string]*AccountConfig)
	for _, ps := range appConfig.Accounts {
		appConfig.accountMap[ps.Name] = &ps
	}

	return &appConfig, nil
}

func (c *AppConfig) GetSmtp(name string) *SmtpConfig {
	return c.smtpMap[name]
}

func (c *AppConfig) GetAccount(name string) *AccountConfig {
	return c.accountMap[name]
}

func (c *AppConfig) SaveToFile(filename string) error {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, bs, 0644)
}

func (c *AppConfig) ToString() (string, error) {
	bs, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
