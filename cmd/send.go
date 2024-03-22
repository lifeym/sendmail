package cmd

import (
	"fmt"
	"slices"
	"time"

	"github.com/lifeym/she/config"
	"github.com/lifeym/she/mail"
	"github.com/spf13/cobra"
)

var (
	_account string
	_mail    string
	_config  string
	_print   bool
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send -c -a -m",
	Long: `she命令行邮件客户端
                   Copyright lifeym 2024`,
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := send(_account, _mail, _config); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	sendCmd.Flags().StringVarP(&_account, "account", "a", "", `Account config name in config file.`)
	sendCmd.Flags().StringVarP(&_mail, "mail", "m", "", `Mail config names in config file to be sent, default to all.`)
	sendCmd.Flags().StringVarP(&_config, "message-file", "f", "", `Mail message config file.`)
	sendCmd.Flags().BoolVarP(&_print, "print", "p", false, `Print mail message content to stdout.`)
	sendCmd.MarkFlagRequired("account")
	// sendCmd.MarkFlagRequired("mail")
	sendCmd.MarkFlagRequired("config")
	// rootCmd.PersistentFlags().StringVarP(&_dsn, "datasource", "s", "", `datasource url for connecting to a database, run hare -h database for more  details.`)
	// rootCmd.Flags().StringVarP(&_output, "output", "o", "", `output file name of generated results(use standard output by default)`)
	// rootCmd.Flags().StringArrayVarP(&_template, "template", "t", nil, `specify templates to use`)
	// rootCmd.Flags().StringVarP(&_ldelim, "ldelim", "l", "", `specify left delimiters for template definition, must be used with rdelim together`)
	// rootCmd.Flags().StringVarP(&_rdelim, "rdelim", "r", "", `specify right delimiters for template definition, must be used with ldelim together`)
	// rootCmd.PersistentFlags().StringArrayVarP(&_var, "var", "v", nil, `var=value`)
	rootCmd.AddCommand(sendCmd)
}

func send(accountRef string, mailRef string, cfgPath string) error {
	cfg, err := config.LoadConfigFile(".sendmail.yaml")
	if err != nil {
		return err
	}

	msgFile, err := config.LoadMessageFile(cfgPath)
	if err != nil {
		return err
	}

	account := cfg.GetAccount(accountRef)
	if account == nil {
		return fmt.Errorf("account definition not found: %s", accountRef)
	}

	mailDef := msgFile.GetMail(mailRef)
	if mailDef == nil {
		return fmt.Errorf("mail definition not found: %s", mailRef)
	}

	msgTpl := msgFile.GetTemplate(mailDef.Template)
	if msgTpl == nil {
		return fmt.Errorf("message definition not found: %s", mailDef.Template)
	}

	msg := mail.NewMessage()

	// construct message header
	for k := range msgTpl.Header {
		for _, v := range msgTpl.Header[k] {
			msg.AddHeader(k, v)
		}
	}

	for k := range mailDef.Spec.Header {
		msg.RemoveHeader(k)
		for _, v := range mailDef.Spec.Header[k] {
			msg.AddHeader(k, v)
		}
	}

	if msg.GetHeader("from") == "" {
		msg.SetHeader("from", account.DefaultFrom)
	}

	if msg.GetHeader("from") == "" {
		return fmt.Errorf("mail: header missing or empty -- %s", "from")
	}

	if msg.GetHeader("to") == "" {
		return fmt.Errorf("mail: header missing or empty -- %s", "to")
	}

	// body
	if mailDef.Spec.Body != "" {
		msg.Body = mailDef.Spec.Body
	} else {
		msg.Body = msgTpl.Body
	}

	// attachements
	for _, att := range slices.Concat(mailDef.Spec.Attachments, msgTpl.Attachments) {
		err = msg.AttachFile(att.Path, att.Name, att.Header.ToMIMEHeader())
		if err != nil {
			return fmt.Errorf("cannot attach file: %s. Err: %s", att.Path, err)
		}
	}

	// date
	if msg.GetHeader("date") == "" {
		msg.SetHeader("date", time.Now().Format(time.RFC1123Z))
	}

	if _print {
		msgData, err := msg.ToBytes()
		if err != nil {
			return nil
		}

		fmt.Println(string(msgData))
	}

	smtpHost := cfg.GetSmtp(account.SmtpRef)
	smtp := mail.New(account.LoginUser, account.Password, smtpHost.Host, smtpHost.Port, smtpHost.StartTLS)
	err = smtp.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
