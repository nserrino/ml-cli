package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// List all ML apps running on this cluster.
	Apps.AddCommand(ListApps)
	// Get more info about a particular ML app running on this cluster.
	Apps.AddCommand(GetApp)
}

// Apps is the parent of all app-related commands.
var Apps = &cobra.Command{
	Use:   "apps",
	Short: "Commands related to ML apps",
}

// ListApps lists all of the ML apps running on the cluster.
var ListApps = &cobra.Command{
	Use:   "list",
	Short: "List ML apps in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PLACEHOLDER: Listing ML apps")
	},
}

// GetApp gets information about a particular ML app.
var GetApp = &cobra.Command{
	Use:   "get",
	Short: "Get more info about a particular ML app in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm apps get` requires 1 argument (app name)")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Getting the ML app %s\n", args[0])
	},
}
