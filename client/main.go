package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k8s-api-demo/proto"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"strings"
	"time"
)

const (
	port = 50051
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	namespace := "default"

	for {
		ctx := context.Background()
		services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
		if errors.IsNotFound(err) {
			fmt.Printf("service not found in default namespace\n")
		} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
			fmt.Printf("Error getting service %v\n", statusError.ErrStatus.Message)
		} else if err != nil {
			panic(err.Error())
		}
		fmt.Printf("Found services in %s namespace\n", namespace)
		fmt.Printf("There are %d services in the cluster\n", len(services.Items))
		for _, service := range services.Items {
			name := service.ObjectMeta.Name
			ip := service.Spec.ClusterIP
			fmt.Printf("service: %v\n ip: %v\n\n", name, ip)
			if strings.Contains(name, "greeter") {
				sendReq(ctx, service.Spec.ClusterIP)
			}
		}
		time.Sleep(20 * time.Second)
	}
}

func sendReq(ctx context.Context, ip string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", ip, port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := "world"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Printf("Greeting: %s\n\n", r.GetMessage())
}
