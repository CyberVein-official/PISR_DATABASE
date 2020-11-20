package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"cybervein.org/CyberveinDB/core"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/models/code"
	"github.com/gin-gonic/gin"
)

func ExecuteCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.ExecuteRequest{}
	ginMsg.DecodeRequestBody(request)
	var err error
	var res interface{}

	if strings.EqualFold(request.Mode, "async") {
		res, err = core.AppService.ExecuteAsync(&models.CommandRequest{request.Cmd})
	} else if strings.EqualFold(request.Mode, "commit") {
		res, err = core.AppService.Execute(&models.CommandRequest{request.Cmd})
	} else if strings.EqualFold(request.Mode, "private") {
		res, err = core.AppService.ExecuteWithPrivateKey(&models.CommandRequest{request.Cmd})
	} else {
		ginMsg.CommonResponse(http.StatusOK, code.CodeTypeInvalidExecuteMode, fmt.Sprintf("Invalid mode : %s", request.Mode))
		return
	}

	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypecyberveinExecuteError, code.CodeTypecyberveinExecuteErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}


func QueryCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.Query(&models.CommandRequest{request.Cmd})
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypecyberveinQueryError, code.CodeTypecyberveinQueryErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func QueryPrivateCommand(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommandRequest{}
	ginMsg.DecodeRequestBody(request)
	addr := c.Query("address")
	var res *models.QueryResponse
	var err error
	if len(addr) != 0 {
		res, err = core.AppService.QueryPrivateDataWithAddress(&models.QueryPrivateWithAddrRequest{request.Cmd, strings.ToUpper(addr)})
	} else {
		res, err = core.AppService.QueryPrivateData(&models.CommandRequest{request.Cmd})
	}
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypecyberveinQueryError, code.CodeTypecyberveinQueryErrorMsg, err)
		return
	}
	ginMsg.Response(http.StatusOK, res)
}


func RestoreLocalDatabase(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	core.LogStoreApp.State.Lock()
	err := core.AppService.RestoreLocalDatabase()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeInternalError, code.CodeTypeInternalErrorMsg, err)
		return
	}
	core.LogStoreApp.State.UnLock()
	ginMsg.Success()
}
