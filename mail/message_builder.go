package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"slices"
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

func (mb *messageBuilder) writeLine(s string) {
	mb.buf.WriteString(fmt.Sprintf("%s\r\n", s))
}

func (mb *messageBuilder) writeEmptyLine() {
	mb.buf.WriteString("\r\n")
}

func (mb *messageBuilder) appendString(s string) {
	mb.buf.WriteString(s)
}

func (mb *messageBuilder) writeFiled(name string, value string) {
	mb.writeLine(fmt.Sprintf("%s: %s", name, value))
}

func (mb *messageBuilder) appendFrom(from string) {
	mb.writeFiled("From", from)
}

func (mb *messageBuilder) appendSubject(subject string) {
	mb.writeFiled("Subject", subject)
}

func (mb *messageBuilder) appendTo(mailto []string) {
	mb.writeFiled("To", strings.Join(mailto, ","))
}

func (mb *messageBuilder) appendCc(cc []string) {
	mb.writeFiled("Cc", strings.Join(cc, ","))
}

func (mb *messageBuilder) appendBcc(bcc []string) {
	mb.writeFiled("Bcc", strings.Join(bcc, ","))
}

func (mb *messageBuilder) writeHeaders(headers textproto.MIMEHeader) {
	if headers == nil {
		return
	}

	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	for _, k := range keys {
		for _, v := range headers[k] {
			mb.writeFiled(textproto.CanonicalMIMEHeaderKey(k), v)
		}
	}
}

func (mb *messageBuilder) appendMessageAttachments(m *Message, mw *multipart.Writer) error {
	for el := m.attachments.Front(); el != nil; el = el.Next() {
		att := el.Value.(*MessageAttachment)
		if att.Headers.Get("Content-Type") == "" {
			att.Headers.Add("Content-Type", http.DetectContentType(att.Content))
		}

		if att.Headers.Get("Content-Transfer-Encoding") == "" {
			att.Headers.Add("Content-Transfer-Encoding", "base64")
		}

		if att.Headers.Get("Content-Disposition") == "" {
			att.Headers.Add("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", att.Name))
		}

		w, err := mw.CreatePart(att.Headers)
		if err != nil {
			return err
		}

		bs := make([]byte, base64.StdEncoding.EncodedLen(len(att.Content)))
		base64.StdEncoding.Encode(bs, att.Content)
		w.Write(bs)
	}

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

	mb.writeHeaders(m.headers)
	// mb.appendFiled("MIME-Version", "1.0")

	if withAttachments {
		mw := multipart.NewWriter(mb.buf)
		boundary := mw.Boundary()
		mb.writeFiled("Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary))
		mb.writeEmptyLine()
		mpHeader := make(textproto.MIMEHeader)
		mpHeader.Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", http.DetectContentType([]byte(m.Body))))
		w, _ := mw.CreatePart(mpHeader)
		w.Write([]byte(m.Body))
		mb.appendMessageAttachments(m, mw)
		mw.Close()
		// mb.appendLine(fmt.Sprintf("--%s", boundary))
	} else {
		mb.writeFiled("Content-Type", "text/plain; charset=utf-8")
		mb.appendString(m.Body)
	}

	return mb.buf.Bytes()
}
