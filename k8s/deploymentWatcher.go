package deployment

import (
	"fmt"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfig = "/home/clanderas/.kube/config"

func StartWatcher(changes chan<- *v1.Deployment ) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if(err != nil) {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}
	deploymentClient := clientset.AppsV1().Deployments("default")


	watcher, err := deploymentClient.Watch(metav1.ListOptions{Watch:true})
	if err != nil {
		panic(err)
	}

	for event := range watcher.ResultChan() {
		deploy, ok := event.Object.(*v1.Deployment)
		if ok {
			fmt.Println(deploy.Name, "has changed")
			fmt.Println(deploy.Name, "has available", deploy.Status.AvailableReplicas, "of", *deploy.Spec.Replicas)
			changes <- deploy
		}
	}
}