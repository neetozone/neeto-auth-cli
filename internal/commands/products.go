package commands

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/neetozone/neeto-cli-commons/client"
	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Inspect products, whether they are enabled, and their available roles",
}

var productsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List every product in the workspace with its enabled state and the roles it exposes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}
		data, err := c.Get("/products", nil)
		if err != nil {
			return fmt.Errorf("Failed to fetch products: %w", err)
		}
		printList(data, "products", productBreadcrumbs())
		return nil
	},
}

var enableProductCmd = &cobra.Command{
	Use:     "enable <product>",
	Short:   "Enable a product",
	Example: `  neetoauth products enable crm`,
	Args:    cobra.ExactArgs(1),
	RunE:    toggleProduct(true),
}

var disableProductCmd = &cobra.Command{
	Use:     "disable <product>",
	Short:   "Disable a product",
	Example: `  neetoauth products disable crm`,
	Args:    cobra.ExactArgs(1),
	RunE:    toggleProduct(false),
}

func toggleProduct(enabled bool) func(cmd *cobra.Command, args []string) error {
	action := "disable"
	verbed := "disabled"
	if enabled {
		action = "enable"
		verbed = "enabled"
	}
	return func(cmd *cobra.Command, args []string) error {
		product := args[0]
		c, err := getClient(cmd)
		if err != nil {
			return err
		}
		if _, err := c.Patch(fmt.Sprintf("/products/%s", url.PathEscape(product)), map[string]any{
			"enabled": enabled,
		}); err != nil {
			if isProductNotFoundError(err) {
				return fmt.Errorf("Product %q not found. Run \"neetoauth products list\" to see available products.", product)
			}
			return fmt.Errorf("Failed to %s product %q: %w", action, product, err)
		}
		printMessage(fmt.Sprintf("Product %q %s.", product, verbed))
		return nil
	}
}

func isProductNotFoundError(err error) bool {
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func productBreadcrumbs() []output.Breadcrumb {
	return []output.Breadcrumb{
		{Label: "Enable a product", Command: "neetoauth products enable <product>"},
		{Label: "Disable a product", Command: "neetoauth products disable <product>"},
		{Label: "Use a role when inviting", Command: "neetoauth users create --email <email> --role non_owner --app <product>:<role>"},
	}
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(productsCmd) })
	productsCmd.AddCommand(productsListCmd)
	productsCmd.AddCommand(enableProductCmd)
	productsCmd.AddCommand(disableProductCmd)
}
