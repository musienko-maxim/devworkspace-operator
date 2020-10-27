package test

import (
	"context"
	"os"

	devworkspace "github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	_ "github.com/devfile/devworkspace-operator/test/e2e/pkg/tests"
	"k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("test/cmd")

func main() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Failed to create client config")
		os.Exit(1)
	}

	client, err := crclient.New(cfg, crclient.Options{})

	if err != nil {
		log.Error(err, "Failed to create client")
		os.Exit(1)
	}

	namespacedName := types.NamespacedName{
		Name:      "name",
		Namespace: "namespace",
	}

	workspace := &devworkspace.DevWorkspace{}
	err = client.Get(context.TODO(), namespacedName, workspace)

	if err != nil {
		panic(err)
	}

	// Here we go. We have workspace fetched
}
