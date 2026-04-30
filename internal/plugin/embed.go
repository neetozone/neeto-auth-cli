package plugin

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

//go:embed skill.md
var skillContent string

// SkillBody returns the SKILL.md content with YAML frontmatter stripped.
func SkillBody() string {
	content := skillContent
	if strings.HasPrefix(content, "---") {
		if idx := strings.Index(content[3:], "---"); idx != -1 {
			content = strings.TrimLeft(content[3+idx+3:], "\n")
		}
	}
	return content
}

// ExtractClaudePlugin writes a complete .claude-plugin/ directory to dest,
// including plugin.json, hooks, commands, and a copy of SKILL.md.
func ExtractClaudePlugin(dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	manifest := map[string]interface{}{
		"name":        "neetoauth",
		"description": "Command-line interface for neetoAuth team member management.",
		"author": map[string]string{
			"name":  "BigBinary",
			"email": "support@bigbinary.com",
		},
		"homepage":   "https://github.com/neetozone/neeto-auth-cli",
		"repository": "https://github.com/neetozone/neeto-auth-cli",
	}
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	if err := writeFile(filepath.Join(dest, "plugin.json"), manifestJSON, 0o644); err != nil {
		return err
	}

	hooks := map[string]interface{}{
		"hooks": map[string]interface{}{
			"SessionStart": []interface{}{
				map[string]interface{}{
					"hooks": []interface{}{
						map[string]interface{}{
							"type":    "command",
							"command": "${CLAUDE_PLUGIN_ROOT}/hooks/session-start.sh",
							"timeout": 5,
						},
					},
				},
			},
		},
	}
	hooksJSON, _ := json.MarshalIndent(hooks, "", "  ")
	if err := writeFile(filepath.Join(dest, "hooks", "hooks.json"), hooksJSON, 0o644); err != nil {
		return err
	}

	sessionStartSh := "#!/bin/sh\n" +
		"# NeetoAuth CLI — session-start hook for Claude Code\n" +
		"# Lightweight auth liveness check. Always exits 0 (informational).\n" +
		"\n" +
		"if ! command -v neetoauth >/dev/null 2>&1; then\n" +
		"  echo \"NeetoAuth CLI is not installed or not on PATH.\"\n" +
		"  exit 0\n" +
		"fi\n" +
		"\n" +
		"if neetoauth whoami >/dev/null 2>&1; then\n" +
		"  echo \"NeetoAuth plugin active.\"\n" +
		"else\n" +
		"  echo \"NeetoAuth CLI installed but not authenticated. Run 'neetoauth login' to authenticate.\"\n" +
		"fi\n" +
		"\n" +
		"exit 0\n"
	if err := writeFile(filepath.Join(dest, "hooks", "session-start.sh"), []byte(sessionStartSh), 0o755); err != nil {
		return err
	}

	doctorMd := "---\n" +
		"name: neetoauth-doctor\n" +
		"description: Check NeetoAuth CLI health — auth, API connectivity.\n" +
		"invocable: true\n" +
		"---\n" +
		"\n" +
		"Run `neetoauth doctor` and report the results to the user.\n"
	if err := writeFile(filepath.Join(dest, "commands", "doctor.md"), []byte(doctorMd), 0o644); err != nil {
		return err
	}

	if err := writeFile(filepath.Join(dest, "skills", "neetoauth", "SKILL.md"), []byte(skillContent), 0o644); err != nil {
		return err
	}

	return nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
