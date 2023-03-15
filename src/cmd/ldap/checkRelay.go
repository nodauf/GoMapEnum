package ldap

import (
	"GoMapEnum/src/logger"
	"github.com/spf13/cobra"
)

// bruteCmd represents the owa command
var checkRelayCmd = &cobra.Command{
	Use:     "checkRelay",
	Short:   "Check if signing and binding are required",
	Long:    ``,
	Example: ``,

	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("checkRelay", "LDAP", ldapOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module LDAP")
		ldapOptions.Log = log

		if ldapOptions.Hash != "" {
			ldapOptions.Passwords = ldapOptions.Hash
			ldapOptions.IsHash = true
		}

		ldapOptions.CheckRelay()

	},
}

func init() {

}
