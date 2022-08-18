package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

const dockerfileTmpl = `# syntax=docker/dockerfile:1
FROM tensorflow/serving:2.8.0

RUN  apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*

COPY model /models/$MODEL_NAME/1

ENV MODEL_NAME=$MODEL_NAME

ENTRYPOINT ["/usr/bin/tf_serving_entrypoint.sh"]`

func createDockerfile(modelName, filePath string) error {
	fileContents := strings.Replace(dockerfileTmpl, "$MODEL_NAME", modelName, -1)

	tmpFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	text := []byte(fileContents)
	if _, err := tmpFile.Write(text); err != nil {
		return fmt.Errorf("Failed to write to Dockerfile", err)
	}

	return nil
}

func createAndPushModelServerImage(modelName, modelPath, imageDestination string) error {
	prevDir, _ := os.Getwd()
	absModelPath, err := filepath.Abs(modelPath)
	if err != nil {
		return err
	}

	dir, err := ioutil.TempDir(os.TempDir(), modelName)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Could not move into the directory (%s)\n", dir)
	}

	cpCmd := exec.Command("cp", "-r", absModelPath, "model")
	if o, err := cpCmd.CombinedOutput(); err != nil {
		fmt.Println("Error copying model into temporary directory:" + string(o))
		return err
	}

	err = createDockerfile(modelName, "Dockerfile")
	if err != nil {
		return err
	}

	buildCmd := exec.Command("docker", "build", ".", "-t", imageDestination)
	if o, err := buildCmd.CombinedOutput(); err != nil {
		fmt.Println("Error building docker image:" + string(o))
		return err
	}

	pushCmd := exec.Command("docker", "push", imageDestination)
	if o, err := pushCmd.CombinedOutput(); err != nil {
		fmt.Println("Error pushing docker image:" + string(o))
		return err
	}
	fmt.Println("Successfully built image: " + imageDestination)

	if err := os.Chdir(prevDir); err != nil {
		return fmt.Errorf("Could not move into the directory (%s)\n", prevDir)
	}
	fmt.Println("Successfully pushed image: " + imageDestination)

	return nil
}

const validationRegex = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"

func validateKubernetesName(input string) error {
	r, err := regexp.Compile(validationRegex)
	if err != nil {
		return fmt.Errorf("Failed to compile regexp %s", validationRegex)
	}
	if !r.MatchString(input) {
		return fmt.Errorf("Input '%s' doesn't meet guidelines: must consist of lower case "+
			"alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character",
			input)
	}
	return nil
}

const kubernetesModelServerTmpl = `---
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: $MODEL_DAEMONSET_NAME
  namespace: mlm
  labels:
    mlm.model.name: $MODEL_NAME
    mlm.model.variant: $VARIANT_NAME
    app: $MODEL_DAEMONSET_NAME
spec:
  selector:
    matchLabels:
      app: $MODEL_DAEMONSET_NAME
  template:
    metadata:
      labels:
        mlm.model.name: $MODEL_NAME
        mlm.model.variant: $VARIANT_NAME      
        app: $MODEL_DAEMONSET_NAME
    spec:
      nodeSelector:
        $NODE_LABEL_KEY: $VARIANT_NAME
      containers:
      - name: mlm
        image: $MODEL_SERVER_VARIANT_IMAGE
        ports:
        - containerPort: 8501
          hostPort: 50311
          name: http
        - containerPort: 8500
          hostPort: 50310
          name: grpc
`

func createTempFile(fileContents string) (string, error) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return "", fmt.Errorf("Cannot create temporary file", err)
	}

	text := []byte(fileContents)
	if _, err = tmpFile.Write(text); err != nil {
		return tmpFile.Name(), fmt.Errorf("Failed to write to temporary file", err)
	}

	return tmpFile.Name(), nil
}

func deployModelServer(modelName, variantName, imageLocation string) error {
	modelDaemonSetName := fmt.Sprintf("%s-%s", modelName, variantName)
	nodeLabelKey := labelName(modelName)
	modelServerVariantImage := imageLocation

	kubernetesResources := kubernetesModelServerTmpl
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_NAME", modelName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$VARIANT_NAME", variantName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_DAEMONSET_NAME", modelDaemonSetName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$NODE_LABEL_KEY", nodeLabelKey, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_SERVER_VARIANT_IMAGE", modelServerVariantImage, -1)

	fileName, err := createTempFile(kubernetesResources)
	if fileName == "" || err != nil {
		return err
	}
	defer os.Remove(fileName)
	cmd := exec.Command("kubectl", "apply", "-f", fileName)
	if o, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Error creating model resources: " + string(o))
		return err
	}
	return nil
}
