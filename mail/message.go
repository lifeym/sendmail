package mail

import (
	"container/list"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
)

type MessageAttachment struct {
	Name    string
	Content []byte
	Header  textproto.MIMEHeader
}

// Message represents a mail message to be sent by smtp server
type Message struct {
	// From        string
	// To          []string
	// Cc          []string
	// Bcc         []string
	// Subject     string
	Body        string
	attachments *list.List
	header      mail.Header
}

func NewMessage() *Message {
	return &Message{
		attachments: list.New(),
		header:      make(mail.Header),
	}
}

func (m *Message) AddHeader(field string, value string) {
	textproto.MIMEHeader(m.header).Add(field, value)
}

func (m *Message) SetHeader(field string, value string) {
	textproto.MIMEHeader(m.header).Set(field, value)
}

func (m *Message) GetHeader(field string) string {
	return m.header.Get(field)
}

func (m *Message) RemoveHeader(field string) {
	textproto.MIMEHeader(m.header).Del(field)
}

func (m *Message) AttachFile(src string, name string, header textproto.MIMEHeader) error {
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

	result := MessageAttachment{attachName, b, header}
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

func (m *Message) AddressList(key string) ([]*mail.Address, error) {
	hdr := m.header.Get(key)
	if hdr == "" {
		return nil, mail.ErrHeaderNotPresent
	}

	return mail.ParseAddressList(hdr)
}

func (m *Message) ToBytes() ([]byte, error) {
	mb := newMessageBuilder()
	return mb.Build(m)
}
