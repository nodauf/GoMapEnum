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
package gather

import (
	"GoMapEnum/src/linkedin"
	"GoMapEnum/src/logger"

	"github.com/spf13/cobra"
)

var linkedinOptions linkedin.Options

// o365Cmd represents the o365 command
var linkedinCmd = &cobra.Command{
	Use:   "linkedin",
	Short: "Search on Linkedin for people working in the specified company",
	Long: `Firstly, it will search for company based on the provided name and then list all the people working at these companies and print them in the specified format.
The session cookie is needed to use the Linkedin features.`,
	Example: `go run main.go gather linkedin -c contoso -f "{f}{last}@contonso.com" -e -s AQEDA...`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "linkedin", "Linkedin")
		log.SetLevel(level)
		log.Info("Starting the module linkedin")
		linkedinOptions.Log = log
		linkedinOptions.Proxy = proxy
		users = linkedinOptions.Gather()
	},
}

func init() {

	linkedinCmd.Flags().StringVarP(&linkedinOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	linkedinCmd.Flags().StringVarP(&linkedinOptions.Company, "company", "c", "", "Company name")
	linkedinCmd.Flags().BoolVarP(&linkedinOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	linkedinCmd.Flags().StringVarP(&linkedinOptions.Cookie, "cookie", "s", "", "Session cookie named li_at")
	linkedinCmd.MarkFlagRequired("company")
	linkedinCmd.MarkFlagRequired("cookie")
	linkedinCmd.MarkFlagRequired("format")
}
