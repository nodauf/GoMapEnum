package brute

import (
	"GoMapEnum/src/logger"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var output string
var noBruteforce bool
var sleep int
var proxy func(*http.Request) (*url.URL, error)

var proxyString string
var validUsers []string

// BruteSprayCmd represents the bruteSpray command
var BruteSprayCmd = &cobra.Command{
	Use:   "bruteSpray",
	Short: "Spray a password or bruteforce a user's password",
	Long:  `Different services are supported. The authentication could be on an ADFS instance, an o365 or an OWA.`,
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
	BruteSprayCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	BruteSprayCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	BruteSprayCmd.PersistentFlags().StringVarP(&output, "output-file", "o", "", "The out file for valid emails")
	BruteSprayCmd.PersistentFlags().BoolVarP(&noBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	BruteSprayCmd.PersistentFlags().IntVarP(&sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	BruteSprayCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Sleep in seconds before sending an authentication request")

	BruteSprayCmd.AddCommand(o365Cmd)
	BruteSprayCmd.AddCommand(adfsCmd)
	BruteSprayCmd.AddCommand(owaCmd)
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
