# NeetoAuth CLI

NeetoAuth CLI lets you manage team members and inspect product roles in your NeetoAuth workspace from the terminal, and integrates with Claude Code as a skill.

## Installation

### macOS / Linux

**Homebrew (recommended on macOS):**

```bash
brew trust neetozone/tap
brew install neetozone/homebrew-tap/neetoauth
```

**Shell script:**

```bash
curl -fsSL https://neetoauth.com/cli/install.sh | sh
```

### Windows

**PowerShell:**

```powershell
irm https://neetoauth.com/cli/install.ps1 | iex
```

**Command Prompt (CMD):**

```cmd
curl -fsSL https://neetoauth.com/cli/install.cmd -o install.cmd && install.cmd
```

### Verify installation

```bash
neetoauth --help
```

## Prerequisites (development)

- [Go](https://go.dev/dl/) 1.26.1+
- Access to a NeetoAuth organization

## Development

```bash
git clone https://github.com/neetozone/neeto-auth-cli.git
cd neeto-auth-cli
bin/setup
```

This installs Go dependencies, golangci-lint, configures git hooks, and builds the binary.

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

### Pointing to a local or staging server

Set `NEETOAUTH_BASE_URL` to override the default `https://<subdomain>.neetoauth.com`:

```bash
export NEETOAUTH_BASE_URL=http://acme.lvh.me:8980
neetoauth login --subdomain acme
```

## Global flags

Every command accepts:

| Flag | Description |
|---|---|
| `--subdomain <name>` | Which logged-in subdomain to use (required when multiple are logged in). |
| `--json` | Force JSON envelope output. |
| `--quiet` | Emit raw data only. Action commands print just the identifier; `delete` prints `success`. |
| `--toon` | TOON (Token-Optimized Output Notation) — compact format for LLMs. |

## Adding product-specific commands

See [`docs/adding-commands.md`](docs/adding-commands.md) for the step-by-step
workflow for adding new resource commands that use the built-in auth, HTTP
client, and output helpers.

Quick API wrapper reference: [`docs/api-wrapper-reference.md`](docs/api-wrapper-reference.md).

## Release

Releases are cut by BigBinary's CI pipeline defined in
`.neetoci/release.yml`. Merging a PR with a `major` / `minor` / `patch`
label to `main` triggers `.scripts/release.sh`, which tags the current
VERSION, runs GoReleaser, uploads artifacts to
`s3://neeto-downloads/cli/NeetoAuth/`, updates the Homebrew tap
(`neetozone/homebrew-tap`), and opens the next-version bump PR.

## AI coding assistants

```bash
neetoauth setup claude      # Register plugin with Claude Code
neetoauth setup cursor      # Write .cursor/rules/neetoauth.mdc
neetoauth setup windsurf    # Write .windsurf/rules/neetoauth.md
neetoauth setup copilot     # Append to .github/copilot-instructions.md
neetoauth setup gemini      # Append to GEMINI.md
neetoauth setup codex       # Append to AGENTS.md
```
