package azure

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/azure"
	"GoMapEnum/src/orchestrator"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var azureOptions azure.Options

// enumCmd represents the azure command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "User enumeration through autlogon API",
	Long: `The authentication process does not seem to work but the error code can still give information if the user's account exist or not
	Credits https://github.com/treebuilder/aad-sso-enum-brute-spray`,
	Example: `go run main.go azure userenum -u john.doe@contoso.com
	go run main.go userenum azure -u users -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "Azure", "https://autologon.microsoftazuread-sso.com")
		log.SetLevel(level)
		log.Info("Starting the module Azure")
		azureOptions.Log = log
		azureOptions.Proxy = proxy

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.UserEnumFunc = azure.UserEnum
		validUsers = orchestratorOptions.UserEnum(&azureOptions)
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(strings.Join(validUsers, "\n")), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	enumCmd.Flags().StringVarP(&azureOptions.Users, "user", "u", "", "User or file containing the emails")
	enumCmd.Flags().IntVar(&azureOptions.Thread, "thread", 2, "Number of threads")
	enumCmd.MarkFlagRequired("user")
}
