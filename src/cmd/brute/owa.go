package brute

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/owa"

	"github.com/spf13/cobra"
)

var owaOptions owa.Options

// owaCmd represents the owa command
var owaCmd = &cobra.Command{
	Use:   "owa",
	Short: "Authentication on Exchange",
	Long: `Authenticate with basic authentication on multiple endpoints. 
Beware of account locking. No locking information is returned and therefore no failsafes could be set up.
Credits: https://github.com/busterb/msmailprobe`,
	Example: `go run main.go bruteSpray owa -u users -p pass -t mail.contoso.com -s 10 -o validUsers
go run main.go bruteSpray owa -u john.doe@contoso.com -p Automn2021! -t mail.contoso.com -v`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "OWA", owaOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module OWA")
		owaOptions.Log = log
		owaOptions.Proxy = proxy
		owaOptions.NoBruteforce = noBruteforce
		owaOptions.Sleep = sleep
		owaOptions.Brute()

	},
}

func init() {

	owaCmd.Flags().StringVarP(&owaOptions.Users, "user", "u", "", "User or file containing the emails")
	owaCmd.Flags().StringVarP(&owaOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	owaCmd.Flags().StringVarP(&owaOptions.Target, "target", "t", "", "Host pointing to the OWA service")
	owaCmd.Flags().IntVar(&owaOptions.Thread, "thread", 2, "Number of threads ")
	owaCmd.MarkFlagRequired("user")
	owaCmd.MarkFlagRequired("password")
}
