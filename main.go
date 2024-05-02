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
	ns      = "default"
	podName = "nginx-test-84b487c7d9-dhkfs"
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
	pod, err := clientSet.CoreV1().Pods(ns).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("HostIP is: %s\n", pod.Status.HostIP)
	for index, ip := range pod.Status.PodIPs {
		fmt.Printf("index is: %d, ip is: %s\n", index, ip.IP)
	}
	fmt.Printf("pod ip is: %s\n", pod.Status.PodIP)
}
