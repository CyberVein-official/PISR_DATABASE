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

func (r *RpcClient) QueryPrivateDataWithAddress(cmd *proto.QueryPrivateWithAddrRequest) (*proto.QueryResponse, error) {
	response, err := r.app.QueryPrivateDataWithAddress(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.QueryPrivateDataWithAddress(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}

func (r *RpcClient) Execute(cmd *proto.CommandRequest) (*proto.ExecuteResponse, error) {
	response, err := r.app.Execute(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.Execute(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}

func (r *RpcClient) ExecuteAsync(cmd *proto.CommandRequest) (*proto.ExecuteAsyncResponse, error) {
	response, err := r.app.ExecuteAsync(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.ExecuteAsync(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}

func (r *RpcClient) ExecuteWithPrivateKey(cmd *proto.CommandRequest) (*proto.ExecuteResponse, error) {
	response, err := r.app.ExecuteWithPrivateKey(r.GetContextWithToken(), cmd)
	if err != nil {
		if r.CheckExpire(err) {
			r.UpdateAuth()
			return r.app.ExecuteWithPrivateKey(r.GetContextWithToken(), cmd)
		}
		return nil, err
	}
	return response, nil
}

func (r *RpcClient) GetContextWithToken() context.Context {
	md := metadata.Pairs("cybervein_token", r.token)
	return metadata.NewOutgoingContext(context.Background(), md)
}

func (r *RpcClient) CheckExpire(err error) bool {
	if strings.Contains(err.Error(), code.CodeTypeTokenTimeoutErrorMsg) {
		return true
	}
	return false
}

func (r *RpcClient) UpdateAuth() error {
	token, err := r.app.Auth(context.Background(), &proto.AuthRequest{Password: r.password})
	if err != nil {
		return err
	}
	r.token = token.Token
	return nil
}
