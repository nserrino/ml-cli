package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// List all model versions currently stored in this cluster.
	Models.AddCommand(ListModels)
	// Get more information about a particular model version in this cluster.
	Models.AddCommand(GetModel)
	// Create a new model version in this cluster by uploading from local files.
	Models.AddCommand(CreateModel)

	CreateModel.Flags().String("local-model-path", "", "The local path of the model to upload to the cluster")
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

// GetModel gets information about a particular ML deployment.
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

// CreateModel updates the given deployment to use a different model versino.
var CreateModel = &cobra.Command{
	Use:   "create",
	Short: "Add a new model to the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm models create` requires 1 positional argument (model name)")
			os.Exit(1)
		}
		localModelPath, err := cmd.Flags().GetString("local-model-path")
		if err != nil || localModelPath == "" {
			fmt.Println("`mlm models create` requires flag --local-model-path")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Creating a new ML model %s which is loading from the directory %s\n", args[0], localModelPath)
	},
}
