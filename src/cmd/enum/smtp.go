package enum

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smtp"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

var smtpOptions smtp.Options

// smtpCmd represents the smtp command
var smtpCmd = &cobra.Command{
	Use:   "smtp",
	Short: "Enumerate email address by connection to the smtp port of the target.",
	Long: `SMTP user enumeration with RCPT, VRFY and EXPN.
	Credits: https://github.com/cytopia/smtp-user-enum`,
	Example: `go run main.go userenum smtp -u users  -t mail.contoso.com -d domain.tld -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "SMTP", smtpOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module SMTP")
		smtpOptions.Log = log
		smtpOptions.Proxy = proxy

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionUserEnum = smtp.PrepareSMTPConnections
		orchestratorOptions.UserEnumFunc = smtp.UserEnum
		orchestratorOptions.PostActionUserEnum = smtp.CloseSMTPConnections
		validUsers = orchestratorOptions.UserEnum(&smtpOptions)

	},
}

func init() {

	smtpCmd.Flags().StringVarP(&smtpOptions.Domain, "domain", "d", "", "Targeted domain ")
	smtpCmd.Flags().StringVarP(&smtpOptions.Mode, "mode", "m", "", "RCPT, VRFY, EXPN (default: RCPT)")
	smtpCmd.Flags().StringVarP(&smtpOptions.Users, "user", "u", "", "Username or file containing the usernames")
	smtpCmd.Flags().StringVarP(&smtpOptions.Target, "target", "t", "", "Host pointing to the SMTP service. If not specified, the first SMTP server in the MX record will be targeted.")
	smtpCmd.Flags().IntVar(&smtpOptions.Thread, "thread", 2, "Number of threads")
	smtpCmd.MarkFlagRequired("user")
	smtpCmd.MarkFlagRequired("domain")
}
