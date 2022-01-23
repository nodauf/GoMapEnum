package enum

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
	"GoMapEnum/src/orchestrator"
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

var o365Options o365.Options

// o365Cmd represents the o365 command
var o365Cmd = &cobra.Command{
	Use:   "o365",
	Short: "Enumerate users on serveral o365 endpoint. One of them does not require authenticate and it is therefore stealthier.",
	Long: `Can enumerate users without authenticate on office mode. On onedrive and oauth2, an authentication attempt will be made (not implemented yet).
	Credits: https://github.com/0xZDH/o365spray`,
	Example: `go run main.go userenum o365  -u users`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		o365Options.Mode = strings.ToLower(o365Options.Mode)
		if o365Options.Mode != "office" && o365Options.Mode != "oauth2" && o365Options.Mode != "onedrive" {
			return errors.New("invalid mode. Should be office, oauth2 or onedrive")
		}
		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "O365", "https://login.microsoftonline.com")
		log.SetLevel(level)
		log.Info("Starting the module O365")
		o365Options.Log = log
		o365Options.Proxy = proxy
		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.CheckBeforeEnumFunc = o365.CheckTenant
		orchestratorOptions.UserEnumFunc = o365.UserEnum
		validUsers = orchestratorOptions.UserEnum(&o365Options)
	},
}

func init() {

	o365Cmd.Flags().StringVarP(&o365Options.Mode, "mode", "m", "office", "Choose a mode between office (Stealthier), oauth2 and onedrive")
	o365Cmd.Flags().StringVarP(&o365Options.Users, "user", "u", "", "User or file containing the emails")
	o365Cmd.Flags().IntVar(&o365Options.Thread, "thread", 2, "Number of threads")
	o365Cmd.MarkFlagRequired("user")
}
