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
	"GoMapEnum/src/owa"

	"github.com/spf13/cobra"
)

var owaOptions owa.Options

// owaCmd represents the owa command
var owaCmd = &cobra.Command{
	Use:   "owa",
	Short: "Enumerate account on a Exchange",
	Long: `The response for invalid user will be significantly longer than for valid account.
An average response time is calculated, and each attempt is then compare to the average response time.
Beware of account locking. No locking information is returned and therefore no failsafes could be set up.
Credits: https://github.com/busterb/msmailprobe`,
	Example: `go run main.go userenum owa -u users  -t mail.contoso.com -o validUsers`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Enumeration", "OWA", owaOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module OWA")
		owaOptions.Log = log
		owaOptions.Proxy = proxy
		validUsers = owaOptions.UserEnum()

	},
}

func init() {

	owaCmd.Flags().StringVarP(&owaOptions.Users, "user", "u", "", "Username or file containing the usernames")
	owaCmd.Flags().StringVarP(&owaOptions.Target, "target", "t", "", "Host pointing to the OWA service")
	owaCmd.Flags().IntVar(&owaOptions.Thread, "thread", 2, "Number of threads")
	owaCmd.MarkFlagRequired("user")
	owaCmd.MarkFlagRequired("target")
}
