package routers


import (
	"github.com/gin-gonic/gin"
)


func InitRouter() *gin.Engine {

	r := gin.Default()

	r.POST("/login", handlers.Login)

	return r
}



