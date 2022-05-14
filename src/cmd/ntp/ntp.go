package ntp

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/ntp"
	"strings"

	"github.com/spf13/cobra"
)

var level logger.Level
var verbose bool
var debug bool
var proxyString string

var ntpOptions ntp.Options

// NTPCmd represents the bruteSpray command
var NTPCmd = &cobra.Command{
	Use:   "ntp",
	Short: "Commands for ntp module",
	Long: `For some modules to works like Kerberos, the systems needs to have the clock sync.
	In some case, it is not possible to directly synchronize and, for example, you have to use a third machine to get the time of the remote server.`,
	Example: strings.Join([]string{`go run main.go ntp -t ntp.example.com`,
		`go run main.go ntp -t ntp.example.com --UTC`}, "\n"),
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New("GetTime", "ntp", ntpOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module ntp")
		ntpOptions.Log = log
		ntpOptions.GetTime()
	},
}

func init() {

	cobra.OnInitialize(initLogger)
	//cobra.OnInitialize(initProxy)
	NTPCmd.PersistentFlags().StringVarP(&ntpOptions.Target, "taget", "t", "", "NTP server to query")
	NTPCmd.PersistentFlags().BoolVar(&ntpOptions.UTC, "utc", false, "Print in UTC format")
	NTPCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose")
	NTPCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug")
	//NTPCmd.PersistentFlags().StringVar(&proxyString, "proxy", "", "Proxy to use (ex: http://localhost:8080)")

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

/*func initProxy() {
	if proxyString != "" {
		url, err := url.Parse(proxyString)
		if err != nil {
			fmt.Println("Fail to parse URL " + proxyString + " - error " + err.Error())
			os.Exit(1)
		}
		ntpOptions.ProxyHTTP = http.ProxyURL(url)
	}
}*/
