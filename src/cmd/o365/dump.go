package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

// dumpCmd represents the o365 command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Authenticate on o365 and dump various informations",
	Long: `
	Credits: https://github.com/dirkjanm/ROADtools`,
	Example: `go run main.go o365 dump  -u john.doe@contoso.com  -p s3crur3P4ssw0rd
go run main.go o365 dump  -u john.doe@contoso.com  -p s3crur3P4ssw0rd --dump users --html=False --json`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// check the objects to dump
		if o365Options.DumpObjects != "all" {
			dumpObjectsSlice := strings.Split(o365Options.DumpObjects, ",")
			for _, object := range dumpObjectsSlice {
				if !o365.IsObjectCanBeDumped(object) {
					return fmt.Errorf("cannot dump %s", object)
				}
			}
		}
		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Dump", "O365", "https://login.microsoftonline.com")
		log.SetLevel(level)
		log.Info("Starting the module O365")
		o365Options.Log = log

		if o365Options.CheckIfValid {
			optionsInterface := reflect.ValueOf(&o365Options).Interface()
			if !o365.Authenticate(&optionsInterface, o365Options.Users, o365Options.Passwords) {
				log.Error("Cannot authenticate with %s / %s", o365Options.Users, o365Options.Passwords)
				return
			}
		}

		o365Options.Dump()
	},
}

func init() {

	dumpCmd.Flags().BoolVarP(&o365Options.CheckIfValid, "check", "c", true, "Check if the user is valid before trying password")
	dumpCmd.Flags().StringVarP(&o365Options.Users, "user", "u", "", "User or file containing the emails")
	dumpCmd.Flags().StringVarP(&o365Options.Passwords, "password", "p", "", "Password or file containing the passwords")
	dumpCmd.Flags().StringVar(&o365Options.DumpObjects, "dump", "all", "Dump objects among computers,users. Could be 'all' keyword  (ex: tenantDetails,policies,servicePrincipals,groups,applications,devices,directoryRoles,roleDefinitions,contacts,oauth2PermissionGrants)")
	dumpCmd.Flags().IntVarP(&o365Options.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	dumpCmd.Flags().IntVar(&o365Options.Thread, "thread", 2, "Number of threads ")
	dumpCmd.Flags().BoolVar(&o365Options.HTML, "html", true, "save the output as html")
	dumpCmd.Flags().BoolVar(&o365Options.JSON, "json", false, "save the output as JSON (all data are kept)")
	dumpCmd.MarkFlagRequired("user")
	dumpCmd.MarkFlagRequired("password")
}
