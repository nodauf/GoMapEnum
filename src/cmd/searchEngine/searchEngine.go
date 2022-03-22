package searchEngine

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/searchEngine"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var validUsers string
var output string
var proxyString string

var searchEngineOptions searchEngine.Options

// SearchEngineCmd represents the bruteSpray command
var SearchEngineCmd = &cobra.Command{
	Use:   "searchEngine",
	Short: "Commands for searchEngine module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(validUsers), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	SearchEngineCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	SearchEngineCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	SearchEngineCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	SearchEngineCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	SearchEngineCmd.AddCommand(gatherCmd)
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
		searchEngineOptions.ProxyHTTP = http.ProxyURL(url)
	}
}
