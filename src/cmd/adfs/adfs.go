package adfs

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
var proxy func(*http.Request) (*url.URL, error)
var validUsers []string
var output string
var proxyString string

// AdfsCmd represents the bruteSpray command
var AdfsCmd = &cobra.Command{
	Use:   "adfs",
	Short: "Commands for adfs module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	AdfsCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	AdfsCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	AdfsCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	AdfsCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	AdfsCmd.AddCommand(bruteCmd)
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
