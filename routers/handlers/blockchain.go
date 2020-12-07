package handlers

import (
	"cybervein.org/CyberveinDB/core"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/models/code"
	"github.com/gin-gonic/gin"

	"net/http"
)

func GetTransactionByHash(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.TxHashRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.GetTransaction(request.Hash)
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetCommittedTxList(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.CommittedTxListRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.GetCommittedTxList(request.Begin, request.End)
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}
// 
func GetBlockByHeight(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.BlockHeightRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.GetBlock(request.Height)
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetChainState(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res, err := core.AppService.GetChainState()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetChainInfo(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.ChainInfoRequest{}
	ginMsg.DecodeRequestBody(request)
	res, err := core.AppService.GetChainInfo(request.Min, request.Max)
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetNetInfo(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res, err := core.AppService.GetNetInfo()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetUnconfirmedTxs(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res, err := core.AppService.GetUnconfirmedTxs()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetKeyLog(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	key := c.Query("key")
	res, err := core.AppService.GetKeyLog(key)
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetGenesis(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res, err := core.AppService.GetGenesis()
	if err != nil {
		ginMsg.Error(http.StatusOK, code.CodeTypeGetChainInfoError, code.CodeTypeGetChainInfoErrorMsg, err)
		return
	}
	ginMsg.SuccessWithData(res)
}

func GetVotingValidator(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	res := core.AppService.QueryVotingValidators()
	ginMsg.Response(http.StatusOK, res)
}

func UpdateValidators(c *gin.Context) {
	ginMsg := models.GinMsg{C: c}
	request := &models.ValidatorUpdateData{}
	ginMsg.DecodeRequestBody(request)
	res := core.AppService.UpdateValidators(request)
	ginMsg.Response(http.StatusOK, res)
}
