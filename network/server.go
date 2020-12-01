package network

import (
	"fmt"
	"net/http"
	"strconv"

	"cybervein.org/CyberveinDB/grpc"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/routers"
	"cybervein.org/CyberveinDB/utils"
	"github.com/gin-gonic/gin"
)

var AppServer *Server

type Server struct {
	httpPort   string
	rpcPort    string
	rpcServer  *grpc.Server
	httpServer *http.Server
}

func NewServer() *Server {

	server := &Server{
		httpPort:   strconv.Itoa(utils.Config.HttpServer.Port),
		rpcPort:    strconv.Itoa(utils.Config.Rpc.Port),
		rpcServer:  nil,
		httpServer: nil,
	}

	gin.SetMode(utils.Config.HttpServer.RunMode)
	routersInit := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", utils.Config.HttpServer.Port)

	server.httpServer = &http.Server{
		Addr:    endPoint,
		Handler: routersInit,
	}

	server.rpcServer = grpc.NewRpcServer(strconv.Itoa(utils.Config.Rpc.Port))
	AppServer = server

	return server
}

func (server *Server) Start() {

	logger.Log.Infof("Rpc Server will be started at :%s...", server.rpcPort)
	go server.rpcServer.StartServer()

	logger.Log.Infof("Start Command Log Writer ...")
	//go core.AppService.StartCommandLogWriter()

	logger.Log.Infof("Http Server will be started at :%s...", server.httpPort)
	if err := server.httpServer.ListenAndServe(); err != nil {
		logger.Log.Error(err)
		return
	}
}
