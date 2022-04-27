package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/kerberos"
	"GoMapEnum/src/orchestrator"
	"errors"
	"fmt"
	"os"

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
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if kerberosOptions.Passwords != "" && kerberosOptions.Hash != "" {
			return errors.New("Only password or hash field should be specified")
		}

		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "kerberos", kerberosOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Kerberos")
		kerberosOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionUserEnum = kerberos.KerberosSession
		orchestratorOptions.UserEnumFunc = kerberos.UserEnum
		validUsers = orchestratorOptions.UserEnum(&kerberosOptions)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(validUsers), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	enumCmd.Flags().StringVarP(&kerberosOptions.Users, "user", "u", "", "User or file containing the emails")
	enumCmd.Flags().IntVar(&kerberosOptions.Thread, "thread", 2, "Number of threads")
	enumCmd.MarkFlagRequired("user")
}
