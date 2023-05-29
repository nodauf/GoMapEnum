package owa

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/owa"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// bruteCmd represents the owa command
var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "Authentication on Exchange",
	Long: `Authenticate with basic authentication on multiple endpoints. 
Beware of account locking. No locking information is returned and therefore no failsafes could be set up.
Credits: https://github.com/busterb/msmailprobe`,
	Example: `go run main.go owa brute -u users -p pass -t mail.contoso.com -s 10 -o validUsers
go run main.go owa brute -u john.doe@contoso.com -p Automn2021! -t mail.contoso.com -v`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "OWA", owaOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module OWA")
		owaOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionBruteforce = owa.PrepareBruteforce
		orchestratorOptions.CustomOptionsForCheckIfValid = owa.PrepareOptions
		orchestratorOptions.PreActionUserEnum = owa.InitAndAverageResponseTime
		orchestratorOptions.UserEnumFunc = owa.UserEnum
		validUsers = orchestratorOptions.Bruteforce(&owaOptions)

	},
}

func init() {

	bruteCmd.Flags().BoolVarP(&owaOptions.CheckIfValid, "check", "c", true, "Check if the user is valid before trying password")
	bruteCmd.Flags().BoolVarP(&owaOptions.NoBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	bruteCmd.Flags().BoolVar(&owaOptions.Basic, "basic", false, "Basic authentication instead of NTLM")
	bruteCmd.Flags().StringVarP(&owaOptions.Users, "user", "u", "", "User or file containing the emails")
	bruteCmd.Flags().StringVarP(&owaOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	bruteCmd.Flags().StringVarP(&owaOptions.Target, "target", "t", "", "Host pointing to the OWA service")
	bruteCmd.Flags().IntVar(&owaOptions.Thread, "thread", 2, "Number of threads")
	bruteCmd.Flags().IntVarP(&owaOptions.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	bruteCmd.MarkFlagRequired("user")
	bruteCmd.MarkFlagRequired("password")
	bruteCmd.MarkFlagRequired("target")
}
