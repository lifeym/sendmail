package config

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/lifeym/she/mail"
)

type CompiledSmtpConfig struct {
	Name     string
	Host     string
	Port     int
	StartTLS bool
}

func compileSmtpConfig(t *SheTemplate, sc *SmtpConfig) (*CompiledSmtpConfig, error) {
	result := CompiledSmtpConfig{}
	s, err := t.Execute(sc.Name, nil)
	if err != nil {
		return nil, err
	}

	// host
	result.Name = s
	s, err = t.Execute(sc.Host, nil)
	if err != err {
		return nil, err
	}

	result.Host = s

	// port
	s, err = t.Execute(sc.Port, nil)
	if err != err {
		return nil, err
	}

	var i int
	i, err = strconv.Atoi(s)
	if err != nil {
		return nil, err
	}

	result.Port = i

	// starttls
	var b bool
	s, err = t.Execute(sc.StartTLS, nil)
	if err != err {
		return nil, err
	}

	b, err = strconv.ParseBool(s)
	if err != nil {
		return nil, err
	}

	result.StartTLS = b
	return &result, nil
}

// type CompiledMailConfig = mailConfig

// func newCompiledMailConfig() *CompiledMailConfig {
// 	result := CompiledMailConfig{}
// 	result.Spec.Header = make(mailHeaderData)
// 	return &result
// }

// func compileMailConfig(t *SheTemplate, mc *mailConfig) (*CompiledMailConfig, error) {
// 	result := newCompiledMailConfig()
// 	var err error
// 	if result.Name, err = t.Execute(mc.Name, nil); err != nil {
// 		return nil, err
// 	}

// 	if result.Template, err = t.Execute(mc.Template, nil); err != nil {
// 		return nil, err
// 	}

// 	if result.Spec.Body, err = t.Execute(mc.Spec.Body, nil); err != nil {
// 		return nil, err
// 	}

// 	for k := range mc.Spec.Header {
// 		for _, v := range mc.Spec.Header[k] {
// 			cv, err := t.Execute(v, nil)
// 			if err != nil {
// 				return nil, err
// 			}

// 			result.Spec.Header.Add(k, cv)
// 		}
// 	}

// 	for _, att := range mc.Spec.Attachments {
// 		compiledAtt := messageAttachment{}
// 		if compiledAtt.Name, err = t.Execute(att.Name, nil); err != nil {
// 			return nil, err
// 		}

// 		if compiledAtt.Path, err = t.Execute(att.Path, nil); err != nil {
// 			return nil, err
// 		}

// 		for k := range att.Header {
// 			for _, v := range att.Header[k] {
// 				cv, err := t.Execute(v, nil)
// 				if err != nil {
// 					return nil, err
// 				}

// 				compiledAtt.Header.Add(k, cv)
// 			}
// 		}

// 		result.Spec.Attachments = append(result.Spec.Attachments, compiledAtt)
// 	}

// 	return result, nil
// }

// type CompiledMessageTemplate = messageTemplate

// func newCompiledMessageTemplate() *CompiledMessageTemplate {
// 	return &CompiledMessageTemplate{
// 		Header: make(mailHeaderData),
// 	}
// }

// func compileMessageTemplate(t *SheTemplate, mt *messageTemplate) (*CompiledMessageTemplate, error) {
// 	result := newCompiledMessageTemplate()
// 	var err error
// 	if result.Name, err = t.Execute(result.Name, nil); err != nil {
// 		return nil, err
// 	}

// 	if result.Body, err = t.Execute(result.Body, nil); err != nil {
// 		return nil, err
// 	}

// 	for k := range mt.Header {
// 		for _, v := range mt.Header[k] {
// 			cv, err := t.Execute(v, nil)
// 			if err != nil {
// 				return nil, err
// 			}

// 			result.Header.Add(k, cv)
// 		}
// 	}

// 	for _, att := range mt.Attachments {
// 		compiledAtt := messageAttachment{}
// 		if compiledAtt.Name, err = t.Execute(att.Name, nil); err != nil {
// 			return nil, err
// 		}

// 		if compiledAtt.Path, err = t.Execute(att.Path, nil); err != nil {
// 			return nil, err
// 		}

// 		for k := range att.Header {
// 			for _, v := range att.Header[k] {
// 				cv, err := t.Execute(v, nil)
// 				if err != nil {
// 					return nil, err
// 				}

// 				compiledAtt.Header.Add(k, cv)
// 			}
// 		}

// 		result.Attachments = append(result.Attachments, compiledAtt)
// 	}

// 	return result, nil
// }

type CompiledMail struct {
	LoginUser string
	Password  string
	Smtp      *CompiledSmtpConfig
	Message   *mail.Message
}

func CompileMail(appCfg *AppConfig, mf *MessageFile, accountName string, mailName string) (*CompiledMail, error) {
	// data := make(map[string]any)
	t := NewTemplate()
	var err error
	result := CompiledMail{}
	account := appCfg.GetAccount(accountName)
	if account == nil {
		return nil, fmt.Errorf("account definition not found: %s", accountName)
	}

	if result.LoginUser, err = t.Execute(account.LoginUser, nil); err != nil {
		return nil, err
	}

	if result.Password, err = t.Execute(account.Password, nil); err != nil {
		return nil, err
	}

	var csmtpRef string
	csmtpRef, err = t.Execute(account.SmtpRef, nil)
	if err != nil {
		return nil, err
	}

	smtpHost := appCfg.GetSmtp(csmtpRef)
	if smtpHost == nil {
		return nil, fmt.Errorf("smtp config not found: %s", account.SmtpRef)
	}

	result.Smtp, err = compileSmtpConfig(t, smtpHost)
	if err != nil {
		return nil, err
	}

	mc := mf.GetMail(mailName)
	if mc == nil {
		return nil, fmt.Errorf("mail definition not found: %s", mailName)
	}

	mt := mf.GetTemplate(mc.Template)
	if mt == nil {
		return nil, fmt.Errorf("message definition not found: %s", mc.Template)
	}

	msg := mail.NewMessage()

	// construct message header
	for k := range mt.Header {
		if _, ok := mc.Spec.Header[k]; ok {
			continue
		}

		for _, v := range mt.Header[k] {
			cv, err := t.Execute(v, nil)
			if err != nil {
				return nil, err
			}

			msg.AddHeader(k, cv)
		}
	}

	for k := range mc.Spec.Header {
		for _, v := range mc.Spec.Header[k] {
			cv, err := t.Execute(v, nil)
			if err != nil {
				return nil, err
			}

			msg.AddHeader(k, cv)
		}
	}

	if msg.GetHeader("from") == "" {
		cv, err := t.Execute(account.DefaultFrom, nil)
		if err != nil {
			return nil, err
		}

		msg.SetHeader("from", cv)
	}

	// body
	var cbody string
	cbody, err = t.Execute(mc.Spec.Body, nil)
	if err != nil {
		return nil, err
	}

	if cbody == "" {
		cbody, err = t.Execute(mt.Body, nil)
		if err != nil {
			return nil, err
		}
	}

	msg.Body = cbody

	// attachements
	for _, att := range slices.Concat(mc.Spec.Attachments, mt.Attachments) {
		compiledAtt, err := compileAttachment(&att, t)
		if err != nil {
			return nil, err
		}

		err = msg.AttachFile(compiledAtt.Path, compiledAtt.Name, compiledAtt.Header.ToMIMEHeader())
		if err != nil {
			return nil, fmt.Errorf("cannot attach file: %s. Err: %s", att.Path, err)
		}
	}

	result.Message = msg
	return &result, nil
}

func compileAttachment(att *messageAttachment, t *SheTemplate) (*messageAttachment, error) {
	var err error
	compiledAtt := messageAttachment{}
	if compiledAtt.Name, err = t.Execute(att.Name, nil); err != nil {
		return nil, err
	}

	if compiledAtt.Path, err = t.Execute(att.Path, nil); err != nil {
		return nil, err
	}

	for k := range att.Header {
		for _, v := range att.Header[k] {
			cv, err := t.Execute(v, nil)
			if err != nil {
				return nil, err
			}

			compiledAtt.Header.Add(k, cv)
		}
	}

	return &compiledAtt, nil
}
