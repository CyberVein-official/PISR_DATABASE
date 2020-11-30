package client

import (
	"context"

	proto "cybervein.org/CyberveinDB/grpc/proto"
	"google.golang.org/grpc"
)

type RpcClient struct {
	app      proto.cyberveinClient
	token    string
	password string
}

func NewRpcClient(address string, password string) (*RpcClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := &RpcClient{app: proto.NewcyberveinClient(conn), token: "", password: password}
	token, err := client.app.Auth(context.Background(), &proto.AuthRequest{Password: password})
	if err != nil {
		return nil, err
	}
	client.token = token.Token
	return client, nil
}

func (r *RpcClient) Query(cmd *proto.CommandRequest) (*proto.QueryResponse, error) {
	response, err := r.app.Query(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.Query(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}

func (r *RpcClient) QueryPrivateData(cmd *proto.CommandRequest) (*proto.QueryResponse, error) {
	response, err := r.app.QueryPrivateData(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.QueryPrivateData(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}
