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

package tests

import (
	"fmt"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/client"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[Create OpenShift Web Terminal Workspace]", func() {
	var podName string
	devCubeConfig := "/tmp/devconfig"
	k8sClient, err := client.NewK8sClient()
	devK8sClient, err := client.NewK8sClientWithCredentials("developer", "developer", devCubeConfig, "https://api.crc.testing:6443")

	ginkgo.It("Wait devworkspace controller Pod", func() {
		controllerLabel := "app.kubernetes.io/name=devworkspace-controller"

		if err != nil {
			ginkgo.Fail("Failed to create k8s client: " + err.Error())
			return
		}

		deploy, err := k8sClient.WaitForPodRunningByLabel(controllerLabel)

		if err != nil {
			ginkgo.Fail(fmt.Sprintf("cannot get the Pod status with label %s: %s", controllerLabel, err.Error()))
			return
		}

		if !deploy {
			fmt.Println("DevWorkspace controller  didn't start properlyб")
		}
	})

	ginkgo.It("Wait webhook controller Pod", func() {
		controllerLabel := "app.kubernetes.io/name=devworkspace-webhook-server"

		deploy, err := k8sClient.WaitForPodRunningByLabel(controllerLabel)

		if err != nil {
			ginkgo.Fail(fmt.Sprintf("cannot get the Pod status with label %s: %s", controllerLabel, err.Error()))
			return
		}

		if !deploy {
			ginkgo.Fail(fmt.Sprintf("Devworkspace webhook  didn't start properly %s"))
		}
	})

	ginkgo.It("Add OpenShift web terminal to cluster and wait running status", func() {

		devK8sClient.CreateProjectWithKubernetesContext(config.DevNameSpace, "max", config.DevNameSpace)

		err = devK8sClient.OcApplyWorkspace("samples/web-terminal.yaml")

		if err != nil {
			ginkgo.Fail("Failed to create OpenShift web terminal workspace: " + err.Error())
			return
		}

		deploy, err := client.WaitDevWsStatus(v1alpha1.WorkspaceStatusRunning)

		if !deploy {
			ginkgo.Fail(fmt.Sprintf("OpenShift Web terminal workspace didn't start properly. Error: %s", err))
		}

	})

	ginkgo.It("Check that pod creator can execute a command in the container", func() {
		podSelector := "controller.devfile.io/workspace_name=web-terminal"
		podName = k8sClient.GetPodNameFromUserNameSpaceByLabel(podSelector)
		resultOfExecCommand := devK8sClient.ExecCommandInContainerAsRegularUser(podName, "echo hello dev")

		gomega.Expect(resultOfExecCommand).To(gomega.ContainSubstring("hello dev"))
	})

	ginkgo.It("Check that not pod owner cannot execute a command in the container", func() {
		expectedMessageSuffix := "denied the request: The only workspace creator has exec access"
		resultOfExecCommand,_ := k8sClient.ExecCommandInPodAsDefaultUser(podName, "echo hello dev")

		gomega.Expect(resultOfExecCommand).To(gomega.ContainSubstring(expectedMessageSuffix))
	})

})
