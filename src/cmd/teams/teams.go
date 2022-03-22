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
var validUsers string
var output string
var proxyString string

var teamsOptions teams.Options

// TeamsCmd represents the teams command
var TeamsCmd = &cobra.Command{
	Use:   "teams",
	Short: "Commands for teams module",
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
		teamsOptions.ProxyHTTP = http.ProxyURL(url)
	}
}
