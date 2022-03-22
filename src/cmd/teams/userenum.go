package teams

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/teams"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// enumCmd represents the teams command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "User enumeration on Microsoft Teams (Stealthier)",
	Long: `Users can be enumerated on Microsoft Teams with the search features.
it will validates an email address or a list of email addresses.
If these emails exist the presence of the user is retrieved as well as the device used to connect`,
	Example: `go run main.go userenum teams -t "eyJ0..." -e emails -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "Teams", teamsOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Teams")
		teamsOptions.Token = "Bearer " + teamsOptions.Token
		teamsOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.UserEnumFunc = teams.UserEnum
		validUsers = orchestratorOptions.UserEnum(&teamsOptions)
	},
}

func init() {

	enumCmd.Flags().StringVarP(&teamsOptions.Users, "user", "u", "", "Email or file containing the email address")
	enumCmd.Flags().StringVarP(&teamsOptions.Token, "token", "t", "", "Bearer token (only the base64 part: eyJ0...). This token can be found on requests made to teams.microsoft.com/api/")
	enumCmd.Flags().IntVar(&teamsOptions.Thread, "thread", 1, "Number of threads")

	enumCmd.MarkFlagRequired("token")
	enumCmd.MarkFlagRequired("user")
}
