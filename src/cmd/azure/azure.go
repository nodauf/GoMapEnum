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
var validUsers string
var output string
var proxyString string

// AzureCmd represents the azure command
var AzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Commands for azure module",
	Long:  `Authentication on Azure could be used to enumerate valid user or not by trying to authenticate.`,
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
		azureOptions.ProxyHTTP = http.ProxyURL(url)
	}
}
