package mail

import (
	"container/list"
	"os"
	"path/filepath"
)

type MessageAttachment struct {
	Name    string
	Content []byte
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
}

func NewMessage() *Message {
	return &Message{
		attachments: list.New(),
	}
}

func (m *Message) AttachFile(src string, name string) error {
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

	m.attachments.PushBack(&MessageAttachment{attachName, b})
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
