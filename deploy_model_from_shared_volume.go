package main

import (
// "fmt"
// "io/ioutil"
// "os"
// "os/exec"
// "regexp"
// "strings"
)

func deployModelFromSharedVolume(modelName, modelVariantName, sharedVolumeName string) error {
	err := validateKubernetesName(modelName)
	if err != nil {
		return err
	}
	return nil
}
