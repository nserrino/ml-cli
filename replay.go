package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	CreateModel.Flags().String("http-port", "31312", "The HTTP port(s) for the client to access the model(s).")
	CreateModel.Flags().String("from-file", "", "A path to the captured JSON requests to replay.")
}

// Replay replays requests from a given model to another model.
var Replay = &cobra.Command{
	Use:   "replay",
	Short: "Replay a request to another instance of a model (local or remote)",
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
