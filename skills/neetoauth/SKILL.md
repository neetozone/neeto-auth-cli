---
name: neetoauth
description: >
  Manage NeetoAuth from the command line.
  Use when the user asks about operations exposed by the NeetoAuth CLI.
---

## Prerequisites

Run `neetoauth doctor` to check authentication and connectivity.
If not authenticated, run `neetoauth login`.

## Authentication & multi-subdomain

Credentials for every logged-in subdomain are stored together in
`~/.config/neetoauth/auth.json`. A command that talks to the API picks which
subdomain to use by these rules:

- 0 subdomains logged in → every credential-using command errors with
  "not logged in. Run 'neetoauth login' to authenticate".
- 1 subdomain logged in → that one is the implicit default; `--subdomain`
  may be omitted.
- 2+ subdomains logged in → **`--subdomain <name>` is required** on every
  credential-using command, including `doctor`. The error lists every
  logged-in subdomain so the agent can offer a choice.

`login` / `logout` / `whoami` have dedicated behavior:

| Command | Behavior |
|---|---|
| `neetoauth login --subdomain <name>` | Adds or refreshes the entry for `<name>`. No flag → prompts for the subdomain. |
| `neetoauth logout --subdomain <name>` | Removes that one entry. |
| `neetoauth logout --all` | Removes every entry. |
| `neetoauth logout` (no flag) | Removes the only entry if exactly one is logged in; errors if multiple. |
| `neetoauth whoami` | Lists every logged-in account. Marks the entry `(default)` when exactly one. |
| `neetoauth whoami --subdomain <name>` | Shows just that one. |

## Global flags (persistent on every command)

| Flag | Purpose |
|---|---|
| `--subdomain <name>` | Select which logged-in subdomain the command targets. Required when multiple are logged in. |
| `--json` | Force JSON envelope output even on a TTY. |
| `--quiet` | Emit only the raw payload — no envelope, no breadcrumbs. For action commands (create/update), emits just the resource identifier; `delete` emits `success`. Designed for scripting. |
| `--toon` | Emit TOON (Token Optimized Output Notation). Preferred for feeding list/show output back to an LLM; ~30–60% fewer tokens than JSON. |

Precedence if multiple are set: `--toon` > `--quiet` > `--json` > pretty.

## Output modes & response envelope

**Pretty (default on a TTY)** — tables for arrays, key-value for objects,
breadcrumbs appended. Not intended for machine consumption.

**JSON envelope** (non-TTY, or `--json`):
```json
{
  "data": <resource body>,
  "breadcrumbs": [{ "label": "List", "command": "neetoauth <resource> list" }],
  "pagination": {
    "current_page_number": 1,
    "total_pages": 10,
    "total_records": 250
  }
}
```
`breadcrumbs` is omitted when empty. `pagination` is present only for list
commands.

**Quiet** (`--quiet`) — `data` contents only, no envelope. For action
commands `PrintQuiet` unwraps a single-key wrapper and prints the first of
`sid` / `id` / `name`. For `delete` it prints `success`.

**TOON** (`--toon`) — same data as JSON, re-encoded into TOON. Shape is
equivalent but whitespace/keys are compressed. Parse by re-reading keys as
you would JSON.

### Pagination

List commands accept `--page` (1-indexed) and `--page-size` (max 100).
The envelope's `pagination` field always exposes:
`current_page_number`, `total_pages`, `total_records`. Agents should loop
by incrementing `--page` until `current_page_number == total_pages`.

## Discovery

The full, always-accurate command tree (including any flags added after
this skill was built) is available as JSON:

```bash
neetoauth commands
```

Each catalog entry has `command`, `description`, optional `flags` (with
`name`, `type`, `default`, `description`, `required`), and `subcommands`.
Use this whenever a user asks about a flag or command not covered below.

## Diagnostics & IDE setup

| Command | Purpose |
|---|---|
| `doctor` | Auth check + API reachability + version. Uses `--subdomain` when multiple are logged in. |
| `version` | Print CLI version / commit / build date. |
| `commands` | Emit the full command/flag catalog as JSON. |
| `setup claude` | Install NeetoAuth plugin into Claude Code (`plugin.json`, hooks, this SKILL.md). |
| `setup cursor` / `windsurf` / `copilot` / `gemini` / `codex` | Write NeetoAuth rule files into the current project directory; re-run after an upgrade to refresh them. |

## Environment variable override

Set `NEETOAUTH_BASE_URL` to point the CLI at a staging or local server:

```bash
export NEETOAUTH_BASE_URL=http://acme.lvh.me:8980
neetoauth login --subdomain acme
```

## Error surface

Every command exits non-zero on failure and writes a single-line message to
stderr. Common errors the agent should expect:

- `not logged in. Run 'neetoauth login' to authenticate` — empty credential store.
- `multiple subdomains logged in (acme, beta); specify --subdomain` — pick one.
- `not logged in to "foo". Logged in subdomains: acme, beta` — bad `--subdomain`.
- `required flag(s) "xxx" not set` (from cobra) — missing required flag.
- API errors come through with the server's message body; inspect the
  JSON envelope (or the `--quiet` payload) for `error` / `errors` / `notice`
  keys and any suggestions the API returns.

## Resource reference

The CLI wraps NeetoAuth's external team-member API
(`/api/external/v2/...`). Two resources are exposed.

### Team members (`neetoauth users`)

| Command | Required flags / args | Optional flags | Returns |
|---|---|---|---|
| `users list` | — | `--page`, `--page-size` | `{ "data": [{ email, role, first_name, last_name }, …] }` (one entry per **active** workspace member). |
| `users create` | `--email *`, `--role *` | `--first-name`, `--last-name`, `--app <product>:<role>` (repeatable), `--json-file <path>` | `201 Created` with the newly-invited user as `{ "data": { email, role, first_name, last_name, apps: [...] } }`. |
| `users delete <email>` | `<email>` (positional) | — | `204 No Content`; the CLI prints `Removed <email> from the workspace.` Quiet mode emits `success`. |

Notes for `users create`:
- `--role` accepts the workspace-wide role (`owner`, `non_owner`).
- Each `--app` flag value is `name:role`. Example:
  `--app neetocal:admin --app neetoform:member`.
- `--json-file` accepts either `{ "user": { ... } }` or a flat user
  object; CLI flags override file values.
- The server rejects creating a duplicate active user — capture the
  error envelope when it does.

Notes for `users delete`:
- Soft-deletes the user (`deactivated_at`) and revokes all server-side
  sessions.
- Cannot remove the last owner of an organization — the API returns
  an error envelope explaining why.

### Products (`neetoauth products`)

| Command | Args | Returns |
|---|---|---|
| `products list` | — | `{ "data": [{ "name": "<app_name>", "roles": ["<role1>", "<role2>", …] }, …] }`. Use the `name`/`roles` pairs to construct valid `--app` values for `users create`. |

### Typical agent flow

1. `neetoauth doctor` — confirm auth is healthy.
2. `neetoauth products list --toon` — discover what apps and roles exist
   in this workspace.
3. `neetoauth users list --toon` — see who is already a member.
4. `neetoauth users create --email new@x.com --role non_owner --app neetocal:admin --quiet`
   — invite, then read the printed identifier.
5. `neetoauth users delete new@x.com --quiet` — remove if needed.
