package main

//
//Temporary file for quick launching
//

func main() {
	//config.Namespace = "devworkspace-controller"
	//k8sClient, err := client.NewK8sClient()
	//podList, err := k8sClient.Kube().CoreV1().Pods(config.Namespace ).List(context.TODO(),metav1.ListOptions{})
	//
	//if err != nil {
	//	log.Fatal("Error!!!")
	//}

	//for _, item := range  podList.Items {
	//	if strings.HasPrefix(item.Name,"workspace"){
	//		found:=item.Name
	//		found2:=item.Spec.Containers
	//		fmt.Println(found, found2)
	//	}
	//
	//}
	//podName := k8sClient.FindPodInNamespaceByNamePrefix("workspace")
	//fmt.Println(client.ExecCommandInPod(podName, "dev", "echo", "hello"))
}
