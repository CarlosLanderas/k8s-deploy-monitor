package deployment

import (
	"fmt"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/tools/clientcmd"
)

const kubeconfig = "/home/clanderas/.kube/config"

type deploymentClient struct {
	deployments v12.DeploymentInterface
}

func NewClient() *deploymentClient {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if(err != nil) {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err)
	}
	return &deploymentClient{
	 deployments: clientset.AppsV1().Deployments("default"),
	}
}
func (d *deploymentClient) GetDeployments() *v1.DeploymentList {
	deployments, err := d.deployments.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	return deployments
}
func (d *deploymentClient) StartWatcher(changes chan<- *v1.Deployment ) {

	watcher, err := d.deployments.Watch(metav1.ListOptions{Watch:true})
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