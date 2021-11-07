package cmd

import (
	"GoMapEnum/src/cmd/brute"
	"GoMapEnum/src/cmd/enum"
	"GoMapEnum/src/cmd/gather"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "GoMapEnum",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Add child
	rootCmd.AddCommand(enum.UserenumCmd)
	rootCmd.AddCommand(brute.BruteSprayCmd)
	rootCmd.AddCommand(gather.GatherCmd)
}
