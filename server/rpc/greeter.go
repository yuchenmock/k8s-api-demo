package rpc

import (
	"context"
	"k8s-api-demo/proto"
)

type RpcServer struct{
	proto.UnimplementedGreeterServer
}

func (r *RpcServer) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	return &proto.HelloReply{
		Message: "hello " + request.GetName(),
	}, nil
}
