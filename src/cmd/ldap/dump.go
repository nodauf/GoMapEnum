package ldap

import (
	"GoMapEnum/src/logger"
	"errors"

	ldaplib "github.com/go-ldap/ldap/v3"
	"github.com/spf13/cobra"
)

// dumpCmd represents the owa command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump LDAP datas",
	Long:  ``,
	Example: `go run main.go owa brute -u users -p pass -t mail.contoso.com -s 10 -o validUsers
go run main.go owa brute -u john.doe@contoso.com -p Automn2021! -t mail.contoso.com -v`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if ldapOptions.Passwords == "" && ldapOptions.Hash == "" {
			return errors.New("The field password or hash is required")
		} else if ldapOptions.Passwords != "" && ldapOptions.Hash != "" {
			return errors.New("Only password or hash field should be specified")
		}

		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Dump", "LDAP", ldapOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module LDAP")
		ldapOptions.Log = log

		// If we bruteforce with hash, password will contains the hashs as we use the orchestrator and we want to keep the same format than CME
		if ldapOptions.Hash != "" {
			ldapOptions.Passwords = ldapOptions.Hash
			ldapOptions.IsHash = true
		}

		outputCmd = ldapOptions.Dump()

	},
}

func init() {
	dumpCmd.Flags().StringVarP(&ldapOptions.BaseDN, "baseDN", "b", "", "The base DN for all ldap queries. Default it will be the defaultNamingContext")
	dumpCmd.Flags().StringVar(&ldapOptions.DumpObjects, "dump", "", "Dump objects among computers,users. Could be 'all' keyword  (ex: computers,users)")
	dumpCmd.Flags().IntVar(&ldapOptions.Thread, "thread", 2, "Number of threads")
	dumpCmd.Flags().IntVar(&ldapOptions.Timeout, "timeout", int(ldaplib.DefaultTimeout.Seconds()), "Timeout for the SMB connection in seconds")
	dumpCmd.Flags().IntVarP(&ldapOptions.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")

}
