package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage workspace team members",
}

var usersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active team members in the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		params := paginationParams(cmd)

		data, err := c.Get("/users", params)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "Invite a new member", Command: "neetoauth users create --email <email> --role <role>"},
			{Label: "Remove a member", Command: "neetoauth users delete <email>"},
		}

		printList(data, "users", breadcrumbs)
		return nil
	},
}

var usersCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Invite a new team member to the workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		user := map[string]interface{}{}

		jsonFile, _ := cmd.Flags().GetString("json-file")
		if jsonFile != "" {
			fileData, err := readJSONFile(jsonFile)
			if err != nil {
				return err
			}
			if raw, ok := fileData["user"]; ok {
				nested, ok := raw.(map[string]interface{})
				if !ok {
					return fmt.Errorf("Invalid \"user\" field in %s: expected an object.", jsonFile)
				}
				user = nested
			} else {
				user = fileData
			}
		}

		email, _ := cmd.Flags().GetString("email")
		if email != "" {
			user["email"] = email
		}
		if _, ok := user["email"]; !ok {
			return fmt.Errorf("--email is required.")
		}

		role, _ := cmd.Flags().GetString("role")
		if role != "" {
			user["role"] = role
		}
		if _, ok := user["role"]; !ok {
			return fmt.Errorf("--role is required (e.g. owner, non_owner).")
		}

		if firstName, _ := cmd.Flags().GetString("first-name"); firstName != "" {
			user["first_name"] = firstName
		}
		if lastName, _ := cmd.Flags().GetString("last-name"); lastName != "" {
			user["last_name"] = lastName
		}

		appFlags, _ := cmd.Flags().GetStringSlice("app")
		if len(appFlags) > 0 {
			apps := make([]map[string]string, 0, len(appFlags))
			for _, raw := range appFlags {
				name, appRole, ok := strings.Cut(raw, ":")
				if !ok || name == "" || appRole == "" {
					return fmt.Errorf("Invalid --app value %q (expected name:role).", raw)
				}
				apps = append(apps, map[string]string{
					"name": strings.TrimSpace(name),
					"role": strings.TrimSpace(appRole),
				})
			}
			user["apps"] = apps
		}

		body := map[string]interface{}{"user": user}

		data, err := c.Post("/users", body)
		if err != nil {
			return err
		}

		breadcrumbs := []output.Breadcrumb{
			{Label: "List members", Command: "neetoauth users list"},
			{Label: "Remove this member", Command: fmt.Sprintf("neetoauth users delete %v", user["email"])},
		}

		printActionResult(data, breadcrumbs)
		return nil
	},
}

var usersDeleteCmd = &cobra.Command{
	Use:   "delete <email>",
	Short: "Remove a team member from the workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := getClient(cmd)
		if err != nil {
			return err
		}

		if err := c.Delete(fmt.Sprintf("/users/%s", url.PathEscape(args[0]))); err != nil {
			return err
		}

		printMessage(fmt.Sprintf("Removed %s from the workspace.", args[0]))
		return nil
	},
}

func init() {
	register(func(root *cobra.Command) { root.AddCommand(usersCmd) })

	usersCmd.AddCommand(usersListCmd)
	addPaginationFlags(usersListCmd)

	usersCmd.AddCommand(usersCreateCmd)
	usersCreateCmd.Flags().String("email", "", "Member email address")
	usersCreateCmd.Flags().String("role", "", "Organization role, e.g. owner or non_owner")
	usersCreateCmd.Flags().String("first-name", "", "Member first name")
	usersCreateCmd.Flags().String("last-name", "", "Member last name")
	usersCreateCmd.Flags().StringSlice("app", nil, "Per-app role assignment as name:role (repeatable)")
	usersCreateCmd.Flags().String("json-file", "", "Path to a JSON file with the full user payload")
	markFlagsRequired(usersCreateCmd, "email", "role")
	allowJSONFileToSatisfyRequiredFlags(usersCreateCmd)

	usersCmd.AddCommand(usersDeleteCmd)
}
