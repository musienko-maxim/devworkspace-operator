package client

import (
	"context"
	"errors"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
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

//get workspace current dev workspace status from the Custom Resource object
func  GetDevWsStatus() (status v1alpha1.WorkspacePhase) {
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

	if err != nil {
		panic(err)
	}
	return workspace.Status.Phase
}

func WaitDevWsStatus(expectedStatus v1alpha1.WorkspacePhase) (bool, error) {
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			logrus.Info("Now current status of developer workspace is: " + GetDevWsStatus())
			currentStatus := GetDevWsStatus()
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

