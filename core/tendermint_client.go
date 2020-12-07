package core

import (
	"context"
	"fmt"

	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/utils"
	c "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"strconv"
	"time"

	"github.com/tendermint/tendermint/types"
)

var tendermintHttpClient *c.HTTP

func InitClient() {
	var host = "tcp://" + utils.Config.Tendermint.Url
	var wsEndpoint = "/websocket"
	tendermintHttpClient = c.NewHTTP(host, wsEndpoint)
}

func BroadcastTxCommit(op *models.TxCommitBody) (*ctypes.ResultBroadcastTxCommit, error) {

	tx := types.Tx(utils.StructToJson(op))
	logger.Log.Info("Tendermint BroadcastTxCommit: " + string(tx))

	start := time.Now()
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxCommit(tx)
	end := time.Now()

	if err != nil {
		err = fmt.Errorf("BroadcastTxCommit command error : %s, %s", tx, err)
		logger.Log.Error(err)
		return nil, err
	}
	logger.Log.Infof("%v | Tendermint BroadcastTxCommit response: %s ", end.Sub(start), string(utils.StructToJson(resultBroadcastTxCommit)))

	return resultBroadcastTxCommit, nil
}

func BroadcastTxSync(op *models.TxCommitBody) (*ctypes.ResultBroadcastTx, error) {

	tx := types.Tx(utils.StructToJson(op))
	logger.Log.Info("Tendermint BroadcastTxSync: " + string(tx))

	start := time.Now()
	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxSync(tx)
	end := time.Now()

	if err != nil {
		err = fmt.Errorf("BroadcastTxSync command error : %s, %s", tx, err)
		logger.Log.Error(err)
		return nil, err
	}

	logger.Log.Infof("%v | Tendermint BroadcastTxCommit response: %s ", end.Sub(start), string(utils.StructToJson(resultBroadcastTxCommit)))

	return resultBroadcastTxCommit, nil
}

func ABCIDataQuery(path string, data []byte) *ctypes.ResultABCIQuery {

	resultABCIQuery, err := tendermintHttpClient.ABCIQuery(path, data)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	return resultABCIQuery
}

func GetTx(hash []byte) (*ctypes.ResultTx, error) {

	logger.Log.Info("Tendermint GetTx: " + string(hash))
	resultTx, err := tendermintHttpClient.Tx(hash, true)
	if err != nil {
		err = fmt.Errorf("get transaction by hash error : %s, %s", utils.ByteToHex(hash), err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultTx, nil
}

func GetChainInfo(min int, max int) (*ctypes.ResultBlockchainInfo, error) {

	minH := int64(min)
	maxH := int64(max)
	logger.Log.Info("Tendermint GetChainInfo: " + strconv.Itoa(min) + strconv.Itoa(max))

	resultBlockchainInfo, err := tendermintHttpClient.BlockchainInfo(minH, maxH)
	if err != nil {
		err = fmt.Errorf("get chain info error : %s", err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultBlockchainInfo, nil
}

func GetChainState() (*ctypes.ResultStatus, error) {

	logger.Log.Info("Tendermint GetChainState ")
	resultStatus, err := tendermintHttpClient.Status()
	if err != nil {
		err = fmt.Errorf("get chain state error : %s", err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultStatus, nil
}

func GetNetInfo() (*ctypes.ResultNetInfo, error) {
	resultNetInfo, err := tendermintHttpClient.NetInfo()
	if err != nil {
		err = fmt.Errorf("get net info error : %s", err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultNetInfo, nil
}

func GetUnconfirmedTxs() (*ctypes.ResultUnconfirmedTxs, error) {
	resultUnconfirmedTxs, err := tendermintHttpClient.NumUnconfirmedTxs()
	if err != nil {
		err = fmt.Errorf("get unconfirmed txs error : %s", err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultUnconfirmedTxs, nil
}

func GetBlockFromHeight(h int64) (*ctypes.ResultBlock, error) {

	logger.Log.Info("Tendermint GetBlockFromHeight : " + strconv.Itoa(int(h)))

	resultBlock, err := tendermintHttpClient.Block(&h)
	if err != nil {
		err = fmt.Errorf("get block error : %s, height : %d", err, h)
		logger.Log.Error(err)
		return nil, err
	}
	return resultBlock, nil
}

func UpdateValidator(update *models.ValidatorUpdateBody) (*ctypes.ResultBroadcastTxCommit, error) {
	tx := types.Tx(utils.StructToJson(update))
	logger.Log.Info("Tendermint UpdateValidator: " + string(tx))

	resultBroadcastTxCommit, err := tendermintHttpClient.BroadcastTxCommit(tx)
	if err != nil {
		err = fmt.Errorf("BroadcastTxCommit command error : %s, %s", tx, err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultBroadcastTxCommit, nil
}

func SubScribeEvent(event string) (out <-chan ctypes.ResultEvent, err error) {

	for !tendermintHttpClient.IsRunning() {
		time.Sleep(time.Second)
		tendermintHttpClient.Start()
	}
	eventQuery := "tm.event= '" + event + "'"
	out, err = tendermintHttpClient.Subscribe(context.Background(), "", eventQuery, 1024)
	if err != nil {
		err = fmt.Errorf("SubScribeEvent error : %s, %s", eventQuery, err)
		logger.Log.Error(err)
		return nil, err
	}
	return out, nil
}

func UnSubScribeEvent(event string) error {
	eventQuery := "tm.event='" + event + "'"
	err := tendermintHttpClient.WSEvents.Unsubscribe(context.Background(), "", eventQuery)
	if err != nil {
		err = fmt.Errorf("UnSubScribeEvent error : %s, %s", eventQuery, err)
		logger.Log.Error(err)
		return err
	}
	return nil
}

func GetGenesis() (*ctypes.ResultGenesis, error) {
	resultGenesis, err := tendermintHttpClient.Genesis()
	if err != nil {
		err = fmt.Errorf("GetGenesis error : % %s", err)
		logger.Log.Error(err)
		return nil, err
	}
	return resultGenesis, nil
}
