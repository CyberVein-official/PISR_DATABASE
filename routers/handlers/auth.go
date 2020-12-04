package handlers

import (
	"net/http"

	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/models/code"
	"cybervein.org/CyberveinDB/utils"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.LoginRequest{}
	ginMsg.DecodeRequestBody(request)

	if request.Name != utils.Config.App.Name || request.Password != utils.Config.App.DbPassword {
		ginMsg.SimpleErrorMsg(http.StatusOK, code.CodeTypeIncorrectPassword, "Incorrect AppName or Password")
		return
	}

	s, err := utils.GenerateToken(c.ClientIP(), request.Name, request.Password)
	if err != nil {
		logger.Log.Error(err)
		ginMsg.SimpleErrorMsg(http.StatusOK, code.CodeTypeInternalError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  code.CodeTypeOK,
		"msg":   code.CodeTypeOKMsg,
		"token": s,
	})
}
