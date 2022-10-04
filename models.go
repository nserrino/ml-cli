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
		"The local path of the variant to upload as the base variant for the model. "+
			"Supports tensorflow serving and openvino format.")
	CreateModel.Flags().String("grpc-port", "31312", "The GRPC port for the client to access the model.")
	CreateModel.Flags().String("http-port", "31313", "The HTTP port for the client to access the model.")
	CreateModel.Flags().String("image-destination", "",
		"The destination to push the container image for the model server")
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
		modelName := args[0]
		variantName := "base"

		localBaseVariantPath, err := cmd.Flags().GetString("from-local")
		if err != nil || localBaseVariantPath == "" {
			fmt.Println("`mlm models create` requires flag --from-local")
			os.Exit(1)
		}
		httpPort, err := cmd.Flags().GetString("http-port")
		if err != nil || httpPort == "" {
			fmt.Println("`mlm models create` requires flag --http-port")
			os.Exit(1)
		}
		grpcPort, err := cmd.Flags().GetString("grpc-port")
		if err != nil || httpPort == "" {
			fmt.Println("`mlm models create` requires flag --grpc-port")
			os.Exit(1)
		}
		imageDestination, err := cmd.Flags().GetString("image-destination")
		if err != nil || imageDestination == "" {
			fmt.Println("`mlm models create` requires flag --image-destination")
			os.Exit(1)
		}

		createVariant(modelName, variantName, localBaseVariantPath, imageDestination, nil,
			grpcPort, httpPort)
	},
}
