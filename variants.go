package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	// List all variant versions currently stored in this cluster.
	Variants.AddCommand(ListVariants)
	// Get more information about a particular variant version for a cluster.
	Variants.AddCommand(GetVariant)
	// Create a new variant version for a given model by uploading from local files.
	Variants.AddCommand(CreateVariant)

	// Flags for all Variant commands
	Variants.PersistentFlags().String("model", "", "The name of the model")

	// Flags for create
	CreateVariant.Flags().StringP("namespace", "n", "mlm", "The namespace for the variant")
	CreateVariant.Flags().String("from-local", "", "The local path of the variant to upload. "+
		"Supports tensorflow serving and openvino format.")
	CreateVariant.Flags().String("nodes", "", "A list of nodes to use this variant for")
	CreateVariant.Flags().String("node-selector", "", "A K8s node selector to identify nodes to use this variant for")
	CreateVariant.Flags().String("grpc-port", "31312", "The GRPC port for the client to access the model.")
	CreateVariant.Flags().String("http-port", "31313", "The HTTP port for the client to access the model.")
	CreateVariant.Flags().String("image-destination", "",
		"The destination to push the container image for the model server")
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

func createVariant(modelName, variantName, namespace, localBaseVariantPath string,
	imageDestination string, nodes []string, nodeSelector, grpcPort, httpPort string) {
	err := createAndPushModelServerImage(modelName, localBaseVariantPath, imageDestination)
	if err != nil {
		fmt.Println("Error creating and pushing model server:" + err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully created and pushed model server image")

	err = deployModelServer(modelName, variantName, namespace, imageDestination, grpcPort, httpPort)
	if err != nil {
		fmt.Println("Error deploying model server:" + err.Error())
		os.Exit(1)
	}
	fmt.Println("Successfully deployed model server")

	// Add node labels based on the user criteria
	if nodeSelector != "" {
		err = addNodeLabelsBySelector(modelName, variantName, nodeSelector)
	} else if len(nodes) > 0 {
		err = addNodeLabelsByList(modelName, variantName, nodes)
	} else {
		err = addNodeLabelsToAll(modelName, variantName)
	}
	if err != nil {
		fmt.Printf("Error adding node labels for variant %s: %s\n", variantName, err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully added node labels for %s variant\n", variantName)
	fmt.Printf("Deployed model to nodes on ports %s (gRPC), %s (HTTP)\n", grpcPort, httpPort)
	fmt.Printf("Model can be accessed at HTTP endpoint `<HOST IP>:%s/v1/models/%s:predict`", httpPort, modelName)
	fmt.Println("Run `kubectl -n mlm get daemonset` for more information")
}

// CreateVariant updates the given variant to use a different model version.
var CreateVariant = &cobra.Command{
	Use:   "create",
	Short: "Add a new variant to the model",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("`mlm variants create` requires 1 positional argument (variant name)")
			os.Exit(1)
		}
		variantName := args[0]
		modelName := getModelOrFail(cmd)

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil || namespace == "" {
			fmt.Println("`mlm variants create` requires flag --namespace")
			os.Exit(1)
		}
		localVariantPath, err := cmd.Flags().GetString("from-local")
		if err != nil || localVariantPath == "" {
			fmt.Println("`mlm variants create` requires flag --from-local")
			os.Exit(1)
		}
		nodeList, _ := cmd.Flags().GetString("nodes")
		nodeSelector, _ := cmd.Flags().GetString("node-selector")
		if nodeSelector == "" && nodeList == "" {
			fmt.Println("`mlm variants create` requires one of the flags: --nodes or --node-selector")
			os.Exit(1)
		}		
		imageDestination, err := cmd.Flags().GetString("image-destination")
		if err != nil || imageDestination == "" {
			fmt.Println("`mlm variants create` requires flag --image-destination")
			os.Exit(1)
		}
		httpPort, err := cmd.Flags().GetString("http-port")
		if err != nil || httpPort == "" {
			fmt.Println("`mlm variants create` requires flag --http-port")
			os.Exit(1)
		}
		grpcPort, err := cmd.Flags().GetString("grpc-port")
		if err != nil || httpPort == "" {
			fmt.Println("`mlm variants create` requires flag --grpc-port")
			os.Exit(1)
		}
		createVariant(modelName, variantName, namespace, localVariantPath, imageDestination,
			strings.Split(nodeList, ","), nodeSelector, grpcPort, httpPort)
	},
}
