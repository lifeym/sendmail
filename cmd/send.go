package cmd

import (
	"fmt"

	"github.com/lifeym/she/config"
	"github.com/lifeym/she/mail"
	"github.com/spf13/cobra"
)

var (
	_account string
	_mail    string
	_config  string
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
	sendCmd.Flags().StringVarP(&_mail, "mail", "m", "", `Mail config name in config file.`)
	sendCmd.Flags().StringVarP(&_config, "config", "c", "", `Mail config file.`)
	sendCmd.MarkFlagRequired("account")
	sendCmd.MarkFlagRequired("mail")
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

	envelope := msgFile.GetEnvelope(mailDef.EnvelopeRef)
	if envelope == nil {
		return fmt.Errorf("envelope definition not found: %s", mailDef.EnvelopeRef)
	}

	msgDef := msgFile.GetMessage(mailDef.MessageRef)
	if msgDef == nil {
		return fmt.Errorf("message definition not found: %s", mailDef.MessageRef)
	}

	msg := mail.NewMessage()
	if envelope.From == "" {
		msg.From = account.DefaultFrom
	} else {
		msg.From = envelope.From
	}

	if msg.From == "" {
		return fmt.Errorf("both Envelop.From and Account.DefaultFrom are empty")
	}

	if msgDef.Subject == "" {
		msg.Subject = envelope.Subject
	} else {
		msg.Subject = msgDef.Subject
	}

	msg.To = envelope.To
	msg.Bcc = envelope.Bcc
	msg.Cc = envelope.Cc
	msg.Body = msgDef.Body
	for _, att := range msgDef.Attachments {
		err = msg.AttachFile(att.Path, att.Name)
		if err != nil {
			return fmt.Errorf("cannot attach file: %s. Err: %s", att.Path, err)
		}
	}

	smtpHost := cfg.GetSmtp(account.SmtpRef)
	smtp := mail.New(account.LoginUser, account.Password, smtpHost.Host, smtpHost.Port, smtpHost.StartTLS)
	err = smtp.Send(msg)
	if err != nil {
		return err
	}

	return nil
}
