package plugins

import (
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
)

type DefaultBlockChainPlugin struct {
}

func (d DefaultBlockChainPlugin) CustomChainInitMethod() {
	logger.Log.Info("Default block chain plugin CustomChainInitMethod: do nothing")
}

func (d DefaultBlockChainPlugin) CustomNewBlockEventMethod(block *models.Block) {
	logger.Log.Info("Default block chain plugin CustomNewBlockEventMethod: do nothing")
}

func (d DefaultBlockChainPlugin) CustomTxValidationCheck(tx []byte) (bool, string) {
	logger.Log.Info("Default transaction plugin CustomTxValidationCheck: do nothing")
	return true, ""
}

func (d DefaultBlockChainPlugin) CustomTransactionDeliverLog(tx []byte, result string) string {
	logger.Log.Info("Default transaction plugin CustomTransactionDeliverLog: do nothing")
	return ""
}

func init() {
	register("default", &DefaultBlockChainPlugin{})
}
