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

//#################### DeliverTx ####################
func (app *LogStoreApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	if app.IsValidatorUpdateTx(req.Tx) {
		return app.execValidatorTx(req.Tx)
	}

	code, info := app.isValid(req.Tx)
	if code != c.CodeTypeOK {
		return abcitypes.ResponseDeliverTx{Code: code, Info: info}
	}

	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(req.Tx, &commitBody)
	operationLog := models.OperationLog{
		Operation: commitBody.Data.Operation,
		Address:   commitBody.Address,
		Signature: commitBody.Signature,
		Time:      time.Now().String(),
		Height:    strconv.Itoa(int(app.State.currentHeight)),
		Sequence:  commitBody.Data.Sequence,
	}

	if app.IsPrivateCommand(commitBody) {
		if app.initFlag || !strings.EqualFold(commitBody.Address, utils.ValidatorKey.Address.String()) {
			if !app.privateTxSet.Contains(commitBody.Data.Sequence) {
				app.ExecuteCommand(commitBody.Data.Operation, ReTry, ReTryTime, true)
			} else {
				app.privateTxSet.Remove(commitBody.Data.Sequence)

				app.State.Lock()
				app.State.UpdateLogSequence()
				app.State.UpdateCurrentTxIndex()
				app.State.UnLock()
			}
		}
		database.UpdateSingleKeyWriteLog(operationLog)
		return abcitypes.ResponseDeliverTx{Code: c.CodeTypeOK, Info: c.CodeTypeOKMsg}
	} else {
		res := app.ExecuteCommand(commitBody.Data.Operation, ReTry, ReTryTime, true)
		log := app.plugin.CustomTransactionDeliverLog(req.Tx, res)
		database.UpdateSingleKeyWriteLog(operationLog)
		return abcitypes.ResponseDeliverTx{Code: c.CodeTypeOK, Data: []byte("Result:" + res), Log: log}
	}
}

//#################### Commit ####################
func (app *LogStoreApplication) Commit() abcitypes.ResponseCommit {
	return abcitypes.ResponseCommit{Data: []byte{}}
}

//#################### EndBlock ####################
func (app *LogStoreApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	app.State.CommitHeight()
	return abcitypes.ResponseEndBlock{ValidatorUpdates: app.valUpdates}
}

func (app *LogStoreApplication) execValidatorTx(tx []byte) abcitypes.ResponseDeliverTx {
	update := models.ValidatorUpdateBody{}
	utils.JsonToStruct(tx, &update)

	pubkeyS, powerS := update.ValidatorUpdate.PublicKey, update.ValidatorUpdate.Power

	// decode the pubkey
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyS)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: c.CodeTypeEncodingError,
			Log:  c.CodeTypeEncodingErrorMsg}
	}

	// decode the power
	power, err := strconv.ParseInt(powerS, 10, 64)
	if err != nil {
		return abcitypes.ResponseDeliverTx{
			Code: c.CodeTypeEncodingError,
			Log:  c.CodeTypeEncodingErrorMsg}
	}

	//publicKey:power
	voteKey := pubkeyS + VoteKeySep + powerS

	app.valAddrVote.addVote(voteKey, update.Address, update)
	count := app.valAddrVote[voteKey]
	if app.valAddrVote.getVoteNum(voteKey)*3 >= len(app.valAddrToPubKeyMap)*2 {
		app.updateValidator(abcitypes.Ed25519ValidatorUpdate(pubkey, power))
		app.committedValidatorVote[voteKey] = count
		delete(app.valAddrVote, voteKey)
	}

	return abcitypes.ResponseDeliverTx{Code: code.CodeTypeOK, Data: utils.StructToJson(count)}
}

func (app *LogStoreApplication) QueryVotingValidators() *Vote {
	return &app.valAddrVote
}

func (app *LogStoreApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	//不使用Tendermint查询接口，查询功能通过本地redis提供
	return
}

func (app *LogStoreApplication) updateValidator(v abcitypes.ValidatorUpdate) abcitypes.ResponseDeliverTx {

	pubkey := ed25519.PubKeyEd25519{}
	copy(pubkey[:], v.PubKey.Data)
	val := database.GetBadgerVal([]byte("val:" + string(v.PubKey.Data)))

	if v.Power == 0 {
		// remove validator
		if val != nil {
			database.DeleteBadgerKey([]byte("val:" + string(v.PubKey.Data)))
			delete(app.valAddrToPubKeyMap, pubkey.Address().String())

		} else {
			pubStr := base64.StdEncoding.EncodeToString(v.PubKey.Data)
			return abcitypes.ResponseDeliverTx{
				Code: c.CodeTypeEncodingError,
				Log:  c.CodeTypeEncodingErrorMsg + " : " + pubStr}

		}

	} else {
		// add or update validator
		value := bytes.NewBuffer(make([]byte, 0))
		if err := abcitypes.WriteMessage(&v, value); err != nil {
			return abcitypes.ResponseDeliverTx{
				Code: c.CodeTypeEncodingError,
				Log:  c.CodeTypeEncodingErrorMsg + " : " + err.Error()}
		}
		database.UpdateBadgerVal([]byte("val:"+string(v.PubKey.Data)), value.Bytes())
		app.valAddrToPubKeyMap[pubkey.Address().String()] = v.PubKey
	}

	app.valUpdates = append(app.valUpdates, v)

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (LogStoreApplication) IsValidatorUpdateTx(tx []byte) bool {
	commitBody := models.TxCommitBody{}
	utils.JsonToStruct(tx, &commitBody)

	if commitBody.Data != nil {
		return false
	}
	return true
}

func (app *LogStoreApplication) IsPrivateCommand(commitBody models.TxCommitBody) bool {
	key, err := database.GetKey(commitBody.Data.Operation)
	if err != nil {
		logger.Log.Error(err)
		return false
	}
	split := strings.Split(key, PrivateSep)
	if _, ok := app.valAddrToPubKeyMap[split[0]]; ok &&  strings.EqualFold(split[0], commitBody.Address) {
		return true
	}
	return false 
}


func (app *LogStoreApplication) ExecuteCommand(command string, retry int, t time.Duration, updateState bool) string {
	app.State.lock.Lock()
	defer app.State.lock.Unlock()
	res, err := database.ExecuteCommand(command)
	fmt.Println(err)
	if err != nil {
		if database.CheckAlive(3) {
			return fmt.Sprintf("Redis Command Error : %s", err)
		} else {
			for i := 0; i < retry; i++ {
				time.Sleep(t)
				res, err = database.ExecuteCommand(command)
				if err == nil {
					break
				}
			}

			if err != nil {
				//重启后还是失败则抛出系统异常
				AppService.RestoreLocalDatabase()
				res, err = database.ExecuteCommand(command)
				if err != nil {
					panic(err)
				}
			}
		}
	}
	if !updateState {
		app.State.UpdateLogSequence()
		app.State.UpdateCurrentTxIndex()
	}
	return res
}
