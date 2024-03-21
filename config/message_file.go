package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Message file
type messageConfig struct {
	Name        string
	Headers     map[string]string
	Body        string
	Attachments []struct {
		Name    string
		Path    string
		Headers map[string]string
	}
}

type envelopConfig struct {
	Name    string
	Subject string
	From    string
	To      []string
	Cc      []string
	Bcc     []string
}

type mailConfig struct {
	Name             string
	EnvelopeRef      string        `yaml:"envelopeRef"`
	EnvelopeOverride envelopConfig `yaml:"envelopeOverride"`
	MessageRef       string        `yaml:"messageRef"`
	MessageOverride  messageConfig `yaml:"messageOverride"`
}

type MessageFile struct {
	Messages    []messageConfig
	Envelopes   []envelopConfig
	Mails       []mailConfig
	messageMap  map[string]*messageConfig
	envelopeMap map[string]*envelopConfig
	mailMap     map[string]*mailConfig
}

func LoadMessageFile(filename string) (*MessageFile, error) {
	bs, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	mf := MessageFile{}
	err = yaml.Unmarshal(bs, &mf)
	if err != nil {
		return nil, err
	}

	mf.messageMap = make(map[string]*messageConfig)
	for _, ms := range mf.Messages {
		mf.messageMap[ms.Name] = &ms
	}

	mf.envelopeMap = make(map[string]*envelopConfig)
	for _, es := range mf.Envelopes {
		mf.envelopeMap[es.Name] = &es
	}

	mf.mailMap = make(map[string]*mailConfig)
	for _, ms := range mf.Mails {
		mf.mailMap[ms.Name] = &ms
	}

	return &mf, nil
}

func (mf *MessageFile) GetMessage(name string) *messageConfig {
	return mf.messageMap[name]
}

func (mf *MessageFile) GetEnvelope(name string) *envelopConfig {
	return mf.envelopeMap[name]
}

func (mf *MessageFile) GetMail(name string) *mailConfig {
	return mf.mailMap[name]
}

func (mf *MessageFile) ToString() (string, error) {
	bs, err := yaml.Marshal(mf)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
