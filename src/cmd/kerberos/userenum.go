package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/kerberos"
	"GoMapEnum/src/orchestrator"

	"github.com/spf13/cobra"
)

// Credits: https://github.com/ropnop/kerbrute

// enumCmd represents the azure command
var enumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "Enumerate valid domain usernames via Kerberos",
	Long: `Will enumerate valid usernames from a list by constructing AS-REQs to requesting a TGT from the KDC.
	If no domain controller is specified, the tool will attempt to look one up via DNS SRV records.
	A full domain is required. This domain will be capitalized and used as the Kerberos realm when attempting the bruteforce.
	Valid usernames will be displayed on stdout.
	It will not increase the badPwdCount counter.`,
	Example: ``,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "kerberos", kerberosOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Kerberos")
		kerberosOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionUserEnum = kerberos.InitSession
		orchestratorOptions.UserEnumFunc = kerberos.UserEnum
		outputCmd = orchestratorOptions.UserEnum(&kerberosOptions)
	},
}

func init() {

	enumCmd.Flags().StringVarP(&kerberosOptions.Users, "user", "u", "", "User or file containing the emails")
	enumCmd.MarkFlagRequired("user")
}
