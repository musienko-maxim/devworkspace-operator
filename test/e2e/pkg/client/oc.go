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
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		"KUBECONFIG=%s oc apply --namespace %s -f %s",
		w.kubeCfgFile,
		config.DevNameSpace,
		filePath))
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
func (w *K8sClient) ExecCommandInContainerAsRegularUser(podName string, commandInContainer string) (commandResult string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		"KUBECONFIG=%s oc exec %s -n %s -c dev %s",
		w.kubeCfgFile,
		podName,
		config.DevNameSpace,
		commandInContainer))

	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if (err != nil) && (!strings.Contains(output, "denied the request: The only workspace creator has exec access")) {
		log.Fatal(err, "Cannot execute command in the dedicated container:", output)
	}
	return output
}

//launch 'exec' oc command in the defined pod and container
func (w *K8sClient) ExecCommandInPodAsDefaultUser(podName, commandInContainer string) (commandResult string, err error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		"oc exec %s -n %s -c dev %s",
		podName,
		config.DevNameSpace,
		commandInContainer))
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil && strings.Contains(output, "The only workspace creator has exec access") {
		return output, nil
	} else {
		return output, err
	}

}

//create a project under login user using oc client
func (w *K8sClient) CreateProjectWithKubernetesContext(projectName, description, displayName string) (error, string)  {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(
		"KUBECONFIG=%s oc new-project %s --description=%s --display-name=%s",
		w.kubeCfgFile,
		projectName,
		description,
		displayName))
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		return err, output
		//log.Fatalf("Cannot create the project %s using oc client %s, %s", projectName, err, output)
	}
	return err, output
}
