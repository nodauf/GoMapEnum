package google

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/google"
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

var googleOptions google.Options

// GoogleCmd represents the google command
var GoogleCmd = &cobra.Command{
	Use:   "google",
	Short: "Commands for google module",
	Long:  ``,
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
	GoogleCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	GoogleCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	GoogleCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")

	GoogleCmd.AddCommand(enumCmd)
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
		googleOptions.ProxyHTTP = http.ProxyURL(url)
	}
}
