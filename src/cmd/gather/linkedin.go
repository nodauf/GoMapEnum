package gather

import (
	"GoMapEnum/src/linkedin"
	"GoMapEnum/src/logger"

	"github.com/spf13/cobra"
)

var linkedinOptions linkedin.Options

// o365Cmd represents the o365 command
var linkedinCmd = &cobra.Command{
	Use:   "linkedin",
	Short: "Search on Linkedin for people working in the specified company",
	Long: `Firstly, it will search for company based on the provided name and then list all the people working at these companies and print them in the specified format.
The session cookie is needed to use the Linkedin features.`,
	Example: `go run main.go gather linkedin -c contoso -f "{f}{last}@contonso.com" -e -s AQEDA...`,
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "linkedin", "Linkedin")
		log.SetLevel(level)
		log.Info("Starting the module linkedin")
		linkedinOptions.Log = log
		linkedinOptions.Proxy = proxy
		users = linkedinOptions.Gather()
	},
}

func init() {

	linkedinCmd.Flags().StringVarP(&linkedinOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	linkedinCmd.Flags().StringVarP(&linkedinOptions.Company, "company", "c", "", "Company name")
	linkedinCmd.Flags().BoolVarP(&linkedinOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	linkedinCmd.Flags().StringVarP(&linkedinOptions.Cookie, "cookie", "s", "", "Session cookie named li_at")
	linkedinCmd.MarkFlagRequired("company")
	linkedinCmd.MarkFlagRequired("cookie")
	linkedinCmd.MarkFlagRequired("format")
}
