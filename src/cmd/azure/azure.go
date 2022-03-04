package azure

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
var validUsers []string
var output string
var proxyString string

// AzureCmd represents the azure command
var AzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Commands for azure module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	AzureCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	AzureCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	AzureCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	AzureCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	AzureCmd.AddCommand(enumCmd)
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
