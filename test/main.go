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

package main

import (
	"fmt"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/client"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
	"log"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	config.DevNameSpace = "user-namespace"
	devCubeConfig := "/tmp/devconfig"
	k8sClient, err := client.NewK8sClient()
	devK8sClient, err := client.NewK8sClientWithCredentials("developer", "developer", devCubeConfig, "https://api.crc.testing:6443")
	//config.DevNameSpace="devns"
	//nsSpec:=&v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "webterminal"}}
	//devK8sClient.Kube().CoreV1().Namespaces().Create(context.TODO(), nsSpec, metav1.CreateOptions{})
	if err != nil {
		log.Fatal("Failed to create userK8sClient client: ")

	}
	//}
	//deploy, err :=	client.WaitDevWsStatus(v1alpha1.WorkspaceStatusRunning)
	//if !deploy {
	//,err:=k8sClient.ListPods(config.DevNameSpace, "controller.devfile.io/workspace_name=web-terminal")
	//fmt.Println(list.Items[0].Name)
	podName:=k8sClient.GetPodNameFromUserNameSpaceByLabel("controller.devfile.io/workspace_name=web-terminal")
	result:=devK8sClient.ExecCommandInContainerAsRegularUser(podName, "echo hello dev")
	fmt.Println("<<<<<<<<<<<<"+result)
}
