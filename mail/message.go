package mail

import (
	"container/list"
	"net/textproto"
	"os"
	"path/filepath"
)

type MessageAttachment struct {
	Name    string
	Content []byte
	Headers textproto.MIMEHeader
}

// Message represents a mail message to be sent by smtp server
type Message struct {
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	attachments *list.List
	headers     textproto.MIMEHeader
}

func NewMessage() *Message {
	return &Message{
		attachments: list.New(),
		headers:     make(textproto.MIMEHeader),
	}
}

func (m *Message) SetHeader(field string, value string) {
	m.headers.Set(field, value)
}

func (m *Message) GetHeader(field string) string {
	return m.headers.Get(field)
}

func (m *Message) RemoveHeader(field string) {
	m.headers.Del(field)
}

func (m *Message) AttachFile(src string, name string, headers map[string]string) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	var attachName string
	if name == "" {
		_, fileName := filepath.Split(src)
		attachName = fileName
	} else {
		attachName = name
	}

	result := MessageAttachment{attachName, b, make(textproto.MIMEHeader)}
	for k, v := range headers {
		result.Headers.Set(k, v)
	}

	m.attachments.PushBack(&result)
	return nil
}

func (m *Message) RemoveAttachmentByIndex(index int) {
	i := 0
	for el := m.attachments.Front(); el != nil; el = el.Next() {
		i++
		if i == index {
			m.attachments.Remove(el)
			break
		}
	}
}

func (m *Message) RemoveAttachmentByName(name string) {
	for el := m.attachments.Front(); el != nil; el = el.Next() {
		att := el.Value.(*MessageAttachment)
		if att.Name == name {
			m.attachments.Remove(el)
			break
		}
	}
}

func (m *Message) ToBytes() []byte {
	mb := newMessageBuilder()
	return mb.Build(m)
}
