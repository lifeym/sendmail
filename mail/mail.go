package mail

import (
	"bytes"
	"container/list"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
)

// SmtpAuth contains informations for connecting to a smtp server
// and functions to interactive with mails
type SmtpAuth struct {
	auth     smtp.Auth
	host     string
	hostPort int
}

func New(username string, password string, host string, hostPort int) *SmtpAuth {
	return &SmtpAuth{
		smtp.PlainAuth("", username, password, host),
		host,
		hostPort,
	}
}

func (s *SmtpAuth) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.hostPort), s.auth, m.From, m.To, m.ToBytes())
}

//-----------------------------------------------------------------------------

// Utility for building mail message
type messageBuilder struct {
	buf *bytes.Buffer
}

func newMessageBuilder() *messageBuilder {
	return &messageBuilder{
		buf: bytes.NewBuffer(nil),
	}
}

func (mb *messageBuilder) appendLine(s string) {
	mb.buf.WriteString(fmt.Sprintf("%s\r\n", s))
}

func (mb *messageBuilder) appendEmptyLine() {
	mb.buf.WriteString("\r\n")
}

func (mb *messageBuilder) appendEmptyLines(count int) {
	mb.buf.WriteString(strings.Repeat("\r\n", count))
}

func (mb *messageBuilder) appendString(s string) {
	mb.buf.WriteString(s)
}

func (mb *messageBuilder) appendFiled(f string, v string) {
	mb.appendLine(fmt.Sprintf("%s: %s", f, v))
}

func (mb *messageBuilder) appendSubject(subject string) {
	mb.appendFiled("Subject", subject)
}

func (mb *messageBuilder) appendTo(mailto []string) {
	mb.appendFiled("To", strings.Join(mailto, ","))
}

func (mb *messageBuilder) appendCc(cc []string) {
	mb.appendFiled("Cc", strings.Join(cc, ","))
}

func (mb *messageBuilder) appendBcc(bcc []string) {
	mb.appendFiled("Bcc", strings.Join(bcc, ","))
}

func (mb *messageBuilder) appendMessageAttachments(m *Message) {
	writer := multipart.NewWriter(mb.buf)
	boundary := writer.Boundary()
	for el := m.attachments.Front(); el != nil; el = el.Next() {
		att := el.Value.(*MessageAttachment)
		mb.appendEmptyLines(2)
		mb.appendLine(fmt.Sprintf("--%s", boundary))
		mb.appendFiled("Content-Type", http.DetectContentType(att.Content))
		mb.appendFiled("Content-Transfer-Encoding", "base64")
		mb.appendFiled("Content-Disposition", fmt.Sprintf("attachment; filename=%s", att.Name))

		bs := make([]byte, base64.StdEncoding.EncodedLen(len(att.Content)))
		base64.StdEncoding.Encode(bs, att.Content)
		mb.buf.Write(bs)
		mb.appendEmptyLine()
		mb.appendString(fmt.Sprintf("--%s", boundary))
	}

	mb.appendString("--")
}

func (mb *messageBuilder) Build(m *Message) []byte {
	withAttachments := m.attachments.Len() > 0
	mb.appendSubject(m.Subject)
	mb.appendTo(m.To)
	if len(m.Cc) > 0 {
		mb.appendCc(m.Cc)
	}

	if len(m.Bcc) > 0 {
		mb.appendBcc(m.Bcc)
	}

	mb.appendFiled("MIME-Version", "1.0")
	writer := multipart.NewWriter(mb.buf)
	boundary := writer.Boundary()
	if withAttachments {
		mb.appendFiled("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", boundary))
		mb.appendLine(fmt.Sprintf("--%s", boundary))
	} else {
		mb.appendFiled("Content-Type", "text/plain; charset=utf-8")
	}

	mb.appendString(m.Body)
	if withAttachments {
		mb.appendMessageAttachments(m)
	}

	return mb.buf.Bytes()
}

//-----------------------------------------------------------------------------

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
