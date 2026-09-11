# NeetoAuth CLI

NeetoAuth CLI lets you manage team members and inspect product roles in your NeetoAuth workspace from the terminal, and integrates with Claude Code as a skill.

<!-- neeto-cli-commons:installation:start -->
## Installation

### macOS / Linux

**Homebrew (recommended on macOS):**

```bash
brew install neetozone/tap/neetoauth
```

**Shell script:**

```bash
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoAuth/latest/install.sh | sh
```

This verifies the download's SHA-256 checksum against the published `SHA256SUMS`,
then installs to `/usr/local/bin` (may prompt for sudo). Set `NEETOAUTH_INSTALL_DIR`
to a directory you own to install without sudo.

### Windows

**PowerShell:**

```powershell
irm https://neeto-downloads.s3.amazonaws.com/cli/NeetoAuth/latest/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://neeto-downloads.s3.amazonaws.com/cli/NeetoAuth/latest/install.cmd -o install.cmd && install.cmd
```

Both verify the download's SHA-256 checksum before installing to
`%LOCALAPPDATA%\Programs\neetoauth` and adding it to your user PATH. Set
`NEETOAUTH_INSTALL_DIR` to install somewhere else.
<!-- neeto-cli-commons:installation:end -->

<!-- neeto-cli-commons:verify-installation:start -->
### Verify installation

```bash
neetoauth --help
```
<!-- neeto-cli-commons:verify-installation:end -->

<!-- neeto-cli-commons:prerequisites:start -->
## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.26.1+
- Access to a NeetoAuth organization
<!-- neeto-cli-commons:prerequisites:end -->

## Development

```bash
git clone https://github.com/neetozone/neeto-auth-cli.git
cd neeto-auth-cli
bin/setup
```

This installs Go dependencies, golangci-lint, configures git hooks, and builds the binary.

<!-- neeto-cli-commons:make-targets:start -->
### Make targets

```bash
make build          # Builds ./neetoauth
make test           # Run tests
make lint           # golangci-lint
make fmt            # gofmt -w
make vet            # go vet
make check          # fmt + vet + test
make install        # Installs to /usr/local/bin
make clean          # Remove built binary
```
<!-- neeto-cli-commons:make-targets:end -->

### Pointing to a local or staging server

Set `NEETOAUTH_BASE_URL` to override the default `https://<subdomain>.neetoauth.com`:

```bash
export NEETOAUTH_BASE_URL=http://acme.lvh.me:8980
neetoauth login --subdomain acme
```

<!-- neeto-cli-commons:global-flags:start -->
## Global flags

Every command accepts:

| Flag | Description |
|---|---|
| `--subdomain <name>` | Which logged-in subdomain to use (required when multiple are logged in). |
| `--json` | Force JSON envelope output. |
| `--quiet` | Emit raw data only. Action commands print just the identifier; `delete` prints `success`. |
| `--toon` | TOON (Token-Optimized Output Notation) — compact format for LLMs. |
| `--verbose` | Expand every field of a record instead of a table. |
<!-- neeto-cli-commons:global-flags:end -->

## Adding product-specific commands

See [`docs/adding-commands.md`](docs/adding-commands.md) for the step-by-step
workflow for adding new resource commands that use the built-in auth, HTTP
client, and output helpers.

Quick API wrapper reference: [`docs/api-wrapper-reference.md`](docs/api-wrapper-reference.md).

<!-- neeto-cli-commons:release:start -->
## Release

Releases are cut by BigBinary's CI pipeline defined in
`.neetoci/release.yml`. Merging a PR with a `major` / `minor` / `patch`
label to `main` triggers the shared release script published by
`neeto-cli-commons`, which bumps and tags VERSION, runs GoReleaser,
uploads artifacts to `s3://neeto-downloads/cli/NeetoAuth/`, updates the
Homebrew tap (`neetozone/tap`), and pushes the version bump commit
straight to `main`.
<!-- neeto-cli-commons:release:end -->

<!-- neeto-cli-commons:ai-coding-assistants:start -->
## AI coding assistants

```bash
neetoauth setup claude      # Register plugin with Claude Code
neetoauth setup cursor      # Write .cursor/rules/neetoauth.mdc
neetoauth setup windsurf    # Write .windsurf/rules/neetoauth.md
neetoauth setup copilot     # Add a NeetoAuth section to .github/copilot-instructions.md
neetoauth setup gemini      # Add a NeetoAuth section to GEMINI.md
neetoauth setup codex       # Add a NeetoAuth section to AGENTS.md
```

Every command except `setup claude` writes into the current project directory, so
run these commands from the root of the project the assistant works in. Re-run
them after every upgrade: `setup cursor` and `setup windsurf` overwrite their rule
file, while `setup copilot`, `setup gemini` and `setup codex` keep the existing
content of their file and replace only the NeetoAuth section instead of adding a
duplicate.
<!-- neeto-cli-commons:ai-coding-assistants:end -->
