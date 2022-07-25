package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	// List all variant versions currently stored in this cluster.
	Variants.AddCommand(ListVariants)
	// Get more information about a particular variant version for a cluster.
	Variants.AddCommand(GetVariant)
	// Create a new variant version for a given model by uploading from local files.
	Variants.AddCommand(CreateVariant)
	// Update a variant to point to a new version.
	Variants.AddCommand(UpdateVariant)

	// Flags for all Variant commands
	Variants.PersistentFlags().String("model", "", "The name of the model")

	// Flags for create
	CreateVariant.Flags().String("from-local", "", "The local path of the variant to upload")
	CreateVariant.Flags().String("nodes", "", "A list of nodes to use this variant for")

	// Flags for update
	UpdateVariant.Flags().String("from-local", "", "The local path of the variant to upload")
}

// Variants is the parent of all variant-related commands.
var Variants = &cobra.Command{
	Use:   "variants",
	Short: "Commands related to ML variants",
}

func getModelOrFail(cmd *cobra.Command) string {
	modelName, err := cmd.Flags().GetString("model")
	if err != nil || modelName == "" {
		fmt.Println("`mlm variants` requires flag --model")
		os.Exit(1)
	}
	return modelName
}

// ListVariants lists all of the ML variants in the current cluster.
var ListVariants = &cobra.Command{
	Use:   "list",
	Short: "List ML variants for a given model",
	Run: func(cmd *cobra.Command, args []string) {
		modelName := getModelOrFail(cmd)
		fmt.Printf("PLACEHOLDER: Listing ML variants for model %s\n", modelName)
	},
}

// GetVariant gets information about a particular ML variant.
var GetVariant = &cobra.Command{
	Use:   "get",
	Short: "Get more info about a particular ML variant for a model",
	Run: func(cmd *cobra.Command, args []string) {
		modelName := getModelOrFail(cmd)
		if len(args) != 1 {
			fmt.Println("`mlm variants get` requires 1 argument (variant name)")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Getting the ML variant %s for model %s\n", args[0], modelName)
	},
}

// CreateVariant updates the given variant to use a different model version.
var CreateVariant = &cobra.Command{
	Use:   "create",
	Short: "Add a new variant to the model",
	Run: func(cmd *cobra.Command, args []string) {
		modelName := getModelOrFail(cmd)
		if len(args) != 1 {
			fmt.Println("`mlm variants create` requires 1 positional argument (variant name)")
			os.Exit(1)
		}
		localVariantPath, err := cmd.Flags().GetString("from-local")
		if err != nil || localVariantPath == "" {
			fmt.Println("`mlm variants create` requires flag --from-local")
			os.Exit(1)
		}
		nodeList, err := cmd.Flags().GetString("nodes")
		if err != nil || nodeList == "" {
			fmt.Println("`mlm variants create` requires flag --nodes")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Creating a new ML variant %s for nodes %s on model %s from the directory %s\n",
			args[0], nodeList, modelName, localVariantPath)
	},
}

// UpdateVariant updates the given variant to point to a new version of the model.
var UpdateVariant = &cobra.Command{
	Use:   "update",
	Short: "Update an existing variant",
	Run: func(cmd *cobra.Command, args []string) {
		modelName := getModelOrFail(cmd)
		if len(args) != 1 {
			fmt.Println("`mlm variants update` requires 1 positional argument (variant name)")
			os.Exit(1)
		}
		localVariantPath, err := cmd.Flags().GetString("from-local")
		if err != nil || localVariantPath == "" {
			fmt.Println("`mlm variants update` requires flag --from-local")
			os.Exit(1)
		}
		fmt.Printf("PLACEHOLDER: Updating ML variant %s for model %s which is loading from the directory %s\n",
			args[0], modelName, localVariantPath)
	},
}
