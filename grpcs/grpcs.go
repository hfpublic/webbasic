package grpcs

import "google.golang.org/grpc"

type addRegisterFunc func(*grpc.Server)

var grpcRegisterFuncs = make([]addRegisterFunc, 0)

func AddGrpcRegisterFuncs(registerFunc addRegisterFunc) {
	grpcRegisterFuncs = append(grpcRegisterFuncs, registerFunc)
}

func RegisterFuncs(serv *grpc.Server) {
	for _, registerFunc := range grpcRegisterFuncs {
		registerFunc(serv)
	}
}
