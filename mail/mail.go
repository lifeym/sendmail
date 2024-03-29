package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

// SmtpAuth contains informations for connecting to a smtp server
// and functions to interactive with mails
type SmtpAuth struct {
	auth     smtp.Auth
	host     string
	hostPort int
	starttls bool
}

func New(username string, password string, host string, hostPort int, tls bool) *SmtpAuth {
	return &SmtpAuth{
		smtp.PlainAuth("", username, password, host),
		host,
		hostPort,
		tls,
	}
}

// Dial returns a new Client connected to an SMTP server at addr.
// The addr must include a port, as in "mail.example.com:smtp".
func DialInsecure(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         addr,
	})

	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func (s *SmtpAuth) SendMail(c *smtp.Client, from string, to []string, msg []byte) error {
	var err error
	if err = validateLine(from); err != nil {
		return err
	}

	for _, recp := range to {
		if err = validateLine(recp); err != nil {
			return err
		}
	}

	if ok, _ := c.Extension("STARTTLS"); ok {
		config := &tls.Config{ServerName: s.host}
		if err = c.StartTLS(config); err != nil {
			return err
		}
	}

	if s.auth != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}

		if err = c.Auth(s.auth); err != nil {
			return err
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	var w io.WriteCloser
	w, err = c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func (s *SmtpAuth) Send(m *Message) error {
	mailfrom, err := mail.ParseAddress(m.GetHeader("from"))
	if err != nil {
		return err
	}

	if err := validateLine(mailfrom.Address); err != nil {
		return err
	}

	var mailto []string
	addrs, _ := m.AddressList("to")
	if addrs == nil {
		return fmt.Errorf("mail: header not in message -- %s", "To")
	}

	for _, a := range addrs {
		mailto = append(mailto, a.Address)
	}

	addrs, _ = m.AddressList("cc")
	for _, a := range addrs {
		mailto = append(mailto, a.Address)
	}

	addrs, _ = m.AddressList("bcc")
	for _, a := range addrs {
		mailto = append(mailto, a.Address)
	}

	for _, recp := range mailto {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.hostPort)

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	var c *smtp.Client
	if s.starttls {
		c, err = smtp.Dial(addr)
	} else {
		c, err = DialInsecure(addr)
	}

	if err != nil {
		return err
	}

	defer c.Close()
	var msgData []byte
	msgData, err = m.ToBytes()
	if err != nil {
		return err
	}

	return s.SendMail(c, mailfrom.Address, mailto, msgData)
}

// validateLine checks to see if a line has CR or LF as per RFC 5321.
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}
