package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
	"GoMapEnum/src/orchestrator"
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

// bruteCmd represents the o365 command
var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "Authenticate on multiple endpoint of o365 (lockout detection available)",
	Long: `Authenticate on three different o365 endpoint: oauth2 or onedrive (not yet implemented).
Beware of account locking. Locking information is only available on oauth2 and therefore failsafe is only set up on oauth2.
By default, if one account is being lock, the all attack will be stopped.
	Credits: https://github.com/0xZDH/o365spray`,
	Example: `go run main.go o365 bruteSpray  -u john.doe@contoso.com  -p passwordFile -s 10 -l 2`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		o365Options.Mode = strings.ToLower(o365Options.Mode)
		if o365Options.Mode != "oauth2" && o365Options.Mode != "autodiscover" {
			return errors.New("invalid mode. Should be oauth2 or autodiscover")
		}
		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "O365", "https://login.microsoftonline.com")
		log.SetLevel(level)
		log.Info("Starting the module O365")
		o365Options.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.CustomOptionsForCheckIfValid = o365.PrepareOptions
		orchestratorOptions.AuthenticationFunc = o365.Authenticate
		orchestratorOptions.UserEnumFunc = o365.UserEnum
		// To check if the user is valid
		orchestratorOptions.CheckBeforeEnumFunc = o365.CheckTenant
		orchestratorOptions.AuthenticationFunc = o365.Authenticate
		validUsers = orchestratorOptions.Bruteforce(&o365Options)
	},
}

func init() {

	bruteCmd.Flags().BoolVarP(&o365Options.CheckIfValid, "check", "c", true, "Check if the user is valid before trying password")
	bruteCmd.Flags().BoolVarP(&o365Options.NoBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	bruteCmd.Flags().StringVarP(&o365Options.Mode, "mode", "m", "oauth2", "Choose a mode between oauth2 and autodiscover (no failsafe for lockout) <- not implemented")
	bruteCmd.Flags().StringVarP(&o365Options.Users, "user", "u", "", "User or file containing the emails")
	bruteCmd.Flags().StringVarP(&o365Options.Passwords, "password", "p", "", "Password or file containing the passwords")
	bruteCmd.Flags().IntVarP(&o365Options.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	bruteCmd.Flags().IntVarP(&o365Options.LockoutThreshold, "lockout-threshold", "l", 1, "Stop the bruteforce when the threshold is meet")
	bruteCmd.Flags().IntVar(&o365Options.Thread, "thread", 2, "Number of threads ")
	bruteCmd.MarkFlagRequired("user")
	bruteCmd.MarkFlagRequired("password")
}
