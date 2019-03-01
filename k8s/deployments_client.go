package deployment

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strconv"
)

type deploymentClient struct {
	deployments v12.DeploymentInterface
	kubeconfig  string
}


func NewClient(config ClusterConfig) *deploymentClient {

		tlsConfig := rest.TLSClientConfig{
			Insecure: !config.VerifyTlsCertificate,
		}
		k8sConfig := rest.Config{
			Host: config.Host,
			TLSClientConfig: tlsConfig,
			BearerToken: config.BearerToken,
		}

		clientset, err := kubernetes.NewForConfig(&k8sConfig)

		if err != nil {
			panic(err)
		}

		return &deploymentClient{
			deployments: clientset.AppsV1().Deployments("default"),
		}
}

func NewClientFromKubeConfig(kubeconfig string) *deploymentClient {

	if _, err := os.Stat(kubeconfig); os.IsNotExist(err) {
		panic(fmt.Sprintf("File %s was not found", kubeconfig))
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
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

func NewFromConfig() (*deploymentClient, error) {

	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Could not read config.ini")
	}

	kubeconfig := cfg.Section("").Key("kubeconfig").Value()
	if kubeconfig != "" {
		fmt.Println("Creating client from kubeconfig:", kubeconfig)
		return NewClientFromKubeConfig(kubeconfig), nil
	}

	k8sHost := cfg.Section("").Key("k8shost").Value()

	if k8sHost != "" {
		fmt.Println("Creating client from configured k8s host:", k8sHost)

		verifytls, err := strconv.ParseBool(cfg.Section("").Key("verifytls").Value())
		if err != nil {
			return nil, errors.New("boolean value verifytls was not found or has incorrect format")
		}

		bearertoken := cfg.Section("").Key("bearertoken").Value()

		if bearertoken == "" {
			return nil, errors.New("bearer token must be configured to connect to the remote cluster")
		}

		return NewClient(ClusterConfig{
			Host: k8sHost,
			VerifyTlsCertificate: verifytls,
			BearerToken: bearertoken,
		}), nil
	}

	return nil, errors.New("No valid configuration found. Use kubeconfig or k8shost section to configure the cluster")

}

func (d *deploymentClient) GetDeployments() *v1.DeploymentList {
	deployments, err := d.deployments.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	return deployments
}
func (d *deploymentClient) StartWatcher(changes chan<- *v1.Deployment) {

	fmt.Println("Start deployments watcher")
	watcher, err := d.deployments.Watch(metav1.ListOptions{Watch: true})
	if err != nil {
		panic(err)
	}

	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		deploy, ok := event.Object.(*v1.Deployment)
		if ok {
			fmt.Println(deploy.Name, "has changed")
			fmt.Println(deploy.Name, "has available", deploy.Status.AvailableReplicas, "of", *deploy.Spec.Replicas)
			changes <- deploy
		}
	}

}
