package main

import (
	"fmt"
	"os/exec"
)

func labelName(modelName string) string {
	return fmt.Sprintf("mlm.variant.%s", modelName)
}

func addNodeLabels(modelName, modelVariant string, nodes []string) error {
	if len(nodes) == 0 {
		cmd := exec.Command("kubectl", "label", "nodes", "--all", "--overwrite",
			fmt.Sprintf("%s=%s", labelName(modelName), modelVariant))
		return cmd.Run()
	}
	for _, nodeName := range nodes {
		cmd := exec.Command("kubectl", "label", "nodes", nodeName, "--overwrite",
			fmt.Sprintf("%s=%s", labelName(modelName), modelVariant))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
