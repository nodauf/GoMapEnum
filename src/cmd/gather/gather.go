package gather

import (
	"GoMapEnum/src/logger"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var proxy func(*http.Request) (*url.URL, error)
var users []string
var output string
var proxyString string

// GatherCmd represents the gather command
var GatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Retrieve a list of email address based on the company name",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(strings.Join(users, "\n")), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	GatherCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	GatherCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	GatherCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	GatherCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Sleep in seconds before sending an authentication request")

	GatherCmd.AddCommand(searchEngineCmd)
	GatherCmd.AddCommand(linkedinCmd)
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
			log.Fatalln("Fail to parse URL " + proxyString + " - error " + err.Error())
		}
		proxy = http.ProxyURL(url)
	}
}
