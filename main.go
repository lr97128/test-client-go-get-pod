package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	ns = "default"
)

func main() {
	ctx := context.TODO()
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		inClusterConfig, err := rest.InClusterConfig()
		if err != nil {
			log.Panic("Can't get Config")
		}
		config = inClusterConfig
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic(err)
	}
	pods, err := clientSet.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Panic(err)
	}
	for _, pod := range pods.Items {

		fmt.Printf("HostIP is: %s\n", pod.Status.HostIP)
		for _, initStatuts := range pod.Status.InitContainerStatuses {
			fmt.Printf("InitContainerStatus is: %s\n", initStatuts.State.Terminated.Reason)
		}
		for index, ip := range pod.Status.PodIPs {
			fmt.Printf("index is: %d, ip is: %s\n", index, ip.IP)
		}
		fmt.Printf("pod ip is: %s\n", pod.Status.PodIP)
		var rsName string
		for index, ownerReference := range pod.ObjectMeta.GetOwnerReferences() {
			rsName = ownerReference.Name
			fmt.Printf("pod's index: %d, ownerReferenceKind: %s, ownerReferenceName: %s\n", index, ownerReference.Kind, rsName)
		}
		rs, err := clientSet.AppsV1().ReplicaSets(ns).Get(ctx, rsName, metav1.GetOptions{})
		if err != nil {
			log.Panic(err)
		}
		for index, ownerReference := range rs.ObjectMeta.GetOwnerReferences() {
			fmt.Printf("ReplicaSet's index: %d, ownerReferenceKind: %s, ownerReferenceName: %s\n", index, ownerReference.Kind, ownerReference.Name)
		}
	}
}
