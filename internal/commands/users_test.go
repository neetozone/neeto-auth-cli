package commands

import (
	"net/http"
	"testing"
)

func TestUsersDeleteEscapesTheEmailInThePath(t *testing.T) {
	recorded, _ := runCommand(t, http.StatusNoContent, ``, "users", "delete", "oliver.smith@example.com?x")

	if want := "/api/external/v2/users/oliver.smith@example.com%3Fx"; recorded.Path != want {
		t.Errorf("path = %q, want %q — an unescaped address silently retargets the request", recorded.Path, want)
	}
}

func TestUsersDeleteLeavesAPlainEmailUnchanged(t *testing.T) {
	recorded, err := runCommand(t, http.StatusNoContent, ``, "users", "delete", "oliver.smith@example.com")
	if err != nil {
		t.Fatalf("users delete: %v", err)
	}
	if recorded.Method != http.MethodDelete {
		t.Errorf("method = %q, want %q", recorded.Method, http.MethodDelete)
	}
	if want := "/api/external/v2/users/oliver.smith@example.com"; recorded.Path != want {
		t.Errorf("path = %q, want %q", recorded.Path, want)
	}
}
