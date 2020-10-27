package test

import (
	"context"
	"fmt"
	"github.com/devfile/api/tree/master/pkg/apis/workspaces/v1alpha1"
	_ "github.com/devfile/devworkspace-operator/test/e2e/pkg/tests"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	workspace := v1alpha1.{}

	fmt.Println(workspace)

	// c is a created client.Client
	cl:= client.Client.Get(context.TODO(), client.ObjectKey{
		Namespace: "namespace",
		Name:      "name"}, workspace)


}
