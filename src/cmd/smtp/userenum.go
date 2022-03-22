package smtp

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smtp"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// enum represents the smtp command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "Enumerate email address by connection to the smtp port of the target.",
	Long: `SMTP user enumeration with RCPT, VRFY and EXPN.
	Credits: https://github.com/cytopia/smtp-user-enum`,
	Example: `go run main.go userenum smtp -u users  -t mail.contoso.com -d domain.tld -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "SMTP", smtpOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module SMTP")
		smtpOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionUserEnum = smtp.PrepareSMTPConnections
		orchestratorOptions.UserEnumFunc = smtp.UserEnum
		orchestratorOptions.PostActionUserEnum = smtp.CloseSMTPConnections
		validUsers = orchestratorOptions.UserEnum(&smtpOptions)

	},
}

func init() {

	enumCmd.Flags().StringVarP(&smtpOptions.Domain, "domain", "d", "", "Targeted domain ")
	enumCmd.Flags().StringVarP(&smtpOptions.Mode, "mode", "m", "", "RCPT, VRFY, EXPN, ALL (default: all)")
	enumCmd.Flags().StringVarP(&smtpOptions.Users, "user", "u", "", "Username or file containing the usernames")
	enumCmd.Flags().StringVarP(&smtpOptions.Target, "target", "t", "", "Host pointing to the SMTP service. If not specified, the first SMTP server in the MX record will be targeted.")
	enumCmd.Flags().IntVar(&smtpOptions.Thread, "thread", 2, "Number of threads")
	enumCmd.MarkFlagRequired("user")
}
