package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

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

func (mb *messageBuilder) appendFiled(name string, body string) {
	mb.appendLine(fmt.Sprintf("%s: %s", name, body))
}

func (mb *messageBuilder) appendFrom(from string) {
	mb.appendFiled("From", from)
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

func (mb *messageBuilder) appendMessageAttachments(m *Message, mw *multipart.Writer) error {

	for el := m.attachments.Front(); el != nil; el = el.Next() {
		att := el.Value.(*MessageAttachment)
		// mb.appendEmptyLines(2)
		header := make(textproto.MIMEHeader)
		header.Add("Content-Type", http.DetectContentType(att.Content))
		header.Add("Content-Transfer-Encoding", "base64")
		header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", att.Name))
		w, err := mw.CreatePart(header)
		if err != nil {
			return err
		}

		// mb.appendLine(fmt.Sprintf("--%s", boundary))
		// mb.appendFiled("Content-Type", http.DetectContentType(att.Content))
		// mb.appendFiled("Content-Transfer-Encoding", "base64")
		// mb.appendFiled("Content-Disposition", fmt.Sprintf("attachment; filename=%s", att.Name))

		bs := make([]byte, base64.StdEncoding.EncodedLen(len(att.Content)))
		base64.StdEncoding.Encode(bs, att.Content)
		w.Write(bs)
		// mb.buf.Write(bs)
		// mb.appendEmptyLine()
		// mb.appendString(fmt.Sprintf("--%s", boundary))
	}

	// mb.appendString("--")
	return nil
}

func (mb *messageBuilder) Build(m *Message) []byte {
	withAttachments := m.attachments.Len() > 0
	mb.appendFrom(m.From)
	mb.appendSubject(m.Subject)
	mb.appendTo(m.To)
	if len(m.Cc) > 0 {
		mb.appendCc(m.Cc)
	}

	if len(m.Bcc) > 0 {
		mb.appendBcc(m.Bcc)
	}

	// mb.appendFiled("MIME-Version", "1.0")

	if withAttachments {
		mw := multipart.NewWriter(mb.buf)
		boundary := mw.Boundary()
		mb.appendFiled("Content-Type", fmt.Sprintf("multipart/mixed; boundary=%s", boundary))
		mb.appendEmptyLine()
		mpHeader := make(textproto.MIMEHeader)
		mpHeader.Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", http.DetectContentType([]byte(m.Body))))
		w, _ := mw.CreatePart(mpHeader)
		w.Write([]byte(m.Body))
		mb.appendMessageAttachments(m, mw)
		mb.appendEmptyLines(2)
		mb.appendString(fmt.Sprintf("--%s--", boundary))
		// mb.appendLine(fmt.Sprintf("--%s", boundary))
	} else {
		mb.appendFiled("Content-Type", "text/plain; charset=utf-8")
		mb.appendString(m.Body)
	}

	fmt.Println(string(mb.buf.Bytes()))
	return mb.buf.Bytes()
}
