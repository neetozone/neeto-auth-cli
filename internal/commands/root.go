package commands

import (
	"fmt"
	"os"

	"github.com/neetozone/neeto-auth-cli/internal/output"
	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:           "neetoauth",
	Short:         "Command-line interface for NeetoAuth team member management.",
	Long:          "NeetoAuth CLI lets you manage team members and inspect product roles in your NeetoAuth workspace from the terminal, and integrates with Claude Code as a skill.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		jsonFlag, _ := cmd.Flags().GetBool("json")
		quietFlag, _ := cmd.Flags().GetBool("quiet")
		toonFlag, _ := cmd.Flags().GetBool("toon")
		output.ForceJSON = jsonFlag
		output.QuietMode = quietFlag
		output.ToonMode = toonFlag
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("json", false, "Output as JSON")
	rootCmd.PersistentFlags().Bool("quiet", false, "Output raw data only (no envelope)")
	rootCmd.PersistentFlags().Bool("toon", false, "Output in TOON format (token-optimized for AI agents)")
	rootCmd.PersistentFlags().String("subdomain", "", "Override saved subdomain")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
