package smb

import (
	"GoMapEnum/src/logger"
	"GoMapEnum/src/modules/smb"
	"GoMapEnum/src/orchestrator"
	"errors"

	"github.com/spf13/cobra"
)

// bruteCmd represents the smtp command
var bruteCmd = &cobra.Command{
	Use:   "brute",
	Short: "Enumerate email address by connection to the smtp port of the target.",
	Long: `SMTP user enumeration with RCPT, VRFY and EXPN.
	Credits: https://github.com/cytopia/smtp-user-enum`,
	Example: `go run main.go smb brute -u users  -t mail.contoso.com -d domain.tld -o validUsers`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if smbOptions.Passwords == "" && smbOptions.Hash == "" {
			return errors.New("The field password or hash is required")
		} else if smbOptions.Passwords != "" && smbOptions.Hash != "" {
			return errors.New("Only password or hash field should be specified")
		}

		return nil
	},
	Run: func(cmdCli *cobra.Command, args []string) {
		log := logger.New("Bruteforce", "SMB", smbOptions.Target)
		log.SetLevel(level)
		log.Info("Starting the module SMB")
		smbOptions.Log = log

		// If we bruteforce with hash, password will contains the hashs as we use the orchestrator and we want to keep the same format than CME
		if smbOptions.Hash != "" {
			smbOptions.Passwords = smbOptions.Hash
			smbOptions.IsHash = true
		}

		orchestratorOptions := orchestrator.Orchestrator{}
		orchestratorOptions.AuthenticationFunc = smb.Authenticate
		orchestratorOptions.PreActionBruteforce = smb.RetrieveTargetInfo
		validUsers = orchestratorOptions.Bruteforce(&smbOptions)

	},
}

func init() {
	bruteCmd.Flags().BoolVarP(&smbOptions.NoBruteforce, "no-bruteforce", "n", false, "No spray when using file for username and password (user1 => password1, user2 => password2)")
	bruteCmd.Flags().StringVarP(&smbOptions.Users, "user", "u", "", "User or file containing the emails")
	bruteCmd.Flags().StringVarP(&smbOptions.Passwords, "password", "p", "", "Password or file containing the passwords")
	bruteCmd.Flags().StringVarP(&smbOptions.Hash, "hash", "H", "", "Hash or file containing the hashes")
	bruteCmd.Flags().StringVarP(&smbOptions.Target, "target", "t", "", "Target host where the port 445 is open")
	bruteCmd.Flags().StringVarP(&smbOptions.Domain, "domain", "d", "", "Domain for the authentication (by default the domain name will be guessed)")
	bruteCmd.Flags().IntVar(&smbOptions.Thread, "thread", 2, "Number of threads")
	bruteCmd.Flags().IntVar(&smbOptions.Timeout, "timeout", 5, "Timeout for the SMB connection in seconds")
	bruteCmd.Flags().BoolVar(&smbOptions.StopOnLockout, "stopOnLockout", true, "Stop the bruteforce if an account is locked out")
	bruteCmd.Flags().IntVarP(&smbOptions.Sleep, "sleep", "s", 0, "Sleep in seconds before sending an authentication request")
	bruteCmd.MarkFlagRequired("user")
	bruteCmd.MarkFlagRequired("target")
}
