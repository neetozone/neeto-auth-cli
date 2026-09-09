package commands

import (
	"net/http"
	"strings"
	"testing"

	"github.com/neetozone/neeto-cli-commons/client"
	"github.com/neetozone/neeto-cli-commons/output"
)

func TestProductsCommandsRegistered(t *testing.T) {
	found := map[string]bool{}
	for _, sub := range productsCmd.Commands() {
		found[sub.Name()] = true
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
	if enableProductCmd.RunE == nil {
		t.Fatal("enableProductCmd.RunE is nil")
	}
	if disableProductCmd.RunE == nil {
		t.Fatal("disableProductCmd.RunE is nil")
	}
}

func TestProductBreadcrumbs(t *testing.T) {
	want := []output.Breadcrumb{
		{Label: "Enable a product", Command: "neetoauth products enable <product>"},
		{Label: "Disable a product", Command: "neetoauth products disable <product>"},
		{Label: "Use a role when inviting", Command: "neetoauth users create --email <email> --role non_owner --app <product>:<role>"},
	}

	crumbs := productBreadcrumbs()
	if len(crumbs) != len(want) {
		t.Fatalf("len(productBreadcrumbs()) = %d, want %d", len(crumbs), len(want))
	}
	for i, crumb := range want {
		if crumbs[i] != crumb {
			t.Errorf("breadcrumb %d = %+v, want %+v", i, crumbs[i], crumb)
		}
	}
}

func TestProductsListMentionsEnabledState(t *testing.T) {
	if !strings.Contains(productsListCmd.Short, "enabled state") {
		t.Errorf("products list Short = %q, want it to mention the enabled state", productsListCmd.Short)
	}
}

func TestIsProductNotFoundError(t *testing.T) {
	if !isProductNotFoundError(&client.APIError{StatusCode: http.StatusNotFound}) {
		t.Error("isProductNotFoundError(404) = false, want true")
	}
	if isProductNotFoundError(&client.APIError{StatusCode: http.StatusInternalServerError}) {
		t.Error("isProductNotFoundError(500) = true, want false")
	}
	if isProductNotFoundError(errString("plain error, not an APIError")) {
		t.Error("isProductNotFoundError(non-APIError) = true, want false")
	}
}

type errString string

func (e errString) Error() string { return string(e) }
