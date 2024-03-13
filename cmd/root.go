import {
	"github.com/spf13/cobra"
}

var rootCmd = &cobra.Command{
     Use:   "she",
     Short: "she cli mail client",
     Long: `she命令行邮件客户端
                   Copyright lifeym 2024`,
     Args: cobra.ArbitraryArgs,
     RunE: func(cmd *cobra.Command, args []string) error {
         // if err := checkGlobalFlags(); err != nil {
         //  return err
         // }
 
         return nil
     },
 }

 func init() {
     // rootCmd.PersistentFlags().StringVarP(&_driverName, "driver", "d", "", `database type to be connected, run hare -h database for more details.`)
     // rootCmd.PersistentFlags().StringVarP(&_dsn, "datasource", "s", "", `datasource url for connecting to a database, run hare -h database for more  details.`)
     // rootCmd.Flags().StringVarP(&_output, "output", "o", "", `output file name of generated results(use standard output by default)`)
     // rootCmd.Flags().StringArrayVarP(&_template, "template", "t", nil, `specify templates to use`)
     // rootCmd.Flags().StringVarP(&_ldelim, "ldelim", "l", "", `specify left delimiters for template definition, must be used with rdelim together`)
     // rootCmd.Flags().StringVarP(&_rdelim, "rdelim", "r", "", `specify right delimiters for template definition, must be used with ldelim together`)
     // rootCmd.PersistentFlags().StringArrayVarP(&_var, "var", "v", nil, `var=value`)
 }
 
 // Execute cmd.
 func Execute() error {
     return rootCmd.Execute()
 }
 
 // For testing purpose.
 func GetRootCmd() *cobra.Command {
     return rootCmd
 }