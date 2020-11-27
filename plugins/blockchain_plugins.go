package plugins

import (
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/utils"
)

var pluginsMap map[string]BlockChainPlugin

type BlockChainPlugin interface {
	CustomChainInitMethod()
	CustomNewBlockEventMethod(block *models.Block)
	CustomTxValidationCheck(tx []byte) (bool, string)
	CustomTransactionDeliverLog(tx []byte, result string) string
}

func register(name string, plugin interface{}) {
	if pluginsMap == nil {
		pluginsMap = make(map[string]BlockChainPlugin, 0)
	}
	pluginsMap[name] = plugin.(BlockChainPlugin)
}

func GetConfigPlugin() BlockChainPlugin {
	var plugin BlockChainPlugin
	customPlugin := utils.Config.App.Plugin
	if pluginsMap == nil || !containsPlugin(customPlugin) {
		plugin = &DefaultBlockChainPlugin{}
	} else {
		plugin = pluginsMap[customPlugin]
	}
	return plugin
}

func containsPlugin(name string) bool {
	_, ok := pluginsMap[name]
	return ok
}
