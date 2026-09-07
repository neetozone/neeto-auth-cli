package commands

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

func newUsersCreateTestCommand(t *testing.T) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "create"}
	cmd.Flags().String("email", "", "")
	cmd.Flags().String("role", "", "")
	cmd.Flags().String("json-file", "", "")
	return cmd
}

func writeTempJSONFile(t *testing.T, content string) string {
	t.Helper()
	path := t.TempDir() + "/user.json"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestSatisfyRequiredFlagsFromNestedUserJSONFile(t *testing.T) {
	t.Run("no json-file is a no-op", func(t *testing.T) {
		cmd := newUsersCreateTestCommand(t)
		if err := satisfyRequiredFlagsFromNestedUserJSONFile(cmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd.Flags().Changed("email") {
			t.Error("email should not be marked changed")
		}
	})

	t.Run("nested user object satisfies email and role", func(t *testing.T) {
		path := writeTempJSONFile(t, `{"user": {"email": "sam.smith@example.com", "role": "owner"}}`)
		cmd := newUsersCreateTestCommand(t)
		if err := cmd.Flags().Set("json-file", path); err != nil {
			t.Fatalf("failed to set json-file flag: %v", err)
		}
		if err := satisfyRequiredFlagsFromNestedUserJSONFile(cmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got, _ := cmd.Flags().GetString("email"); got != "sam.smith@example.com" {
			t.Errorf("email = %q, want sam.smith@example.com", got)
		}
		if got, _ := cmd.Flags().GetString("role"); got != "owner" {
			t.Errorf("role = %q, want owner", got)
		}
		if !cmd.Flags().Changed("email") || !cmd.Flags().Changed("role") {
			t.Error("expected email and role to be marked changed so cobra's required-flag check passes")
		}
	})

	t.Run("explicit flag is not overridden by the file", func(t *testing.T) {
		path := writeTempJSONFile(t, `{"user": {"email": "sam.smith@example.com", "role": "owner"}}`)
		cmd := newUsersCreateTestCommand(t)
		if err := cmd.Flags().Set("json-file", path); err != nil {
			t.Fatalf("failed to set json-file flag: %v", err)
		}
		if err := cmd.Flags().Set("email", "oliver.smith@example.com"); err != nil {
			t.Fatalf("failed to set email flag: %v", err)
		}
		if err := satisfyRequiredFlagsFromNestedUserJSONFile(cmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got, _ := cmd.Flags().GetString("email"); got != "oliver.smith@example.com" {
			t.Errorf("email = %q, want oliver.smith@example.com (explicit flag must win)", got)
		}
	})

	t.Run("flat json file is left to the shared json-file shim", func(t *testing.T) {
		path := writeTempJSONFile(t, `{"email": "sam.smith@example.com", "role": "owner"}`)
		cmd := newUsersCreateTestCommand(t)
		if err := cmd.Flags().Set("json-file", path); err != nil {
			t.Fatalf("failed to set json-file flag: %v", err)
		}
		if err := satisfyRequiredFlagsFromNestedUserJSONFile(cmd, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cmd.Flags().Changed("email") {
			t.Error("flat payloads should be handled by allowJSONFileToSatisfyRequiredFlags, not this hook")
		}
	})
}
