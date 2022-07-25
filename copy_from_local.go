package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func createNamespaceIfNotExists() error {
	cmd := exec.Command("kubectl", "create", "namespace", "mlm")
	return cmd.Run()
}

const sharedVolumeReadWriteOnceTmpl = `---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: $MODEL_VERSION_NAME-init
  namespace: mlm
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 100M
`

const sharedVolumeReadOnlyManyTmpl = `---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: $MODEL_VERSION_NAME
  namespace: mlm
spec:
  accessModes:
  - ReadOnlyMany
  dataSource:
    name: $MODEL_VERSION_NAME-init
    kind: PersistentVolumeClaim  
  resources:
    requests:
      storage: 100M
`

const validationRegex = "^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"

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

func createSharedVolume(sharedVolumeTmpl, modelVersionName string) error {
	sharedVolume := strings.Replace(sharedVolumeTmpl, "$MODEL_VERSION_NAME", modelVersionName, -1)
	fileName, err := createTempFile(sharedVolume)
	if fileName == "" || err != nil {
		return err
	}
	defer os.Remove(fileName)
	cmd := exec.Command("kubectl", "apply", "-f", fileName)
	if o, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Error creating shared volume: " + string(o))
		return err
	}
	return nil
}

const tmpPodTmpl = `---
apiVersion: v1
kind: Pod
metadata:
  name: $MODEL_VERSION_NAME
  namespace: mlm
spec:
  containers:
  - image: alpine:3.2
    command:
      - /bin/sh
      - "-c"
      - "sleep 5m"
    imagePullPolicy: IfNotPresent
    name: alpine
    volumeMounts:
    - mountPath: /models
      name: $MODEL_VERSION_NAME-init
  volumes:
  - name: $MODEL_VERSION_NAME-init
    persistentVolumeClaim:
      claimName: $MODEL_VERSION_NAME-init
  restartPolicy: Always
`

func createTmpPodWithSharedVolume(modelVersionName string) error {
	// kubectl cp tmp-dir/ssd_mobilenet_v1_coco_2017_11_17/saved_model alpine:/models -n mlm
	tmpPod := strings.Replace(tmpPodTmpl, "$MODEL_VERSION_NAME", modelVersionName, -1)
	fileName, err := createTempFile(tmpPod)
	if fileName == "" || err != nil {
		return err
	}
	defer os.Remove(fileName)
	cmd := exec.Command("kubectl", "apply", "-f", fileName)
	if o, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Error creating tmp pod with shared volume: " + string(o))
		return err
	}

	cmd = exec.Command("kubectl", "-n", "mlm", "wait", "--for=condition=Ready=true",
		fmt.Sprintf("pod/%s", modelVersionName))
	fmt.Println(fmt.Sprintf("Waiting for pod %s to be ready", modelVersionName))
	if o, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Waiting for pod to be ready: " + string(o))
		return err
	}
	return nil
}

func copyToSharedVolumeViaTmpPod(localPath, modelVersionName string) error {
	// kubectl cp $localPath alpine:/models -n mlm
	cmd := exec.Command("kubectl", "cp", "-n", "mlm", localPath, fmt.Sprintf("%s:/models/model", modelVersionName))
	if o, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Error copying to shared volume via tmp pod: " + string(o))
		return err
	}
	return nil
}

func validateKubernetesName(input string) error {
	r, err := regexp.Compile(input)
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

func copyFromLocalToSharedVolume(localPath string, modelVersionName string) error {
	fmt.Println("Creating namespace mlm if it doesn't exist")
	createNamespaceIfNotExists()

	fmt.Println("Validating resource name " + modelVersionName)
	err := validateKubernetesName(modelVersionName)
	if err != nil {
		return err
	}

	// Create the ReadWriteOnce shared volume
	fmt.Println("Creating the tmp ReadWriteOnce shared volume to initialize " + modelVersionName)
	err = createSharedVolume(sharedVolumeReadWriteOnceTmpl, modelVersionName)
	if err != nil {
		return err
	}

	fmt.Println("Creating the tmp pod to initialize data for " + modelVersionName)
	err = createTmpPodWithSharedVolume(modelVersionName)
	if err != nil {
		return err
	}

	fmt.Println("Copying " + localPath + " to the tmp shared volume for " + modelVersionName)
	err = copyToSharedVolumeViaTmpPod(localPath, modelVersionName)
	if err != nil {
		return err
	}

	// Create the ReadOnlyMany shared volume
	fmt.Println("Creating the ReadOnlyMany shared volume to finalize " + modelVersionName)
	return createSharedVolume(sharedVolumeReadOnlyManyTmpl, modelVersionName)
}
