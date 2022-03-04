package teams

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/teams"
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

var teamsOptions teams.Options

// TeamsCmd represents the owa command
var TeamsCmd = &cobra.Command{
	Use:   "smtp",
	Short: "Commands for owa module",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	TeamsCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	TeamsCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	TeamsCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")

	TeamsCmd.AddCommand(enumCmd)
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
