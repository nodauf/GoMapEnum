package cmd

import (
	"GoMapEnum/src/cmd/adfs"
	"GoMapEnum/src/cmd/azure"
	"GoMapEnum/src/cmd/ldap"
	"GoMapEnum/src/cmd/linkedin"
	"GoMapEnum/src/cmd/o365"
	"GoMapEnum/src/cmd/owa"
	"GoMapEnum/src/cmd/searchEngine"
	"GoMapEnum/src/cmd/smb"
	"GoMapEnum/src/cmd/smtp"
	"GoMapEnum/src/cmd/teams"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "GoMapEnum",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Add child
	//rootCmd.AddCommand(enum.UserenumCmd)
	//rootCmd.AddCommand(brute.BruteSprayCmd)
	//rootCmd.AddCommand(gather.GatherCmd)
	// New module
	rootCmd.AddCommand(azure.AzureCmd)
	rootCmd.AddCommand(adfs.AdfsCmd)
	rootCmd.AddCommand(ldap.LdapCmd)
	rootCmd.AddCommand(linkedin.LinkedinCmd)
	rootCmd.AddCommand(o365.O365Cmd)
	rootCmd.AddCommand(owa.OwaCmd)
	rootCmd.AddCommand(searchEngine.SearchEngineCmd)
	rootCmd.AddCommand(smb.SmbCmd)
	rootCmd.AddCommand(smtp.SMTPCmd)
	rootCmd.AddCommand(teams.TeamsCmd)
}
