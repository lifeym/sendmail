package mail

import (
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/lifeym/she/genericlist"
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
	Attachments genericlist.GenericList[*MessageAttachment]
	Header      mail.Header
}

func NewMessage() *Message {
	return &Message{
		Header: make(mail.Header),
	}
}

func (m *Message) AddHeader(field string, value string) {
	textproto.MIMEHeader(m.Header).Add(field, value)
}

func (m *Message) SetHeader(field string, value string) {
	textproto.MIMEHeader(m.Header).Set(field, value)
}

func (m *Message) GetHeader(field string) string {
	return m.Header.Get(field)
}

func (m *Message) RemoveHeader(field string) {
	textproto.MIMEHeader(m.Header).Del(field)
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
	m.Attachments.Append(&result)
	return nil
}

func (m *Message) AddressList(key string) ([]*mail.Address, error) {
	hdr := m.Header.Get(key)
	if hdr == "" {
		return nil, mail.ErrHeaderNotPresent
	}

	return mail.ParseAddressList(hdr)
}

func (m *Message) ToBytes() ([]byte, error) {
	mb := newMessageBuilder()
	return mb.Build(m)
}
