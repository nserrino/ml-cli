package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// List all ML models (N models can back 1 ML app) on this cluster.
	Models.AddCommand(ListModels)
	// Get more info about a particular ML model on this cluster.
	Models.AddCommand(GetModel)
	// Update a particular model (to use a different model version).
	Models.AddCommand(CreateModel)

	CreateModel.Flags().String("from-local", "",
		"The local path of the variant to upload as the base variant for the model")
}

// Models is the parent of all model-related commands.
var Models = &cobra.Command{
	Use:   "models",
	Short: "Commands related to ML models",
}

// ListModels lists all of the ML models in the current cluster.
var ListModels = &cobra.Command{
	Use:   "list",
	Short: "List ML models in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PLACEHOLDER: Listing ML models")
	},
}

// GetModel gets information about a particular ML model.
var GetModel = &cobra.Command{
	Use:   "get",
	Short: "Get more info about a particular ML model in the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm models get` requires 1 argument (model name)")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Getting the ML model %s\n", args[0])
	},
}

// CreateModel creates a new model (with a base variant uploaded from local)
var CreateModel = &cobra.Command{
	Use:   "create",
	Short: "Create a new model with a base variant",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm models create` requires 1 argument (model name)")
			os.Exit(1)
		}
		localBaseVariantPath, err := cmd.Flags().GetString("from-local")
		if err != nil || localBaseVariantPath == "" {
			fmt.Println("`mlm models create` requires flag --from-local")
			os.Exit(1)
		}

		modelVersionName := fmt.Sprintf("%s-base", args[0])

		// Create the shared volume for the model
		err = copyFromLocalToSharedVolume(localBaseVariantPath, modelVersionName)
		if err != nil {
			fmt.Println("Error copying from local to shared volume:" + err.Error())
			os.Exit(1)
		}

		// Create the daemonset and service for the model
		err = deployModelFromSharedVolume(args[0], "base", modelVersionName)
		if err != nil {
			fmt.Println("Error deploying model from shared volume:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("Successfully deployed model from shared volume")

		// Add node labels for the base variant
		err = addNodeLabels(args[0], "base", nil)
		if err != nil {
			fmt.Println("Error adding node labels for base variant:" + err.Error())
			os.Exit(1)
		}
		fmt.Println("Successfully added node labels for base variant")
	},
}
