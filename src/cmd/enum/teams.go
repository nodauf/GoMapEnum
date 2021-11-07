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
	"GoMapEnum/src/logger"
	"GoMapEnum/src/teams"

	"github.com/spf13/cobra"
)

var teamsOptions teams.Options

// teamsCmd represents the teams command
var teamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "User enumeration on Microsoft Teams (Stealthier)",
	Long: `Users can be enumerated on Microsoft Teams with the search features.
it will validates an email address or a list of email addresses.
If these emails exist the presence of the user is retrieved as well as the device used to connect`,
	Example: `go run main.go userenum teams -t "eyJ0..." -e emails -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "Teams", owaOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module Teams")
		teamsOptions.Log = log
		teamsOptions.Proxy = proxy
		validUsers = teamsOptions.UserEnum(log)
	},
}

func init() {

	teamsCmd.Flags().StringVarP(&teamsOptions.Email, "email", "e", "", "Email or file containing the email address")
	teamsCmd.Flags().StringVarP(&teamsOptions.Token, "token", "t", "", "Bearer token (only the base64 part: eyJ0...)")
	teamsCmd.Flags().IntVar(&teamsOptions.Thread, "thread", 1, "Number of threads")

	teamsCmd.MarkFlagRequired("token")
	teamsCmd.MarkFlagRequired("email")
}
