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
)

func main()	{
	//devK8sClient, err := client.NewK8sClientWithContext("developer", "developer", "/tmp/admin123-kubeconfig")
	client.LoginIntoClusterWithCredentials("","","")
	devK8sClient, err := client.NewK8sClientWithContext("")
	//adminK8sClient, err := client.NewK8sClientWithConfigFile("~/.kube/config")

	if err != nil {
		fmt.Println(err)
	}
	devK8sClient.FindPodInNamespaceByNamePrefix("")




}
