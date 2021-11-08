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
	"GoMapEnum/src/logger"
	searchengine "GoMapEnum/src/searchEngine"

	"github.com/spf13/cobra"
)

var searchEngineOptions searchengine.Options

// o365Cmd represents the o365 command
var searchEngineCmd = &cobra.Command{
	Use:     "searchEngine",
	Short:   "Search the company name on Bing and Google and parse the Linkedin results",
	Long:    `Credits: https://github.com/m8r0wn/CrossLinked`,
	Example: `go run main.go gather searchEngine -c contoso -f "{f}{last}@contonso.com" -v`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "searchEngine", "Google and Bing search engine")
		log.SetLevel(level)
		log.Info("Starting the module searchEngine")
		searchEngineOptions.Log = log
		searchEngineOptions.Proxy = proxy
		users = searchEngineOptions.Gather()
	},
}

func init() {

	searchEngineCmd.Flags().StringVarP(&searchEngineOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	searchEngineCmd.Flags().StringVarP(&searchEngineOptions.Company, "company", "c", "", "Company name")
	searchEngineCmd.Flags().BoolVarP(&searchEngineOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	searchEngineCmd.MarkFlagRequired("company")
	searchEngineCmd.MarkFlagRequired("format")
}
