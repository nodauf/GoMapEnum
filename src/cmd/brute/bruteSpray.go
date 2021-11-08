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
package brute

import (
	"GoMapEnum/src/logger"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var noBruteforce bool
var sleep int
var proxy func(*http.Request) (*url.URL, error)

var proxyString string

// BruteSprayCmd represents the bruteSpray command
var BruteSprayCmd = &cobra.Command{
	Use:   "bruteSpray",
	Short: "Spray a password or bruteforce a user's password",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	BruteSprayCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	BruteSprayCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	BruteSprayCmd.PersistentFlags().BoolVarP(&noBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	BruteSprayCmd.PersistentFlags().IntVarP(&sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	BruteSprayCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Sleep in seconds before sending an authentication request")

	BruteSprayCmd.AddCommand(o365Cmd)
	BruteSprayCmd.AddCommand(adfsCmd)
	BruteSprayCmd.AddCommand(owaCmd)
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
			fmt.Println("Fail to parse URL " + proxyString + " - error " + err.Error())
			os.Exit(1)
		}
		proxy = http.ProxyURL(url)
	}
}
