package enum

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"GoMapEnum/src/logger"
)

var level logger.Level
var verbose bool
var debug bool
var output string
var validUsers []string
var proxy func(*http.Request) (*url.URL, error)

var proxyString string

// UserenumCmd represents the userenum command
var UserenumCmd = &cobra.Command{
	Use:   "userenum",
	Short: "User enumeration",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if output != "" {
			if err := os.WriteFile(output, []byte(strings.Join(validUsers, "\n")), 0666); err != nil {
				fmt.Println(err)
			}
		}
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	cobra.OnInitialize(initProxy)
	UserenumCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	UserenumCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	UserenumCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	UserenumCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Sleep in seconds before sending an authentication request")

	// Add child
	UserenumCmd.AddCommand(azureCmd)
	UserenumCmd.AddCommand(o365Cmd)
	UserenumCmd.AddCommand(owaCmd)
	UserenumCmd.AddCommand(teamsCmd)

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
