package linkedin

import (
	"GoMapEnum/src/logger"
	"errors"

	"github.com/spf13/cobra"
)

// gatherCmd represents the gather command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Search on Linkedin for people working in the specified company",
	Long: `Firstly, it will search for company based on the provided name and then list all the people working at these companies and print them in the specified format.
The session cookie is needed to use the Linkedin features.`,
	Example: `go run main.go linkedin gather -c contoso -f "{f}{last}@contonso.com" -e -s AQEDA...`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if linkedinOptions.Format == "" && linkedinOptions.Email {
			return errors.New("format flag is requre when email should be guess")
		}
		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Gather", "linkedin", "Linkedin")
		log.SetLevel(level)
		log.Info("Starting the module linkedin")
		linkedinOptions.Log = log
		validUsers = linkedinOptions.Gather()
	},
}

func init() {

	gatherCmd.Flags().StringVarP(&linkedinOptions.Format, "format", "f", "", "Format (ex:{first}.{last}@domain.com, domain\\{f}{last}")
	gatherCmd.Flags().StringVarP(&linkedinOptions.Company, "company", "c", "", "Company name")
	gatherCmd.Flags().BoolVar(&linkedinOptions.Email, "email", true, "Guess the email according to the format. If false print the first name and last name")
	gatherCmd.Flags().BoolVarP(&linkedinOptions.ExactMatch, "exactMatch", "e", false, "Exact match of the company's name")
	gatherCmd.Flags().StringVarP(&linkedinOptions.Cookie, "cookie", "s", "", "Session cookie named li_at")
	gatherCmd.MarkFlagRequired("company")
	gatherCmd.MarkFlagRequired("cookie")
}
