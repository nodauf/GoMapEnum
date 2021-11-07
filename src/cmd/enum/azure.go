/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package enum

import (
	"GoMapEnum/src/azure"
	"GoMapEnum/src/logger"

	"github.com/spf13/cobra"
)

var azureOptions azure.Options

// azureCmd represents the azure command
var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "User enumeration through autlogon API",
	Long: `The authentication process does not seem to work but the error code can still give information if the user's account exist or not
	Credits https://github.com/treebuilder/aad-sso-enum-brute-spray`,
	Example: `go run main.go userenum azure -u john.doe@contoso.com
	go run main.go userenum azure -u users -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "Azure", "https://autologon.microsoftazuread-sso.com")
		log.SetLevel(level)
		log.Info("Starting the module Azure")
		azureOptions.Log = log
		azureOptions.Proxy = proxy
		validUsers = azureOptions.UserEnum()
	},
}

func init() {

	azureCmd.Flags().StringVarP(&azureOptions.Users, "user", "u", "", "User or file containing the emails")
	azureCmd.Flags().IntVar(&azureOptions.Thread, "thread", 2, "Number of threads")
	azureCmd.MarkFlagRequired("user")
}
