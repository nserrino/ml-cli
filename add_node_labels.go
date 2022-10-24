package main

import (
	"fmt"
	"os/exec"
)

func labelName(modelName string) string {
	return fmt.Sprintf("mlm.variant.%s", modelName)
}

func addNodeLabelsToAll(modelName, modelVariant string) error {
	cmd := exec.Command("kubectl", "label", "nodes", "--all", "--overwrite",
		fmt.Sprintf("%s=%s", labelName(modelName), modelVariant))
	return cmd.Run()
}

func addNodeLabelsBySelector(modelName, modelVariant string, selector string) error {
	cmd := exec.Command("kubectl", "label", "nodes", "--selector", selector, "--overwrite",
		fmt.Sprintf("%s=%s", labelName(modelName), modelVariant))
	return cmd.Run()
}

func addNodeLabelsByList(modelName, modelVariant string, nodes []string) error {
	for _, nodeName := range nodes {
		cmd := exec.Command("kubectl", "label", "nodes", nodeName, "--overwrite",
			fmt.Sprintf("%s=%s", labelName(modelName), modelVariant))
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
