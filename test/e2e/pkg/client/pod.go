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
	"context"
	"errors"
	"fmt"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
	_ "github.com/devfile/devworkspace-operator/test/e2e/pkg/tests"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	configs "sigs.k8s.io/controller-runtime/pkg/client/config"
	"time"
)

var (
	Scheme             = runtime.NewScheme()
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: v1alpha1.SchemeGroupVersion.Group, Version: v1alpha1.SchemeGroupVersion.Version}
)

func (w *K8sClient) WaitForPodRunningByLabel(label string) (deployed bool, err error) {
	timeout := time.After(6 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			err := w.WaitForRunningPodBySelector(config.Namespace, label, 3*time.Minute)
			if err == nil {
				return true, nil
			}
		}
	}
}

// Wait up to timeout seconds for all pods in 'namespace' with given 'selector' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func (w *K8sClient) WaitForRunningPodBySelector(namespace, selector string, timeout time.Duration) error {
	podList, err := w.ListPods(namespace, selector)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		fmt.Println("Pod not created yet with selector " + selector + " in namespace " + namespace)

		return fmt.Errorf("Pod not created yet in %s with label %s", namespace, selector)
	}

	for _, pod := range podList.Items {
		fmt.Println("Pod " + pod.Name + " created in namespace " + namespace + "...Checking startup data.")
		if err := w.waitForPodRunning(namespace, pod.Name, timeout); err != nil {
			return err
		}
	}

	return nil
}

//get workspace current dev workspace status from the Custom Resource object
func (w *K8sClient) GetDevWsStatus () (status v1alpha1.WorkspacePhase) {
	if err := AddToScheme(scheme.Scheme); err != nil {
		logrus.Fatalf("Failed to add CRD to scheme")
	}
	if err := api.AddToScheme(Scheme); err != nil {
		logrus.Fatalf("Failed to add CRD to scheme")
	}

	cfg, err := configs.GetConfig()
	if err != nil {
		logrus.Error(err, "Failed to create client config")
		os.Exit(1)
	}

	client, err := crclient.New(cfg, crclient.Options{})

	if err != nil {
		logrus.Error(err, "Failed to create client")
		os.Exit(1)
	}

	namespacedName := types.NamespacedName{
		Name:      "web-terminal",
		Namespace: "devworkspace-controller",
	}

	workspace := &v1alpha1.DevWorkspace{}
	err = client.Get(context.TODO(), namespacedName, workspace)

	logrus.Info("Workspace status is: " + workspace.Status.Phase)

	if err != nil {
		panic(err)
	}
	return  workspace.Status.Phase
}

func (w *K8sClient) WaitDevWsStatus(expectedStatus v1alpha1.WorkspacePhase) (bool, error) {
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			currentStatus := w.GetDevWsStatus()
			if currentStatus == expectedStatus {
				return true, nil
			}
		}
	}
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1alpha1.DevWorkspace{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// Returns the list of currently scheduled or running pods in `namespace` with the given selector
func (w *K8sClient) ListPods(namespace, selector string) (*v1.PodList, error) {
	listOptions := metav1.ListOptions{LabelSelector: selector}
	podList, err := w.Kube().CoreV1().Pods(namespace).List(listOptions)

	if err != nil {
		return nil, err
	}
	return podList, nil
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func (w *K8sClient) waitForPodRunning(namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, w.isPodRunning(podName, namespace))
}

// return a condition function that indicates whether the given pod is
// currently running
func (w *K8sClient) isPodRunning(podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, _ := w.Kube().CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		age := time.Since(pod.GetCreationTimestamp().Time).Seconds()

		switch pod.Status.Phase {
		case v1.PodRunning:
			fmt.Println("Pod started after", age, "seconds")
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, nil
		}
		return false, nil
	}
}
