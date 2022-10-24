package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	Apply.Flags().StringP("filename", "f", "", "The path to the YAML file to apply.")
}

type baseVariantSpec struct {
	ImageDestination string `yaml:"imageDestination"`
	SourceModelPath  string `yaml:"sourceModelPath"`
}

type variantSpec struct {
	Name             string   `yaml:"name"`
	Nodes            []string `yaml:"nodes",flow`
	NodeSelector     string   `yaml:"nodeSelector"`
	ImageDestination string   `yaml:"imageDestination"`
	SourceModelPath  string   `yaml:"sourceModelPath"`
}

type modelSpec struct {
	Model       string           `yaml:"model"`
	Namespace   string           `yaml:"namespace"`
	HttpPort    string           `yaml:"httpPort"`
	GrpcPort    string           `yaml:"grpcPort"`
	BaseVariant *baseVariantSpec `yaml:"baseVariant"`
	Variants    []variantSpec    `yaml:"variants",flow`
}

// Apply deploys new models based on the state specified in the config file.
var Apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply model(s) to your clusters based on a YAML file",
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("filename")
		if err != nil || filename == "" {
			fmt.Println("`mlm apply` requires flag -f")
			os.Exit(1)
		}

		yamlFile, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Printf("Could not read file %s: %v\n", filename, err)
			os.Exit(1)
		}

		var m modelSpec
		err = yaml.Unmarshal(yamlFile, &m)
		if err != nil {
			fmt.Printf("Unmarshal failed on %s: %v", filename, err)
			os.Exit(1)
		}

		if m.BaseVariant != nil {
			fmt.Printf("Building and deploying base variant of %s to namespace %s\n",
				m.Model, m.Namespace)
			createVariant(m.Model, "base", m.Namespace, m.BaseVariant.SourceModelPath,
				m.BaseVariant.ImageDestination, nil, "", m.GrpcPort, m.HttpPort)
		}

		for _, v := range m.Variants {
			fmt.Printf("Building and deploying %s variant of %s to namespace %s\n",
				v.Name, m.Model, m.Namespace)
			createVariant(m.Model, v.Name, m.Namespace, v.SourceModelPath, v.ImageDestination,
				v.Nodes, v.NodeSelector, m.GrpcPort, m.HttpPort)
		}
	},
}
