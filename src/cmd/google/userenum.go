package google

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/google"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// enumCmd represents the teams command
var enumCmd = &cobra.Command{
	Use:     "userenum",
	Short:   "User enumeration on Google Mail (Stealthier)",
	Long:    ``,
	Example: `go run main.go google userenum -e emails -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "Google", "https://mail.google.com/mail/gxlu?email")
		log.SetLevel(level)
		log.Info("Starting the module Google")
		googleOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.UserEnumFunc = google.UserEnum
		validUsers = orchestratorOptions.UserEnum(&googleOptions)
	},
}

func init() {

	enumCmd.Flags().StringVarP(&googleOptions.Users, "user", "u", "", "Email or file containing the email address")
	enumCmd.Flags().IntVar(&googleOptions.Thread, "thread", 1, "Number of threads")

	enumCmd.MarkFlagRequired("user")
}
