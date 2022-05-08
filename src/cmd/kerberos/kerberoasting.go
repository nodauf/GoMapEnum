package kerberos

import (
	"GoMapEnum/src/logger"
	"errors"

	"github.com/spf13/cobra"
)

var targetUser string

// kerberoastingCmd represents the azure command
var kerberoastingCmd = &cobra.Command{
	Use:     "kerberoasting",
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
		log := logger.New("Kerberoasting", "kerberos", kerberosOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Kerberos")
		kerberosOptions.Log = log

		outputCmd = kerberosOptions.Kerberoasting(targetUser, "")
	},
}

func init() {

	kerberoastingCmd.Flags().StringVarP(&kerberosOptions.Users, "user", "u", "", "Username")
	kerberoastingCmd.Flags().StringVar(&targetUser, "targetuser", "", "Target username, if not set all users with SPN will be targeted")
	kerberoastingCmd.Flags().StringVarP(&kerberosOptions.Passwords, "password", "p", "", "Password")
	kerberoastingCmd.Flags().StringVarP(&kerberosOptions.Hash, "hash", "H", "", "Hash (not implemented yet)")
	kerberoastingCmd.MarkFlagRequired("user")
}
