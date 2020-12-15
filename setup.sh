#!/bin/bash
export GOOS=linux
eval "$(minikkube docker-env)"
pushd server || exit
go build -o app .
docker build -t greeter-service:1.0 .
popd || exit
pushd client || exit
go build -o app .
docker build -t greeter-client:1.0 .

kubectl create clusterrolebinding default-view --clusterrole=view --serviceaccount=default:default

kubectl create deployment greeter-demo --image=greeter-service:1.0
kubectl expose deployment greeter-demo --type=NodePort --port=50051
kubectl create deployment greeter-demo2 --image=greeter-service:1.0
kubectl expose deployment greeter-demo2 --type=NodePort --port=50051
# this could be executed in another cmd manually to verify the service discovery
# kubectl create deployment greeter-demo3 --image=greeter-service:1.0
# kubectl expose deployment greeter-demo3 --type=NodePort --port=50051

kubectl run --generator=run-pod/v1 -i --rm client-demo --image=greeter-client:1.0

kubectl delete service greeter-demo
kubectl delete service greeter-demo2
kubectl delete deployment greeter-demo
kubectl delete deployment greeter-demo2
kubectl delete pod client-demo