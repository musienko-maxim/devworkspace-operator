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

package cmd

import (
	"fmt"
	"context"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/client"

	"path/filepath"
	"testing"

	workspaceWebhook "github.com/devfile/devworkspace-operator/webhook/workspace"

	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/deploy"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "github.com/devfile/devworkspace-operator/test/e2e/pkg/tests"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

//Create Constant file
const (
	testResultsDirectory = "/tmp/artifacts"
	jUnitOutputFilename  = "junit-workspaces-operator.xml"
)

//SynchronizedBeforeSuite blocks is executed before run all test suites
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	fmt.Println("Starting to setup objects before run ginkgo suite")
	config.Namespace = "devworkspace-controller"

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		fmt.Println("Failed to create workspace client")
		panic(err)
	}

	controller := deploy.NewDeployment(k8sClient)
	err = controller.DeployWorkspacesController()

	if err != nil {
		fmt.Println("Cannot deploy Web Operator using Make file")
		panic(err)
	}

	return nil
}, func(data []byte) {})

var _ = ginkgo.SynchronizedAfterSuite(func() {
	k8sClient, err := client.NewK8sClient()

	if err != nil {
		_ = fmt.Errorf("Failed to uninstall workspace controller %s", err)
	}

	if err = k8sClient.Kube().AdmissionregistrationV1().MutatingWebhookConfigurations().Delete(context.TODO(),workspaceWebhook.MutateWebhookCfgName, metav1.DeleteOptions{}); err != nil {
		_ = fmt.Errorf("Failed to delete mutating webhook configuration %s", err)
	}

	if err = k8sClient.Kube().AdmissionregistrationV1().ValidatingWebhookConfigurations().Delete(context.TODO(), workspaceWebhook.ValidateWebhookCfgName,  metav1.DeleteOptions{}); err != nil {
		_ = fmt.Errorf("Failed to delete validating webhook configuration %s", err)
	}
}, func() {})

func TestWorkspaceController(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	fmt.Println("Creating ginkgo reporter for Test Harness: Junit and Debug Detail reporter")
	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(testResultsDirectory, jUnitOutputFilename)))

	fmt.Println("Running Workspace Controller e2e tests...")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Workspaces Controller Operator Tests", r)
}
