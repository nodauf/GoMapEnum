package linkedin

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/linkedin"
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

var linkedinOptions linkedin.Options

// LinkedinCmd represents the linkedin command
var LinkedinCmd = &cobra.Command{
	Use:   "linkedin",
	Short: "Commands for linkedin module",
	Long:  `Use linkedin API with the session cookie to an actions.`,
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
	LinkedinCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	LinkedinCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	LinkedinCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	LinkedinCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

	LinkedinCmd.AddCommand(gatherCmd)
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
		linkedinOptions.ProxyHTTP = http.ProxyURL(url)
	}
}
