package brute

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/adfs"
	"GoMapEnum/src/orchestrator"
	"errors"

	"github.com/spf13/cobra"
)

var adfsOptions adfs.Options

// adfsCmd represents the adfs command
var adfsCmd = &cobra.Command{
	Use:   "adfs",
	Short: "Bruteforce or spray password on a ADFS instance",
	Long: `Authenticate on https://<target>/adfs/ls/idpinitiatedsignon.aspx?client-request-id=<randomGUID>&pullStatus=0.
The hostname of the target can be discovered using https://login.microsoftonline.com/getuserrealm.srf?login=<company domain>
Beware of account locking. No locking information is returned and therefore no failsafes could be set up.`,
	Example: `go run main.go bruteSpray adfs -d contoso.com -u users  -p pass
go run main.go bruteSpray adfs -t adfs.contoso.com -u users  -p pass -s 5 -o validUsers
go run main.go bruteSpray adfs -t adfs.contoso.com -u john.doe@contoso.com  -p Autumn2021! -s 5 -v`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if adfsOptions.Target == "" && adfsOptions.Domain == "" {
			return errors.New("either Target or Domain option should be specified")
		}
		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "ADFS", adfsOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module ADFS")
		adfsOptions.Log = log
		adfsOptions.NoBruteforce = noBruteforce
		adfsOptions.Sleep = sleep
		adfsOptions.Proxy = proxy

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionBruteforce = adfs.CheckTarget
		orchestratorOptions.AuthenticationFunc = adfs.Authenticate
		validUsers = orchestratorOptions.Bruteforce(&adfsOptions)
	},
}

func init() {

	adfsCmd.Flags().StringVarP(&adfsOptions.Users, "user", "u", "", "User or file containing the emails")
	adfsCmd.Flags().StringVarP(&adfsOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	adfsCmd.Flags().StringVarP(&adfsOptions.Target, "target", "t", "", "Host pointing to the ADFS service (if not specified the ADFS will be guess on login.microsoftonline.com)")
	adfsCmd.Flags().StringVarP(&adfsOptions.Domain, "domain", "d", "", "If the target is not specified, the domain will be used to guess the ADFS url")
	adfsCmd.Flags().IntVar(&adfsOptions.Thread, "thread", 2, "Number of threads ")
	adfsCmd.MarkFlagRequired("user")
	adfsCmd.MarkFlagRequired("password")

}
