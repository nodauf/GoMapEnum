package owa

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/owa"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// enumCmd represents the owa command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "Enumerate account on a Exchange",
	Long: `The response for invalid user will be significantly longer than for valid account.
An average response time is calculated, and each attempt is then compare to the average response time.
Beware of account locking. No locking information is returned and therefore no failsafes could be set up.
Credits: https://github.com/busterb/msmailprobe`,
	Example: `go run main.go owa userenum -u users  -t mail.contoso.com -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("User enumeration", "OWA", owaOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module OWA")
		owaOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionUserEnum = owa.InitAndAverageResponseTime
		orchestratorOptions.UserEnumFunc = owa.UserEnum
		validUsers = orchestratorOptions.UserEnum(&owaOptions)

	},
}

func init() {

	enumCmd.Flags().StringVarP(&owaOptions.Users, "user", "u", "", "Username or file containing the usernames")
	enumCmd.Flags().StringVarP(&owaOptions.Target, "target", "t", "", "Host pointing to the OWA service")
	enumCmd.Flags().IntVar(&owaOptions.Thread, "thread", 2, "Number of threads")
	enumCmd.MarkFlagRequired("user")
	enumCmd.MarkFlagRequired("target")
}
