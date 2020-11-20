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
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("[Create OpenShift Web Terminal Workspace]", func() {
	//devK8sClient, err := client.NewK8sClientWithCredentials("developer", "developer")
	//adminK8sClient, err := client.NewK8sClientWithConfigFile("~/.kube/config")
	k8sClient, err := client.NewK8sClient()

	ginkgo.It("Wait devworkspace controller Pod", func() {
		controllerLabel := "app.kubernetes.io/name=devworkspace-controller"
		if err != nil {
			ginkgo.Fail("Failed to create k8s client: " + err.Error())
			return
		}
		//deploy, err := adminK8sClient.WaitForPodRunningByLabel(controllerLabel)
		deploy, err := k8sClient.WaitForPodRunningByLabel(controllerLabel)
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("cannot get the Pod status with label %s: %s", controllerLabel, err.Error()))
			return
		}
		if !deploy {
			fmt.Println("DevWorkspace controller  didn't start properly")
		}
	})

	ginkgo.It("Wait webhook controller Pod", func() {
		controllerLabel := "app.kubernetes.io/name=devworkspace-webhook-server"
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
			fmt.Println("Devworkspace webhook  didn't start properly")
		}
	})

	ginkgo.It("Add OpenShift web terminal to cluster", func() {
		client.LoginIntoClusterWithCredentials("developer", "developer", "https://api.crc.testing:6443")

		//devK8sClient.CreateProject("web-terminal", "webterminal-test-project", "web-terminal")
		client.CreateProjectWithOcClient("web-terminal", "webterminal-test-project", "web-terminal")

		userK8sClient, err := client.NewK8sClient()

		if err != nil {
			ginkgo.Fail("Failed to create userK8sClient client: ")
			return
		}

		//devK8sClient.OcApplyWorkspace("samples/web-terminal.yaml")
		err = userK8sClient.OcApplyWorkspace("samples/web-terminal.yaml")

		if err != nil {
			ginkgo.Fail("Failed to create OpenShift web terminal workspace: " + err.Error())
			return
		}

		//TODO Apart from waiting to Running state it makes sense to early fail when status changed to failed.
		deploy, err := client.WaitDevWsStatus(v1alpha1.WorkspaceStatusRunning)

		if !deploy {
			fmt.Println("OpenShift Web terminal workspace didn't start properly")
		}

		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

})
