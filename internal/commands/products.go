package commands

import (
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
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "Use a role when inviting", Command: "neetoauth users create --email <email> --role non_owner --app <product>:<role>"},
		}

		printList(data, "products", breadcrumbs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(productsCmd)
	productsCmd.AddCommand(productsListCmd)
}
