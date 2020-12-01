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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/exec"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

)

type K8sClient struct {
	kubeClient  *kubernetes.Clientset
	kubeCfgFile string // generate when client is created and store config there
}

// NewK8sClient creates kubernetes client wrapper with helper functions and direct access to k8s go client
func NewK8sClient() (*K8sClient, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	h := &K8sClient{kubeClient: client}
	return h, nil
}
// create kubernetes instance using a kubernetes context
func NewK8sClientWithContext(pathToContextFile string) (*K8sClient, error) {
	//generate kubeconfig file name
	cfg, err := clientcmd.BuildConfigFromFlags("kubeconfig", pathToContextFile)
	//cfg, err := config.GetConfigWithContext(kubeCfgFile)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	h := &K8sClient{kubeClient: client, kubeCfgFile: pathToContextFile}
	return h, nil
}
// generate kubernetes config file using a user cluster credentials and create kubernetes instance
// with kube context under the user
func NewK8sClientWithCredentials(login, password, pathToCfgFile, clusterConsoleUrl string) (*K8sClient, error) {
	cmd := exec.Command("bash",
		"-c", fmt.Sprintf(
			"KUBECONFIG=%s"+
		" oc login -u %s -p %s --insecure-skip-tls-verify=true %s",
		pathToCfgFile, login, password, clusterConsoleUrl))
	outBytes, err := cmd.CombinedOutput()
	output := string(outBytes)
	if err != nil {
		log.Fatal(fmt.Sprintf("Cannot login into the cluster with oc client %s %s", err, output))
	}
	cfg, err := clientcmd.BuildConfigFromFlags("kubeconfig", pathToCfgFile)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	h := &K8sClient{kubeClient: client, kubeCfgFile: pathToCfgFile}
	return h, nil
}

// Kube returns the clientset for Kubernetes upstream.
func (c *K8sClient) Kube() kubernetes.Interface {
	return c.kubeClient
}
