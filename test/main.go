package main

import (
	"github.com/devfile/devworkspace-operator/test/e2e/pkg/client"
	v1 "k8s.io/api/core/v1"
	//"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	//restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"log"
	"k8s.io/client-go/tools/clientcmd"
	"io"
)

func main (){
k8sClient, err := client.NewK8sClient()
	cmd := []string{
		"sh",
		"-c",
		"echo Hello",
	}
if err != nil{
	log.Fatal("Cannot create k8s klient inst")
}
k8sClient.Kube().
}

func ExecCmdExample() error {
	config, _ := clientcmd.BuildConfigFromFlags("", "")
	cmd := []string{
		"sh",
		"-c",
		"echo Hello",
	}
	k8sClient, err := client.NewK8sClient()


	req := k8sClient.Kube().CoreV1().RESTClient().Post().Resource("pods").Name("").
		Namespace("default").SubResource("exec")

	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}

	return nil
}