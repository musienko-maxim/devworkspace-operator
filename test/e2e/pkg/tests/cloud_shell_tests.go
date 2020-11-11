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

var _ = ginkgo.Describe("[Create Openshift Web Terminal Workspace]", func() {
	k8sClient, err := client.NewK8sClient()

	ginkgo.It("Wait devworkspace controller Pod", func(){
	controllerLabel:= " app.kubernetes.io/name=devworkspace-controller"
		if err != nil {
			ginkgo.Fail("Failed to create k8s client: " + err.Error())
			return
		}
		deploy, err := k8sClient.WaitForPodRunningByLabel(controllerLabel)
		if err != nil {
			ginkgo.Fail(fmt.Sprint("Cannot get the Pod status with label %s: \n" , controllerLabel) +  err.Error())
			return
		}
		if !deploy {
			fmt.Println("Devworkspace controller  didn't start properly")
		}
	})


	ginkgo.It("Wait webhook controller Pod", func(){
		controllerLabel:= "app.kubernetes.io/name=devworkspace-webhook-server"
		if err != nil {
			ginkgo.Fail("Failed to create k8s client: " + err.Error())
			return
		}
		deploy, err := k8sClient.WaitForPodRunningByLabel(controllerLabel)
		if err != nil {
			ginkgo.Fail(fmt.Sprint("Cannot get the Pod status with label %s: \n" , controllerLabel) +  err.Error())
			return
		}
		if !deploy {
			fmt.Println("Devworkspace controller  didn't start properly")
		}
	})

	ginkgo.It("Add openshift web terminal to cluster", func() {
		label := "controller.devfile.io/workspace_name=web-terminal"

		if err != nil {
			ginkgo.Fail("Failed to create k8s client: " + err.Error())
			return
		}
		err = k8sClient.OcApplyWorkspace("samples/web-terminal.yaml")
		if err != nil {
			ginkgo.Fail("Failed to create openshift web terminal workspace: " + err.Error())
			return
		}
		deploy, err := k8sClient.WaitForPodRunningByLabel(label)

		if !deploy {
			fmt.Println("Openshift Web terminal workspace didn't start properly")
		}

		gomega.Expect(err).NotTo(gomega.HaveOccurred())
	})

	ginkgo.It("Wait the running status of workspace", func() {
				client.WaitDevWsStatus(v1alpha1.WorkspaceStatusRunning)
	})

})
