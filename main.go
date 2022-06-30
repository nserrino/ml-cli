package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Commands related to ML apps running on this cluster.
	RootCmd.AddCommand(Apps)
	// Commands related to ML deployments running on this cluster (N deployments can back 1 app).
	RootCmd.AddCommand(Deployments)
	// Commands related to ML model versions running on this cluster.
	RootCmd.AddCommand(Models)
}

// RootCmd is the base command for Cobra.
var RootCmd = &cobra.Command{
	Use:   "mlm",
	Short: "ML Manager CLI",
	Long:  `The ML Manager CLI interface.`,
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Printf("Error executing command: %s\n", err.Error())
	}
}
