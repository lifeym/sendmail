package mail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/textproto"
	"slices"
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

func (mb *messageBuilder) writeLine(s string) (int, error) {
	return mb.buf.WriteString(fmt.Sprintf("%s\r\n", s))
}

func (mb *messageBuilder) writeEmptyLine() (int, error) {
	return mb.buf.WriteString("\r\n")
}

func (mb *messageBuilder) writeString(s string) (int, error) {
	return mb.buf.WriteString(s)
}

func (mb *messageBuilder) writeFiled(name string, value string) (int, error) {
	return mb.writeLine(fmt.Sprintf("%s: %s", name, value))
}

// func (mb *messageBuilder) appendFrom(from string) {
// 	mb.writeFiled("From", from)
// }

// func (mb *messageBuilder) appendSubject(subject string) {
// 	mb.writeFiled("Subject", subject)
// }

// func (mb *messageBuilder) appendTo(mailto []string) {
// 	mb.writeFiled("To", strings.Join(mailto, ","))
// }

// func (mb *messageBuilder) appendCc(cc []string) {
// 	mb.writeFiled("Cc", strings.Join(cc, ","))
// }

func (mb *messageBuilder) writeHeader(h mail.Header) error {
	if h == nil {
		return nil
	}

	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	for _, k := range keys {
		for _, v := range h[k] {
			if _, err := mb.writeFiled(textproto.CanonicalMIMEHeaderKey(k), v); err != nil {
				return err
			}
		}
	}

	return nil
}

func writeMessageAttachments(m *Message, mw *multipart.Writer) error {
	for _, att := range m.Attachments.Data {
		headerPatchDefault(att.Header, "Content-Type", http.DetectContentType(att.Content))
		headerPatchDefault(att.Header, "Content-Transfer-Encoding", "base64")
		headerPatchDefault(att.Header, "Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", att.Name))

		w, err := mw.CreatePart(att.Header)
		if err != nil {
			return err
		}

		bs := make([]byte, base64.StdEncoding.EncodedLen(len(att.Content)))
		base64.StdEncoding.Encode(bs, att.Content)
		w.Write(bs)
	}

	return nil
}

func (mb *messageBuilder) Build(m *Message) ([]byte, error) {
	// mb.appendFrom(m.From)
	// mb.appendSubject(m.Subject)
	// mb.appendTo(m.To)
	// if len(m.Cc) > 0 {
	// 	mb.appendCc(m.Cc)
	// }
	if err := mb.writeHeader(m.Header); err != nil {
		return nil, err
	}

	// mb.appendFiled("MIME-Version", "1.0")

	if len(m.Attachments.Data) > 0 {
		mw := multipart.NewWriter(mb.buf)
		boundary := mw.Boundary()
		if _, err := mb.writeFiled("Content-Type", fmt.Sprintf("multipart/mixed; boundary=\"%s\"", boundary)); err != nil {
			return nil, err
		}

		if _, err := mb.writeEmptyLine(); err != nil {
			return nil, err
		}

		mpHeader := make(textproto.MIMEHeader)
		mpHeader.Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", http.DetectContentType([]byte(m.Body))))
		w, err := mw.CreatePart(mpHeader)
		if err != nil {
			return nil, err
		}

		if _, err = w.Write([]byte(m.Body)); err != nil {
			return nil, err
		}

		if err = writeMessageAttachments(m, mw); err != nil {
			return nil, err
		}

		if err = mw.Close(); err != nil {
			return nil, err
		}
	} else {
		if _, err := mb.writeFiled("Content-Type", "text/plain; charset=utf-8"); err != nil {
			return nil, err
		}

		if _, err := mb.writeEmptyLine(); err != nil {
			return nil, err
		}

		if _, err := mb.writeString(m.Body); err != nil {
			return nil, err
		}
	}

	return mb.buf.Bytes(), nil
}

func headerPatchDefault(header textproto.MIMEHeader, k string, v string) {
	if header.Get(k) == "" {
		header.Add(k, v)
	}
}
