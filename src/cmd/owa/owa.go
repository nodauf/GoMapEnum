package owa

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/owa"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var proxy func(*http.Request) (*url.URL, error)
var validUsers []string
var output string
var proxyString string

var owaOptions owa.Options

// OwaCmd represents the owa command
var OwaCmd = &cobra.Command{
	Use:   "owa",
	Short: "Commands for owa module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	OwaCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	OwaCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	OwaCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	OwaCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	OwaCmd.AddCommand(enumCmd)
	OwaCmd.AddCommand(bruteCmd)
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
