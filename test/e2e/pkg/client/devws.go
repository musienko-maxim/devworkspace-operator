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
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/config"
	"os"
	"time"

	"github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	configs "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Scheme             = runtime.NewScheme()
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: v1alpha1.SchemeGroupVersion.Group, Version: v1alpha1.SchemeGroupVersion.Version}
)

//get workspace current dev workspace status from the Custom Resource object
func GetDevWsStatus() (*v1alpha1.WorkspacePhase, error) {
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
		Namespace: config.DevNameSpace,
	}

	workspace := &v1alpha1.DevWorkspace{}
	err = client.Get(context.TODO(), namespacedName, workspace)

	if err != nil {
		return nil, err
	}
	return &workspace.Status.Phase, nil
}

func WaitDevWsStatus(expectedStatus v1alpha1.WorkspacePhase) (bool, error) {
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			currentStatus, err := GetDevWsStatus()
			logrus.Info("Now current status of developer workspace is: " + *currentStatus)
			if err != nil  {
				return false, err
			}
			if *currentStatus == v1alpha1.WorkspaceStatusFailed{
				return false, errors.New("workspace has been failed unexpectedly")
			}
			if *currentStatus == expectedStatus {
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
