package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// List all ML deployments (N deployments can back 1 ML app) on this cluster.
	Deployments.AddCommand(ListDeployments)
	// Get more info about a particular ML deployment on this cluster.
	Deployments.AddCommand(GetDeployment)
	// Update a particular deployment (to use a different model version).
	Deployments.AddCommand(UpdateDeployment)
	// Fork a particular deployment to use a different model version (based on node label.)
	Deployments.AddCommand(ForkDeployment)

	UpdateDeployment.Flags().String("model", "", "The name of the model to update the deployment to use")
	ForkDeployment.Flags().String("node-label", "", "The value the node label should use for the forked deployment")
}

// Deployments is the parent of all deployment-related commands.
var Deployments = &cobra.Command{
	Use:   "deployments",
	Short: "Commands related to ML deployments",
}

// ListDeployments lists all of the ML deployments in the current cluster.
var ListDeployments = &cobra.Command{
	Use:   "list",
	Short: "List ML deployments in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PLACEHOLDER: Listing ML deployments")
	},
}

// GetDeployment gets information about a particular ML deployment.
var GetDeployment = &cobra.Command{
	Use:   "get",
	Short: "Get more info about a particular ML deployment in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm deployments get` requires 1 argument (deployment name)")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Getting the ML deployment %s\n", args[0])
	},
}

// UpdateDeployment updates the given deployment to use a different model versino.
var UpdateDeployment = &cobra.Command{
	Use:   "update",
	Short: "Update a deployment to use a different model version",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm deployments update` requires 1 positional argument (deployment name)")
			os.Exit(1)
		}
		newModel, err := cmd.Flags().GetString("model")
		if err != nil || newModel == "" {
			fmt.Println("`mlm deployments update` requires flag --model")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Updating the ML deployment %s to use model %s\n", args[0], newModel)
	},
}

// ForkDeployment forks a particular deployment, with a different node label selector.
var ForkDeployment = &cobra.Command{
	Use:   "fork",
	Short: "Fork an existing deployment, with a new node label selector",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			fmt.Println("`mlm deployments update` requires 2 positional argument (input deployment name, fork deployment name)")
			os.Exit(1)
		}
		labelSelector, err := cmd.Flags().GetString("node-label")
		if err != nil || labelSelector == "" {
			fmt.Println("`mlm deployments fork` requires flag --node-label")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Forking the ML deployment %s to %s using node-label %s\n", args[0], args[1], labelSelector)
	},
}
