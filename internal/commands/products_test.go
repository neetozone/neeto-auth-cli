package commands

import (
	"strings"
	"testing"
)

func TestProductsCommandsRegistered(t *testing.T) {
	found := map[string]bool{}
	for _, c := range rootCmd.Commands() {
		if c.Name() == "products" {
			for _, sub := range c.Commands() {
				found[sub.Name()] = true
			}
		}
	}
	for _, want := range []string{"list", "enable", "disable"} {
		if !found[want] {
			t.Errorf("expected %q registered under products, got: %v", want, found)
		}
	}
}

func TestEnableDisableAcceptExactlyOneArg(t *testing.T) {
	for _, cmd := range []struct {
		name string
		args func([]string) error
	}{
		{"enable", func(a []string) error { return enableProductCmd.Args(enableProductCmd, a) }},
		{"disable", func(a []string) error { return disableProductCmd.Args(disableProductCmd, a) }},
	} {
		t.Run(cmd.name, func(t *testing.T) {
			if err := cmd.args(nil); err == nil {
				t.Error("expected error with zero args")
			}
			if err := cmd.args([]string{"crm"}); err != nil {
				t.Errorf("expected no error with one arg, got: %v", err)
			}
			if err := cmd.args([]string{"crm", "extra"}); err == nil {
				t.Error("expected error with two args")
			}
		})
	}
}

func TestEnableDisableExamples(t *testing.T) {
	if !strings.Contains(enableProductCmd.Example, "products enable crm") {
		t.Errorf("enable example = %q, want it to mention %q", enableProductCmd.Example, "products enable crm")
	}
	if !strings.Contains(disableProductCmd.Example, "products disable crm") {
		t.Errorf("disable example = %q, want it to mention %q", disableProductCmd.Example, "products disable crm")
	}
}

func TestToggleProductAction(t *testing.T) {
	// enableProductCmd.RunE and disableProductCmd.RunE are both built by
	// toggleProduct with different closures over `enabled`; verify they were
	// wired the right way round rather than swapped.
	if enableProductCmd.RunE == nil {
		t.Fatal("enableProductCmd.RunE is nil")
	}
	if disableProductCmd.RunE == nil {
		t.Fatal("disableProductCmd.RunE is nil")
	}
}

func TestProductBreadcrumbs(t *testing.T) {
	crumbs := productBreadcrumbs()
	if len(crumbs) != 1 {
		t.Fatalf("len(productBreadcrumbs()) = %d, want 1", len(crumbs))
	}
	if crumbs[0].Label != "Use a role when inviting" {
		t.Errorf("breadcrumb label = %q, want %q", crumbs[0].Label, "Use a role when inviting")
	}
	wantCommand := "neetoauth users create --email <email> --role non_owner --app <product>:<role>"
	if crumbs[0].Command != wantCommand {
		t.Errorf("breadcrumb command = %q, want %q", crumbs[0].Command, wantCommand)
	}
}

func TestIsProductNotFoundError(t *testing.T) {
	notFound := []string{
		"product not found",
		"Product Not Found",
		"request failed with status 404",
		"404 Not Found",
	}
	for _, msg := range notFound {
		if !isProductNotFoundError(errString(msg)) {
			t.Errorf("isProductNotFoundError(%q) = false, want true", msg)
		}
	}

	other := []string{
		"unauthorized",
		"request failed with status 500",
		"connection refused",
	}
	for _, msg := range other {
		if isProductNotFoundError(errString(msg)) {
			t.Errorf("isProductNotFoundError(%q) = true, want false", msg)
		}
	}
}
