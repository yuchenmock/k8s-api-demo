## Demo for service discover with k8s API

#### try the project
install minikube locally
execute setup.sh

#### k8s API

core logic
```go
import (
    "k8s.io/apimachinery/pkg/api/errors"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
)
config, err := rest.InClusterConfig()
clientset, err := kubernetes.NewForConfig(config)
services, err := clientset.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
for _, service := range services.Items {
    name := service.ObjectMeta.Name
    ip := service.Spec.ClusterIP
}
```