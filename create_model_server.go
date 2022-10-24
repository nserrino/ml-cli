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

const tensorflowModelServerType = "tensorflow"
const openvinoModelServerType = "openvino"

const tensorflowDockerfileTmpl = `# syntax=docker/dockerfile:1
FROM tensorflow/serving:2.8.0

RUN  apt-get update \
  && apt-get install -y wget \
  && rm -rf /var/lib/apt/lists/*

COPY model /models/$MODEL_NAME/1

ENV MODEL_NAME=$MODEL_NAME

ENTRYPOINT ["/usr/bin/tf_serving_entrypoint.sh"]`

const openvinoDockerfileTmpl = `# syntax=docker/dockerfile:1
FROM openvino/model_server:2022.1

ADD model /models/$MODEL_NAME/1

CMD ["/ovms/bin/ovms", "--model_path", "/models/$MODEL_NAME", "--model_name", "$MODEL_NAME", "--port", "8500", "--rest_port", "8501", "--shape", "auto"]`

func createDockerfile(modelServer, modelName, filePath string) error {
	var fileContents string
	if modelServer == openvinoModelServerType {
		fileContents = strings.Replace(openvinoDockerfileTmpl, "$MODEL_NAME", modelName, -1)
	} else if modelServer == tensorflowModelServerType {
		fileContents = strings.Replace(tensorflowDockerfileTmpl, "$MODEL_NAME", modelName, -1)
	} else {
		return fmt.Errorf("Unsupported model server type %s", modelServer)
	}

	fmt.Printf("Building %s type model server\n", modelServer)

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

func getModelServerTypeFromModelFormat(absModelPath string) (string, error) {
	// Determine whether to use tensorflow serving or openvino model server based on input
	modelDirStat, _ := os.Lstat(absModelPath)

	// Validate that the path is a directory
	if isDir := modelDirStat.IsDir(); !isDir {
		return "", fmt.Errorf("Expected model path to point to directory, received file: %s", absModelPath)
	}

	srcDirContent, _ := os.ReadDir(absModelPath)

	// ovms
	hasXml := false
	hasBin := false
	// tensorflow
	hasSavedModel := false

	for _, contentFile := range srcDirContent {
		if contentFile.IsDir() {
			continue
		}
		if contentFile.Name() == "saved_model.pb" {
			hasSavedModel = true
		}
		ext := filepath.Ext(contentFile.Name())
		if ext == ".bin" {
			hasBin = true
		} else if ext == ".xml" {
			hasXml = true
		}
	}

	if hasSavedModel {
		return tensorflowModelServerType, nil
	}
	if hasXml && hasBin {
		return openvinoModelServerType, nil
	}
	return "", fmt.Errorf("Model directory did not contain expected format. " +
		"(Tensorflow format: requires saved_model.pb. OpenVino format: requires xml and bin files.)")
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

	modelServer, err := getModelServerTypeFromModelFormat(absModelPath)
	if err != nil {
		return err
	}

	err = createDockerfile(modelServer, modelName, "Dockerfile")
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
  namespace: $NAMESPACE
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
          hostPort: $HOST_HTTP_PORT
          name: http
        - containerPort: 8500
          hostPort: $HOST_GRPC_PORT
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

func deployModelServer(modelName, variantName, namespace, imageLocation, grpcPort, httpPort string) error {
	modelDaemonSetName := fmt.Sprintf("%s-%s", modelName, variantName)
	nodeLabelKey := labelName(modelName)
	modelServerVariantImage := imageLocation

	kubernetesResources := kubernetesModelServerTmpl
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_NAME", modelName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$VARIANT_NAME", variantName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_DAEMONSET_NAME", modelDaemonSetName, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$NODE_LABEL_KEY", nodeLabelKey, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$MODEL_SERVER_VARIANT_IMAGE", modelServerVariantImage, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$HOST_GRPC_PORT", grpcPort, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$HOST_HTTP_PORT", httpPort, -1)
	kubernetesResources = strings.Replace(kubernetesResources, "$NAMESPACE", namespace, -1)

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
