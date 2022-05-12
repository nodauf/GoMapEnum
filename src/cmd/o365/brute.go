package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
	"GoMapEnum/src/orchestrator"

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
	Example: `go run main.go o365 brute -u john.doe@contoso.com  -p passwordFile -s 10 --stopOnLockout=False`,
	/*PreRunE: func(cmd *cobra.Command, args []string) error {
		// Issue (https://github.com/spf13/cobra/issues/1047) with the default value of a flag that is used in two subcommand
		// Some ugly workaround here. If the mode is the default of the other command we switch to the default of this one: oauth2
		if o365Options.Mode == "office" {
			o365Options.Mode = "oauth2"
		}

		o365Options.Mode = strings.ToLower(o365Options.Mode)
		if o365Options.Mode != "oauth2" && o365Options.Mode != "autodiscover" {
			return errors.New(o365Options.Mode + " is an invalid mode. Should be oauth2 or autodiscover")
		}
		return nil
	},*/
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "O365", "https://login.microsoftonline.com")
		log.SetLevel(level)
		log.Info("Starting the module O365")
		o365Options.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.CustomOptionsForCheckIfValid = o365.PrepareOptions
		orchestratorOptions.AuthenticationFunc = o365.Authenticate
		orchestratorOptions.UserEnumFunc = o365.UserEnum
		validUsers = orchestratorOptions.Bruteforce(&o365Options)
	},
}

func init() {

	bruteCmd.Flags().BoolVarP(&o365Options.CheckIfValid, "check", "c", true, "Check if the user is valid before trying password")
	bruteCmd.Flags().BoolVarP(&o365Options.NoBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	//bruteCmd.Flags().StringVarP(&o365Options.Mode, "mode", "m", "oauth2", "oauth2")
	bruteCmd.Flags().StringVarP(&o365Options.Users, "user", "u", "", "User or file containing the emails")
	bruteCmd.Flags().StringVarP(&o365Options.Passwords, "password", "p", "", "Password or file containing the passwords")
	bruteCmd.Flags().IntVarP(&o365Options.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	bruteCmd.Flags().BoolVar(&o365Options.StopOnLockout, "stopOnLockout", true, "Stop the bruteforce if an account is locked out")
	bruteCmd.Flags().IntVar(&o365Options.Thread, "thread", 2, "Number of threads ")
	bruteCmd.MarkFlagRequired("user")
	bruteCmd.MarkFlagRequired("password")
}
