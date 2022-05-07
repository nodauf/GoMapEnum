package kerberos

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/kerberos"
	"GoMapEnum/src/orchestrator"
	"errors"

	"github.com/spf13/cobra"
)

// Credits: https://github.com/ropnop/kerbrute

// bruteCmd represents the azure command
var bruteCmd = &cobra.Command{
	Use:     "brute",
	Short:   "",
	Long:    ``,
	Example: ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if (kerberosOptions.Passwords != "" && kerberosOptions.Hash != "") || (kerberosOptions.Passwords == "" && kerberosOptions.Hash == "") {
			return errors.New("Only password or hash field should be specified")
		}

		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "kerberos", kerberosOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Kerberos")
		kerberosOptions.Log = log

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.PreActionBruteforce = kerberos.KerberosSession
		orchestratorOptions.AuthenticationFunc = kerberos.Authenticate
		orchestratorOptions.UserEnumFunc = kerberos.UserEnum
		orchestratorOptions.CustomOptionsForCheckIfValid = kerberos.PrepareOptions
		outputCmd = orchestratorOptions.Bruteforce(&kerberosOptions)
	},
}

func init() {

	bruteCmd.Flags().BoolVarP(&kerberosOptions.CheckIfValid, "check", "c", true, "Check if the user is valid before trying password")
	bruteCmd.Flags().StringVarP(&kerberosOptions.Users, "user", "u", "", "User or file containing the usernames")
	bruteCmd.Flags().StringVarP(&kerberosOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	bruteCmd.Flags().StringVarP(&kerberosOptions.Hash, "hash", "H", "", "Hash or file containing the hashes (not implemented yet)")
	bruteCmd.Flags().BoolVar(&kerberosOptions.StopOnLockout, "stopOnLockout", true, "Stop the bruteforce if an account is locked out")
	bruteCmd.MarkFlagRequired("user")
}
