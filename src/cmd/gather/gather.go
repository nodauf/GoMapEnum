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
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var proxy func(*http.Request) (*url.URL, error)
var users []string
var output string
var proxyString string

// GatherCmd represents the gather command
var GatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Retrieve a list of email address based on the company name",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(strings.Join(users, "\n")), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	GatherCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	GatherCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	GatherCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	GatherCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Sleep in seconds before sending an authentication request")

	GatherCmd.AddCommand(searchEngineCmd)
	GatherCmd.AddCommand(linkedinCmd)
}

func initLogger() {
	if debug {
		level = logger.DebugLevel
	} else if verbose {
		level = logger.VerboseLevel
	} else {
		level = logger.InfoLevel
	}

}

func initProxy() {
	if proxyString != "" {
		url, err := url.Parse(proxyString)
		if err != nil {
			log.Fatalln("Fail to parse URL " + proxyString + " - error " + err.Error())
		}
		proxy = http.ProxyURL(url)
	}
}
