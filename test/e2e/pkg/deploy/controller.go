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

package deploy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
)

func (w *Deployment) CreateNamespace() error {
	_, err := w.kubeClient.Kube().CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: config.Namespace,
		},
	}, metav1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func (w *Deployment) DeployWorkspacesController() error {
	cmd := exec.Command("make", "install")
	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))

	if err != nil && !strings.Contains(string(output), "AlreadyExists") {
		fmt.Println(err)
		return err
	}
	return nil
}

func (w *Deployment) CustomResourceDefinitions() error {
	devWorkspaceCRD := exec.Command("make", "install_crds")
	output, err := devWorkspaceCRD.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "AlreadyExists") {
		fmt.Println(err)
		return err
	}
	return nil
}
