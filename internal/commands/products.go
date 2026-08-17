package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/neetozone/neeto-auth-cli/internal/output"
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Inspect products and their available roles",
}

var productsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List products available in the workspace and the roles each one exposes",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}
		data, err := c.Get("/products", nil)
		if err != nil {
			return fmt.Errorf("failed to fetch products: %w", err)
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
		if _, err := c.Patch(fmt.Sprintf("/products/%s", product), map[string]any{
			"enabled": enabled,
		}); err != nil {
			if isProductNotFoundError(err) {
				return fmt.Errorf("product %q not found. Run \"neetoauth products list\" to see available products", product)
			}
			return fmt.Errorf("failed to %s product %q: %w", action, product, err)
		}
		result, _ := json.Marshal(map[string]string{
			"product": product,
			"status":  verbed,
		})
		printResource(result, nil)
		return nil
	}
}

func isProductNotFoundError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") || strings.Contains(msg, "404")
}
func productBreadcrumbs() []output.Breadcrumb {
	return []output.Breadcrumb{
		{Label: "Use a role when inviting", Command: "neetoauth users create --email <email> --role non_owner --app <product>:<role>"},
	}
}

func init() {
	rootCmd.AddCommand(productsCmd)
	productsCmd.AddCommand(productsListCmd)
	productsCmd.AddCommand(enableProductCmd)
	productsCmd.AddCommand(disableProductCmd)
}
