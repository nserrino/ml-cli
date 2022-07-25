package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	// Commands related to ML models running on this cluster.
	RootCmd.AddCommand(Models)
	// Commands related to ML model versions running on this cluster (N variants can back 1 model).
	RootCmd.AddCommand(Variants)
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
