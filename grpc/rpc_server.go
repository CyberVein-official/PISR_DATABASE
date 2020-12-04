package grpc

import (
	context "context"
	"cybervein.org/CyberveinDB/core"
	proto "cybervein.org/CyberveinDB/grpc/proto"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/models/code"
	"cybervein.org/CyberveinDB/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/reflection"
	"net"
	"strings"
)



type Server struct {
	port   string
	server *grpc.Server
	app    *cyberveinService
}



func NewRpcServer(port string) *Server {
	s := &Server{
		server: grpc.NewServer(),
		app:    &cyberveinService{},
		port:   port,
	}
	proto.RegistercyberveinServer(s.server, s.app)
	reflection.Register(s.server)
	return s
}


func (s *Server) StartServer() {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		logger.Log.Error("failed to listen: %v", err)
		return
	}
	if err := s.server.Serve(lis); err != nil {
		logger.Log.Error("failed to serve: %v", err)
		return
	}
}

type cyberveinService struct {
}

func (r *cyberveinService) CheckToken(c context.Context) error {
	respCode := code.CodeTypeOK
	respMsg := code.CodeTypeOKMsg
	md, ok := metadata.FromIncomingContext(c)
	if !ok {
		return fmt.Errorf("Error :rpc context parse error")
	}

	values := md.Get("cybervein_token")
	if values == nil || len(values) == 0 {
		respCode, respMsg = code.CodeTypeTokenInvalidError, code.CodeTypeTokenInvalidErrorMsg+" : empty"
	}
	token := values[0]

	p, _ := peer.FromContext(c)
	addr := p.Addr.String()
	claims, err := utils.ParseToken(token)
	if err != nil {
		switch err.(*jwt.ValidationError).Errors {
		case jwt.ValidationErrorExpired:
			respCode, respMsg = code.CodeTypeTokenTimeoutError, code.CodeTypeTokenTimeoutErrorMsg
		default:
			respCode, respMsg = code.CodeTypeTokenInvalidError, code.CodeTypeTokenInvalidErrorMsg+" : "+token
		}
	} else if strings.EqualFold(claims.Address, addr) {
		respCode, respMsg = code.CodeTypeTokenInvalidError, code.CodeTypeTokenInvalidErrorMsg+" : ip address"
	}

	if respCode != code.CodeTypeOK {
		return fmt.Errorf("Error %d : %s", respCode, respMsg)
	}
	return nil
}

func (r *cyberveinService) Auth(c context.Context, req *proto.AuthRequest) (*proto.Token, error) {
	p, _ := peer.FromContext(c)
	if req.Password != utils.Config.App.DbPassword {
		logger.Log.Error(fmt.Sprintf("Authorization Error from %s , password : %s", p, req.Password))
		return nil, fmt.Errorf("Authorization Error %d : %s ", code.CodeTypeDBPasswordIncorrectError, code.CodeTypeDBPasswordIncorrectErrorMsg)
	}
	s, err := utils.GenerateToken(p.Addr.String(), "", req.Password)
	if err != nil {
		return nil, fmt.Errorf("Authorization Error %d : %s ", code.CodeTypeInternalError, code.CodeTypeInternalErrorMsg)
	}
	return &proto.Token{Token: s}, nil
}

func (r *cyberveinService) Query(c context.Context, req *proto.CommandRequest) (*proto.QueryResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.Query(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.QueryResponse{
		Result: response.Result,
	}, nil
}

func (r *cyberveinService) QueryPrivateData(c context.Context, req *proto.CommandRequest) (*proto.QueryResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.QueryPrivateData(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.QueryResponse{
		Result: response.Result,
	}, nil
}

func (r *cyberveinService) QueryPrivateDataWithAddress(c context.Context, req *proto.QueryPrivateWithAddrRequest) (*proto.QueryResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.QueryPrivateDataWithAddress(&models.QueryPrivateWithAddrRequest{req.Cmd, req.Address})
	if err != nil {
		return nil, err
	}
	return &proto.QueryResponse{
		Result: response.Result,
	}, nil
}

func (r *cyberveinService) Execute(c context.Context, req *proto.CommandRequest) (*proto.ExecuteResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.Execute(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.ExecuteResponse{
		Cmd:           response.Cmd,
		ExecuteResult: response.ExecuteResult,
		Signature:     response.Signature,
		Sequence:      response.Sequence,
		TimeStamp:     response.TimeStamp,
		Hash:          response.Hash,
		Height:        response.Height,
	}, nil
}

func (r *cyberveinService) ExecuteAsync(c context.Context, req *proto.CommandRequest) (*proto.ExecuteAsyncResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.ExecuteAsync(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.ExecuteAsyncResponse{
		Cmd:       response.Cmd,
		Signature: response.Signature,
		Sequence:  response.Sequence,
		TimeStamp: response.TimeStamp,
		Hash:      response.Hash,
	}, nil
}

func (r *cyberveinService) ExecuteWithPrivateKey(c context.Context, req *proto.CommandRequest) (*proto.ExecuteResponse, error) {
	err := r.CheckToken(c)
	if err != nil {
		return nil, err
	}
	response, err := core.AppService.ExecuteWithPrivateKey(&models.CommandRequest{Cmd: req.Cmd})
	if err != nil {
		return nil, err
	}
	return &proto.ExecuteResponse{
		Cmd:           response.Cmd,
		ExecuteResult: response.ExecuteResult,
		Signature:     response.Signature,
		Sequence:      response.Sequence,
		TimeStamp:     response.TimeStamp,
		Hash:          response.Hash,
		Height:        response.Height,
	}, nil
}
