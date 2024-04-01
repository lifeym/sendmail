package cmd

import (
	"fmt"
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
	Short: "send -f -a -m",
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
	sendCmd.Flags().StringVarP(&_mail, "message", "m", "", `Message names in message file to be sent, default to all.`)
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

	var compiledMail *config.CompiledMail
	if compiledMail, err = config.CompileMail(cfg, msgFile, accountRef, mailRef); err != nil {
		return err
	}

	if compiledMail.Message.GetHeader("from") == "" {
		return fmt.Errorf("mail: header missing or empty -- %s", "from")
	}

	if compiledMail.Message.GetHeader("to") == "" {
		return fmt.Errorf("mail: header missing or empty -- %s", "to")
	}

	// date
	if compiledMail.Message.GetHeader("date") == "" {
		compiledMail.Message.SetHeader("date", time.Now().Format(time.RFC1123Z))
	}

	if _print {
		msgData, err := compiledMail.Message.ToBytes()
		if err != nil {
			return nil
		}

		fmt.Println(string(msgData))
	}

	smtp := mail.New(compiledMail.LoginUser, compiledMail.Password, compiledMail.Smtp.Host, compiledMail.Smtp.Port, compiledMail.Smtp.StartTLS)
	err = smtp.Send(compiledMail.Message)
	if err != nil {
		return err
	}

	return nil
}
