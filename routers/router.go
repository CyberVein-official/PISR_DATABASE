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
	return r
}
