#!/bin/sh
# NeetoAuth CLI — session-start hook for Claude Code
# Lightweight auth liveness check. Always exits 0 (informational).

if ! command -v neetoauth >/dev/null 2>&1; then
  echo "NeetoAuth CLI is not installed or not on PATH."
  exit 0
fi

if neetoauth whoami >/dev/null 2>&1; then
  echo "NeetoAuth plugin active."
else
  echo "NeetoAuth CLI installed but not authenticated. Run 'neetoauth login' to authenticate."
fi

exit 0
