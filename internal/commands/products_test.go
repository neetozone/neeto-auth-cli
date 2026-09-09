package commands

import (
	"net/http"
	"strings"
	"testing"

	"github.com/neetozone/neeto-cli-commons/client"
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
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

func TestEnableDisableAcceptAtMostOneArg(t *testing.T) {
	for _, cmd := range []struct {
		name string
		args func([]string) error
	}{
		{"enable", func(a []string) error { return enableProductCmd.Args(enableProductCmd, a) }},
		{"disable", func(a []string) error { return disableProductCmd.Args(disableProductCmd, a) }},
	} {
		t.Run(cmd.name, func(t *testing.T) {
			if err := cmd.args(nil); err != nil {
				t.Errorf("expected no error with zero args, the product comes from --product, got: %v", err)
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
	if !strings.Contains(enableProductCmd.Example, "products enable --product crm") {
		t.Errorf("enable example = %q, want it to mention %q", enableProductCmd.Example, "products enable --product crm")
	}
	if !strings.Contains(disableProductCmd.Example, "products disable --product crm") {
		t.Errorf("disable example = %q, want it to mention %q", disableProductCmd.Example, "products disable --product crm")
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
		{Label: "Enable a product", Command: "neetoauth products enable --product <product>"},
		{Label: "Disable a product", Command: "neetoauth products disable --product <product>"},
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

func TestProductsEnableSendsPatchWithEnabledTrue(t *testing.T) {
	recorded, err := runCommand(t, http.StatusOK, `{}`, "products", "enable", "cal")
	if err != nil {
		t.Fatalf("products enable: %v", err)
	}
	if recorded.Method != http.MethodPatch {
		t.Errorf("method = %q, want %q", recorded.Method, http.MethodPatch)
	}
	if want := "/api/external/v2/products/cal"; recorded.Path != want {
		t.Errorf("path = %q, want %q", recorded.Path, want)
	}
	if recorded.Body["enabled"] != true {
		t.Errorf("body = %v, want enabled true", recorded.Body)
	}
}

func TestProductsDisableSendsPatchWithEnabledFalse(t *testing.T) {
	recorded, err := runCommand(t, http.StatusOK, `{}`, "products", "disable", "cal")
	if err != nil {
		t.Fatalf("products disable: %v", err)
	}
	if recorded.Body["enabled"] != false {
		t.Errorf("body = %v, want enabled false", recorded.Body)
	}
}

func TestProductsEnableEscapesTheProductNameInThePath(t *testing.T) {
	recorded, _ := runCommand(t, http.StatusNotFound, `{"error":"not found"}`, "products", "enable", "cal?x")

	if want := "/api/external/v2/products/cal%3Fx"; recorded.Path != want {
		t.Errorf("path = %q, want %q — an unescaped name silently retargets the request", recorded.Path, want)
	}
}

func TestProductsListReadsTheEnabledFlag(t *testing.T) {
	recorded, err := runCommand(t, http.StatusOK,
		`{"products":[{"name":"Cal","enabled":true,"roles":["Admin"]},{"name":"Git","enabled":false,"roles":[]}],"pagination":{"total_records":2,"total_pages":1,"current_page_number":1,"page_size":30}}`,
		"products", "list")
	if err != nil {
		t.Fatalf("products list: %v", err)
	}
	if want := "/api/external/v2/products"; recorded.Path != want {
		t.Errorf("path = %q, want %q", recorded.Path, want)
	}
}

func newProductFlagCommand(t *testing.T, value string) *cobra.Command {
	t.Helper()
	cmd := &cobra.Command{Use: "enable"}
	cmd.Flags().String("product", "", "Name of the product to enable")
	if value != "" {
		if err := cmd.Flags().Set("product", value); err != nil {
			t.Fatalf("set --product: %v", err)
		}
	}
	return cmd
}

func TestSelectedProductPrefersTheFlag(t *testing.T) {
	got, err := selectedProduct(newProductFlagCommand(t, "cal"), nil)
	if err != nil {
		t.Fatalf("selectedProduct: %v", err)
	}
	if got != "cal" {
		t.Errorf("selectedProduct = %q, want %q", got, "cal")
	}
}

func TestSelectedProductStillAcceptsThePositionalArgument(t *testing.T) {
	got, err := selectedProduct(newProductFlagCommand(t, ""), []string{"cal"})
	if err != nil {
		t.Fatalf("selectedProduct: %v", err)
	}
	if got != "cal" {
		t.Errorf("selectedProduct = %q, want %q", got, "cal")
	}
}

func TestSelectedProductRejectsBothForms(t *testing.T) {
	if _, err := selectedProduct(newProductFlagCommand(t, "cal"), []string{"desk"}); err == nil {
		t.Error("expected an error when the product is passed twice")
	}
}

func TestSelectedProductRequiresTheProduct(t *testing.T) {
	_, err := selectedProduct(newProductFlagCommand(t, ""), nil)
	if err == nil {
		t.Fatal("expected an error when no product is given")
	}
	if !strings.Contains(err.Error(), "--product") {
		t.Errorf("error = %q, want it to name --product", err.Error())
	}
}

func TestProductsEnableWithTheFlagSendsPatch(t *testing.T) {
	recorded, err := runCommand(t, http.StatusOK, `{}`, "products", "enable", "--product", "cal")
	if err != nil {
		t.Fatalf("products enable --product: %v", err)
	}
	if want := "/api/external/v2/products/cal"; recorded.Path != want {
		t.Errorf("path = %q, want %q", recorded.Path, want)
	}
	if recorded.Body["enabled"] != true {
		t.Errorf("body = %v, want enabled true", recorded.Body)
	}
}
