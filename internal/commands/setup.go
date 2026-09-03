package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/neetozone/neeto-auth-cli/internal/plugin"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Set up NeetoAuth for AI coding assistants",
	Long: "Set up NeetoAuth for AI coding assistants.\n\n" +
		"Every subcommand except \"claude\" writes into the current project directory,\n" +
		"so run it from the root of the project you want the assistant to use NeetoAuth in.",
}

// --- Claude Code ---

var setupClaudeCmd = &cobra.Command{
	Use:   "claude",
	Short: "Register NeetoAuth plugin with Claude Code",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}

		if _, err := os.Stat(filepath.Join(home, ".claude")); os.IsNotExist(err) {
			return fmt.Errorf("claude Code not found (~/.claude/ does not exist)")
		}

		dest := filepath.Join(home, ".config", "neetoauth", "claude-plugin")
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("could not clean destination: %w", err)
		}

		if err := plugin.ExtractClaudePlugin(dest); err != nil {
			return fmt.Errorf("could not extract plugin: %w", err)
		}

		// Make hooks executable
		hooksDir := filepath.Join(dest, "hooks")
		entries, _ := os.ReadDir(hooksDir)
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".sh") {
				_ = os.Chmod(filepath.Join(hooksDir, e.Name()), 0o755)
			}
		}

		fmt.Printf("NeetoAuth plugin extracted to %s\n", dest)
		fmt.Println()
		fmt.Println("To finish installation, open Claude Code and run these slash commands:")
		fmt.Printf("  /plugin marketplace add %s\n", dest)
		fmt.Printf("  /plugin install %s@%s\n", plugin.PluginName, plugin.MarketplaceName)
		fmt.Println()
		fmt.Println("(Claude Code installs plugins via interactive slash commands — there is no shell equivalent today.)")
		return nil
	},
}

// --- Cursor ---

var setupCursorCmd = &cobra.Command{
	Use:   "cursor",
	Short: "Write NeetoAuth rules for Cursor IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		target := filepath.Join(".cursor", "rules", "neetoauth.mdc")
		return writeCreateMode(target, cursorContent())
	},
}

func cursorContent() string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("description: \"Command-line interface for NeetoAuth team member management.\"\n")
	b.WriteString("alwaysApply: true\n")
	b.WriteString("---\n\n")
	b.WriteString(plugin.SkillBody())
	return b.String()
}

// --- Windsurf ---

var setupWindsurfCmd = &cobra.Command{
	Use:   "windsurf",
	Short: "Write NeetoAuth rules for Windsurf IDE",
	RunE: func(cmd *cobra.Command, args []string) error {
		target := filepath.Join(".windsurf", "rules", "neetoauth.md")
		return writeCreateMode(target, windsurfContent())
	},
}

func windsurfContent() string {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("trigger: always_on\n")
	b.WriteString("description: \"Command-line interface for NeetoAuth team member management.\"\n")
	b.WriteString("---\n\n")
	b.WriteString(plugin.SkillBody())
	return b.String()
}

// --- Copilot ---

var copilotInstructionsPath = filepath.Join(".github", "copilot-instructions.md")

var setupCopilotCmd = &cobra.Command{
	Use:   "copilot",
	Short: "Add NeetoAuth instructions for GitHub Copilot",
	Long:  sectionHelp(copilotInstructionsPath),
	RunE: func(cmd *cobra.Command, args []string) error {
		return writeSection(cmd.OutOrStdout(), copilotInstructionsPath, plugin.SkillBody())
	},
}

// --- Gemini ---

var setupGeminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "Add NeetoAuth instructions for Gemini CLI",
	Long:  sectionHelp("GEMINI.md"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return writeSection(cmd.OutOrStdout(), "GEMINI.md", plugin.SkillBody())
	},
}

// --- Codex ---

var setupCodexCmd = &cobra.Command{
	Use:   "codex",
	Short: "Add NeetoAuth instructions for OpenAI Codex",
	Long:  sectionHelp("AGENTS.md"),
	RunE: func(cmd *cobra.Command, args []string) error {
		return writeSection(cmd.OutOrStdout(), "AGENTS.md", plugin.SkillBody())
	},
}

func sectionHelp(target string) string {
	return fmt.Sprintf(
		"Add a NeetoAuth section to %s in the current project directory.\n\n"+
			"Existing content in the file is kept. Re-running replaces the NeetoAuth section\n"+
			"instead of adding a duplicate, so run it again after every upgrade.",
		target,
	)
}

// --- Helpers ---

// writeCreateMode writes a file, creating parent dirs. Skips if already present.
func writeCreateMode(target, content string) error {
	if _, err := os.Stat(target); err == nil {
		fmt.Printf("Already installed: %s\n", target)
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		return err
	}

	fmt.Printf("Wrote %s\n", target)
	return nil
}

func sectionMarkers() (string, string) {
	name := rootCmd.Name()
	return "<!-- " + name + ":start -->", "<!-- " + name + ":end -->"
}

func writeSection(w io.Writer, target, body string) error {
	_, statErr := os.Stat(target)
	start, end := sectionMarkers()
	if _, err := upsertBlock(target, start, end, "## NeetoAuth CLI\n\n"+body); err != nil {
		return err
	}
	if statErr == nil {
		fmt.Fprintf(w, "Updated %s\n", target)
	} else {
		fmt.Fprintf(w, "Wrote %s\n", target)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
	setupCmd.AddCommand(setupClaudeCmd)
	setupCmd.AddCommand(setupCursorCmd)
	setupCmd.AddCommand(setupWindsurfCmd)
	setupCmd.AddCommand(setupCopilotCmd)
	setupCmd.AddCommand(setupGeminiCmd)
	setupCmd.AddCommand(setupCodexCmd)
}
