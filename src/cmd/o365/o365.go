package o365

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/o365"
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

var o365Options o365.Options

// O365Cmd represents the bruteSpray command
var O365Cmd = &cobra.Command{
	Use:   "o365",
	Short: "Commands for o365 module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
	PostRun: func(cmd *cobra.Command, args []string) {
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
	O365Cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	O365Cmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	O365Cmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	O365Cmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	O365Cmd.AddCommand(bruteCmd)
	O365Cmd.AddCommand(enumCmd)
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
		o365Options.ProxyHTTP = http.ProxyURL(url)
	}
}
