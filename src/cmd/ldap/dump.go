package ldap

import (
	"GoMapEnum/src/logger"
	"errors"

	"github.com/spf13/cobra"
)

// dumpCmd represents the owa command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump LDAP datas",
	Long:  `Several LDAP requests has been implemented to dump various of data`,
	Example: `go run main.go ldap dump -t 192.168.1.175 -u user -p securePassword --dump computers,users
	go run main.go ldap dump -t 192.168.1.175 -u user -p securePassword --dump all`,
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
	dumpCmd.Flags().IntVarP(&ldapOptions.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")

}
