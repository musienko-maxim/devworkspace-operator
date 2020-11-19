//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package client

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
)

func (w *K8sClient) OcApplyWorkspace(filePath string) (err error) {
	cmd := exec.Command("oc", "apply", "--namespace", config.Namespace, "-f", filePath)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if strings.Contains(output, "failed calling webhook") {
		fmt.Println("Seems DevWorkspace Webhook Server is not ready yet. Will retry in 2 seconds. Cause: " + output)
		time.Sleep(2 * time.Second)
		return w.OcApplyWorkspace(filePath)
	}
	if err != nil && !strings.Contains(output, "AlreadyExists") {
		fmt.Println(err)
	}
	return err
}

//launch 'exec' oc command in the defined pod and container
func ExecCommandInPod(podName string, containerName string, commandArgs ...string) (commandResult string) {
	commandsArg := append([]string{"exec", podName, "--namespace", "user-devvs", "-c", containerName}, commandArgs...)
	cmd := exec.Command("oc", commandsArg...)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if (err != nil) && (!strings.Contains(output, "denied the request: The only workspace creator has exec access")) {
		log.Fatal(err, "Cannot execute command in the dedicated container:", output)
	}
	return output
}

// login into a cluster using oc client ant return output of executed command
func LoginIntoClusterWithCredentials(loginName string, loginPass string, clusterConsoleUrl string) string {
	cmd := exec.Command("oc", "login", "-u", loginName, "-p", loginPass, clusterConsoleUrl)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		log.Fatal("Cannot login into the cluster with oc client", err)
	}
	return output
}

//create a project under login user using oc client
func CreateProjectWithOcClient(projectName string, description string, displayName string) string {
	cmd := exec.Command("oc", "new-project", projectName, "--description="+description, "--display-name="+displayName)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		log.Fatalf("Cannot create the project %s using oc client %s, %s", projectName, err, output)
	}
	return output
}

// login into a cluster using oc client and return output of executed command
func LoginIntoClusterWithToken(token string, clusterConsoleUrl string) string {
	cmd := exec.Command("oc", "login", "--token", token, "--server=", clusterConsoleUrl)
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		log.Fatal("Cannot login into the cluster with oc client", err)
	}
	return output
}
