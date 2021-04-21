package datas

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
)

type GRPCConfig struct {
	Host string
	Port int
}

var grpcClient *GRPCClient

type GRPCClient struct {
	*grpc.ClientConn
}

func InitGRPCClient(conf *GRPCConfig) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%v", conf.Host, conf.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalln("dial error:", err)
	}
	grpcClient = &GRPCClient{ClientConn: conn}
}

func GetGRPCClient() *GRPCClient {
	return grpcClient
}

func CloseGRPCClient() error {
	return grpcClient.Close()
}
