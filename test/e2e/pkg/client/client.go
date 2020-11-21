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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func NewK8sClientWithContext(contextFile string) (*K8sClient, error) {
	//generate kubeconfig file name
	kubeCfgFile := "/tmp/admin123-kubeconfig"

	//copy contextFile to kubeCfgFile
	cfg,err := clientcmd.BuildConfigFromFlags("kubeconfig", kubeCfgFile)
	//cfg, err := config.GetConfigWithContext(kubeCfgFile)
	if err != nil {
		return nil, err
	}
	client, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	h := &K8sClient{kubeClient: client, kubeCfgFile: kubeCfgFile}
	return h, nil
}

func NewK8sClientWithCredentials(login, password, pathToCfgFile string) (*K8sClient, error) {
	//generate kubeconfig file name

	//LoginIntoClusterWithCredentials(login,password, "https://api.crc.testing:6443")
	//execute: KUBECONFIG=/tmp/admin123-kubeconfig oc login -u $login -p $password --insecure-skip-tls-verify=true https://api.crc.testing:6443

	cfg, err := config.GetConfigWithContext(pathToCfgFile)
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
