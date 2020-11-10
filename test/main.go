package main

import (
	"context"
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	_ "github.com/devfile/devworkspace-operator/test/e2e/pkg/tests"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	Scheme             = runtime.NewScheme()
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemeGroupVersion = schema.GroupVersion{Group: v1alpha1.SchemeGroupVersion.Group, Version: v1alpha1.SchemeGroupVersion.Version}
)

func main() {
	if err := AddToScheme(scheme.Scheme); err != nil {
		logrus.Fatalf("Failed to add CRD to scheme")
	}
	if err := api.AddToScheme(Scheme); err != nil {
		logrus.Fatalf("Failed to add CRD to scheme")
	}

	cfg, err := config.GetConfig()
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
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&v1alpha1.DevWorkspace{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
