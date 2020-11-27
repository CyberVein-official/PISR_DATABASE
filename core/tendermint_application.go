package core

import (
	"fmt"
	"strings"
	"time"

	"cybervein.org/CyberveinDB/database"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	c "cybervein.org/CyberveinDB/models/code"
	"cybervein.org/CyberveinDB/plugins"
	"cybervein.org/CyberveinDB/utils"
	"github.com/emirpasic/gods/sets/hashset"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
)

type LogStoreApplication struct {
	valUpdates             []abcitypes.ValidatorUpdate
	valAddrToPubKeyMap     map[string]abcitypes.PubKey
	valAddrVote            Vote
	committedValidatorVote Vote
	privateTxSet           *hashset.Set

	plugin   plugins.BlockChainPlugin
	initFlag bool

	State *AppState
}

var _ abcitypes.Application = (*LogStoreApplication)(nil)
var LogStoreApp *LogStoreApplication

const PrivateSep string = "_"
const VoteKeySep string = ":"
const SocketAddr string = "unix://tendermint.sock"
const ReTry int = 3
const ReTryTime = 10 * time.Millisecond

func InitLogStoreApplication() {
	LogStoreApp = &LogStoreApplication{
		valAddrToPubKeyMap:     make(map[string]abcitypes.PubKey),
		valAddrVote:            NewVote(),
		committedValidatorVote: NewVote(),
		initFlag:               true,
		privateTxSet:           hashset.New(),
		plugin:                 plugins.GetConfigPlugin(),
		State: &AppState{
			logSequence:            0,
			currentHeight:          1,
			currentTxIndex:         0,
			currentCommittedHeight: 1,
		},
	}
}

func (app *LogStoreApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *LogStoreApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *LogStoreApplication) isValid(tx []byte) (uint32, string) {
	var data []byte
	var sign string
	var address string

	if !app.IsValidatorUpdateTx(tx) {
		commitBody := models.TxCommitBody{}
		utils.JsonToStruct(tx, &commitBody)
		data = utils.StructToJson(commitBody.Data)
		sign = commitBody.Signature
		address = commitBody.Address
		if !database.IsValidCmd(commitBody.Data.Operation) {
			return c.CodeTypeInvalidTx, fmt.Sprintf("Invalid redis command : %s", commitBody.Data.Operation)
		}
		if database.IsQueryCmd(commitBody.Data.Operation) {
			return c.CodeTypeInvalidTx, fmt.Sprintf("Read only command can not commit to tendermint : %s", commitBody.Data.Operation)
		}
	} else {
		commitBody := models.ValidatorUpdateBody{}
		utils.JsonToStruct(tx, &commitBody)
		data = utils.StructToJson(commitBody.ValidatorUpdate)
		sign = commitBody.Signature
		address = commitBody.Address
	}

	if _, ok := app.valAddrToPubKeyMap[address]; !ok {
		return c.CodeTypeInvalidValidator, c.CodeTypeInvalidValidatorMsg
	}

	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], app.valAddrToPubKeyMap[address].Data)

	if pubkey.VerifyBytes(data, utils.HexToByte(sign)) != true {
		return c.CodeTypeInvalidSign, c.CodeTypeInvalidSignMsg
	}

	if b, msg := app.plugin.CustomTxValidationCheck(tx); !b {
		return c.CodeTypeInternalError, c.CodeTypeInternalErrorMsg + " : " + msg
	}

	return c.CodeTypeOK, c.CodeTypeOKMsg
}

//#################### InitChain ####################
func (app *LogStoreApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	//TODO 清空Redis, 重启的是时候会初始化
	for _, v := range req.Validators {
		r := app.updateValidator(v)
		if r.IsErr() {
			logger.Log.Error(r)
		}
	}
	app.plugin.CustomChainInitMethod()
	return abcitypes.ResponseInitChain{}
}

//#################### CheckTx ####################
func (app *LogStoreApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	app.initFlag = false
	code, info := app.isValid(req.Tx)
	if code != c.CodeTypeOK {
		return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info}
	}

	if !app.IsValidatorUpdateTx(req.Tx) {
		commitBody := models.TxCommitBody{}
		utils.JsonToStruct(req.Tx, &commitBody)
		if app.IsPrivateCommand(commitBody) && strings.EqualFold(commitBody.Address, utils.ValidatorKey.Address.String()) {
			res := app.ExecuteCommand(commitBody.Data.Operation, ReTry, ReTryTime, false)
			app.privateTxSet.Add(commitBody.Data.Sequence)
			return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info, Data: []byte("Result:" + res)}
		}
	}

	return abcitypes.ResponseCheckTx{Code: code, GasWanted: 1, Info: info}
}

//#################### BeginBlock ####################
func (app *LogStoreApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.State.UpdateCurrentHeight()

	//reset valUpdates
	//TODO 从上一轮的拜占庭节点去判断是否包含本节点, 如果包含, 做一次数据库修复
	app.valUpdates = make([]abcitypes.ValidatorUpdate, 0)
	for _, ev := range req.ByzantineValidators {
		if ev.Type == tmtypes.ABCIEvidenceTypeDuplicateVote {
			// decrease voting power by 1
			if ev.TotalVotingPower == 0 {
				continue
			}
			app.updateValidator(abcitypes.ValidatorUpdate{
				PubKey: app.valAddrToPubKeyMap[string(ev.Validator.Address)],
				Power:  ev.TotalVotingPower - 1,
			})
		}
	}
	return abcitypes.ResponseBeginBlock{}
}
