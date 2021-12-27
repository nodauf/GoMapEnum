package gather

import (
	"GoMapEnum/src/logger"
	searchengine "GoMapEnum/src/searchEngine"

	"github.com/spf13/cobra"
)

var searchEngineOptions searchengine.Options

// o365Cmd represents the o365 command
var searchEngineCmd = &cobra.Command{
	Use:     "searchEngine",
	Short:   "Search the company name on Bing and Google and parse the Linkedin results",
	Long:    `Credits: https://github.com/m8r0wn/CrossLinked`,
	Example: `go run main.go gather searchEngine -c contoso -f "{f}{last}@contonso.com" -v`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "searchEngine", "Google and Bing search engine")
		log.SetLevel(level)
		log.Info("Starting the module searchEngine")
		searchEngineOptions.Log = log
		searchEngineOptions.Proxy = proxy
		users = searchEngineOptions.Gather()
	},
}

func init() {

	searchEngineCmd.Flags().StringVarP(&searchEngineOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	searchEngineCmd.Flags().StringVarP(&searchEngineOptions.Company, "company", "c", "", "Company name")
	searchEngineCmd.Flags().BoolVarP(&searchEngineOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	searchEngineCmd.Flags().StringVarP(&searchEngineOptions.SearchEngine, "searchEngine", "s", "", "Select on which search engine the query will be made, ex: bing,google (default: all)")
	searchEngineCmd.MarkFlagRequired("company")
	searchEngineCmd.MarkFlagRequired("format")
}
