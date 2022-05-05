package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
	"GoMapEnum/src/orchestrator"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// enumCmd represents the azure command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "User enumeration through autlogon API",
	Long: `The authentication process does not seem to work but the error code can still give information if the user's account exist or not
	Credits https://github.com/treebuilder/aad-sso-enum-brute-spray`,
	Example: `go run main.go azure userenum -u john.doe@contoso.com
	go run main.go userenum azure -u users -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "o365", "https://login.microsoftonline.com")
		log.SetLevel(level)
		log.Info("Starting the module Azure")
		o365Options.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.UserEnumFunc = o365.UserEnum
		validUsers = orchestratorOptions.UserEnum(&o365Options)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(validUsers), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	enumCmd.Flags().StringVarP(&o365Options.Users, "user", "u", "", "User or file containing the emails")
	enumCmd.Flags().IntVar(&o365Options.Thread, "thread", 2, "Number of threads")
	enumCmd.MarkFlagRequired("user")
}
