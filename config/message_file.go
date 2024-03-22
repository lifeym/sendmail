package config

import (
	"net/textproto"
	"os"

	"gopkg.in/yaml.v3"
)

type mailHeaderData map[string]StringArray

func (h mailHeaderData) ToMIMEHeader() textproto.MIMEHeader {
	result := make(textproto.MIMEHeader)
	for k := range h {
		for _, v := range h[k] {
			result.Add(k, v)
		}
	}

	return result
}

// Message file
type messageTemplate struct {
	Name        string
	Header      mailHeaderData
	Body        string
	Attachments []struct {
		Name   string
		Path   string
		Header mailHeaderData
	}
}

type messageSpec struct {
	Header      mailHeaderData
	Body        string
	Attachments []struct {
		Name   string
		Path   string
		Header mailHeaderData
	}
}

type mailConfig struct {
	Name     string
	Template string
	Spec     messageSpec
}

type MessageFile struct {
	Templates          []messageTemplate
	Mails              []mailConfig
	messageTemplateMap map[string]*messageTemplate
	mailMap            map[string]*mailConfig
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

	mf.messageTemplateMap = make(map[string]*messageTemplate)
	for _, ts := range mf.Templates {
		mf.messageTemplateMap[ts.Name] = &ts
	}

	mf.mailMap = make(map[string]*mailConfig)
	for _, ms := range mf.Mails {
		mf.mailMap[ms.Name] = &ms
	}

	return &mf, nil
}

func (mf *MessageFile) GetTemplate(name string) *messageTemplate {
	return mf.messageTemplateMap[name]
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
