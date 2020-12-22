package routers

import (
	"cybervein.org/CyberveinDB/routers/handlers"
	"cybervein.org/CyberveinDB/routers/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/login", handlers.Login)
	db.Use(middleware.Authorization())
	db.Use(middleware.Log())
	{
		db.GET("/query", handlers.QueryCommand)
		db.GET("/query_private", handlers.QueryPrivateCommand)
		db.POST("/execute", handlers.ExecuteCommand)
		db.POST("/restore", handlers.RestoreLocalDatabase)
	}

	chain := r.Group("/chain")
	db.Use(middleware.Authorization())
	db.Use(middleware.Log())

	{
		chain.GET("/transaction", handlers.GetTransactionByHash)
		chain.GET("/transaction_list", handlers.GetCommittedTxList)
		chain.GET("/block", handlers.GetBlockByHeight)
		chain.GET("/state", handlers.GetChainState)
		chain.GET("/info", handlers.GetChainInfo)
		chain.GET("/net", handlers.GetNetInfo)
		chain.GET("/unconfirmed_tx", handlers.GetUnconfirmedTxs)
		chain.GET("/key_log", handlers.GetKeyLog)
		chain.GET("/genesis", handlers.GetGenesis)
		chain.GET("/voting_validators", handlers.GetVotingValidator)
		chain.POST("/update_validators", handlers.UpdateValidators)
	}

	return r
}
