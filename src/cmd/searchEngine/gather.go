package searchEngine

import (
	"GoMapEnum/src/logger"

	"github.com/spf13/cobra"
)

// o365Cmd represents the o365 command
var gatherCmd = &cobra.Command{
	Use:     "gather",
	Short:   "Search the company name on Bing and Google and parse the Linkedin results",
	Long:    `Credits: https://github.com/m8r0wn/CrossLinked`,
	Example: `go run main.go gather searchEngine -c contoso -f "{f}{last}@contonso.com" -v`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "searchEngine", "Google and Bing search engine")
		log.SetLevel(level)
		log.Info("Starting the module searchEngine")
		searchEngineOptions.Log = log
		validUsers = searchEngineOptions.Gather()
	},
}

func init() {

	gatherCmd.Flags().StringVarP(&searchEngineOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	gatherCmd.Flags().StringVarP(&searchEngineOptions.Company, "company", "c", "", "Company name")
	gatherCmd.Flags().BoolVarP(&searchEngineOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	gatherCmd.Flags().StringVarP(&searchEngineOptions.SearchEngine, "searchEngine", "s", "", "Select on which search engine the query will be made, ex: bing,google (default: all)")
	gatherCmd.MarkFlagRequired("company")
	gatherCmd.MarkFlagRequired("format")
}
