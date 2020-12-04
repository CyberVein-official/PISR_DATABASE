package core

import (
	"fmt"
	"cybervein.org/CyberveinDB/database"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/utils"
	uuid "github.com/satori/go.uuid"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	"strconv"
	"strings"
	"sync"
	"time"
)

var AppService Service

type ApplicationService struct {
}

func InitService() {
	AppService = &ApplicationService{}
}

func (s *ApplicationService) MakeTxCommitBody(request *models.CommandRequest) *models.TxCommitBody {
	op := &models.TxCommitBody{}
	op.Data = &models.TxCommitData{}

	//Sequence
	u := uuid.NewV1()
	op.Data.Sequence = utils.ByteToHex(u.Bytes())

	//cmd
	op.Data.Operation = request.Cmd

	//Signature
	op.Signature = utils.ValidatorStringSign(utils.StructToJson(op.Data))

	//address
	op.Address = utils.ValidatorKey.Address.String()

	return op
}
func (s *ApplicationService) Execute(request *models.CommandRequest) (*models.ExecuteResponse, error) {
	time.Sleep(time.Duration(200)*time.Millisecond)

	op := s.MakeTxCommitBody(request)
	timestamp := time.Now().UnixNano() / 1e6
	commitMsg, err := BroadcastTxCommit(op)
	if err != nil {
		return nil, err
	}


	return &models.ExecuteResponse{
		Cmd:           request.Cmd,
		ExecuteResult: string(commitMsg.DeliverTx.Data),
		Signature:     op.Signature,
		Sequence:      op.Data.Sequence,
		TimeStamp:     strconv.FormatInt(timestamp, 10),
		Hash:          utils.ByteToHex(commitMsg.Hash),
		Height:        commitMsg.Height}, nil
}

func (s *ApplicationService) ExecuteAsync(request *models.CommandRequest) (*models.ExecuteAsyncResponse, error) {
	op := s.MakeTxCommitBody(request)
	timestamp := time.Now().UnixNano() / 1e6

	sync, err := BroadcastTxSync(op)
	if err != nil {
		return nil, err
	}

	return &models.ExecuteAsyncResponse{
		Cmd:       op.Data.Operation,
		Signature: op.Signature,
		Sequence:  op.Data.Sequence,
		TimeStamp: strconv.FormatInt(timestamp, 10),
		Hash:      utils.ByteToHex(sync.Hash),
	}, nil
}

func (s *ApplicationService) ExecuteWithPrivateKey(request *models.CommandRequest) (*models.ExecuteResponse, error) {
	time.Sleep(time.Duration(5)*time.Millisecond)

	key, err := database.GetKey(request.Cmd)
	if err != nil {
		logger.Log.Error(err)
	}
	key = utils.ValidatorKey.Address.String() + PrivateSep + key
	cmd, err := database.ReplaceKey(request.Cmd, key)
	if err != nil {
		logger.Log.Error(err)
	}
	request.Cmd = cmd

	op := s.MakeTxCommitBody(request)
	timestamp := time.Now().UnixNano() / 1e6

	sync, err := BroadcastTxSync(op)
	if err != nil {
		return nil, err
	}

	return &models.ExecuteResponse{
		Cmd:           op.Data.Operation,
		ExecuteResult: string(sync.Data),
		Signature:     op.Signature,
		Sequence:      op.Data.Sequence,
		TimeStamp:     strconv.FormatInt(timestamp, 10),
		Hash:          utils.ByteToHex(sync.Hash),
	}, nil
}

func (s *ApplicationService) Query(request *models.CommandRequest) (*models.QueryResponse, error) {
	result, err := database.ExecuteCommand(request.Cmd)
	if err != nil {
		return nil, err
	}
	return &models.QueryResponse{Result: result}, nil
}

func (s *ApplicationService) QueryPrivateData(request *models.CommandRequest) (*models.QueryResponse, error) {

	key, err := database.GetKey(request.Cmd)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	key = utils.ValidatorKey.Address.String() + PrivateSep + key
	cmd, err := database.ReplaceKey(request.Cmd, key)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	request.Cmd = cmd
	result, err := database.ExecuteCommand(request.Cmd)
	if err != nil {
		return nil, err
	}
	return &models.QueryResponse{Result: result}, nil
}

func (s *ApplicationService) QueryPrivateDataWithAddress(request *models.QueryPrivateWithAddrRequest) (*models.QueryResponse, error) {

	key, err := database.GetKey(request.Cmd)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	key = strings.ToUpper(request.Address) + PrivateSep + key
	cmd, err := database.ReplaceKey(request.Cmd, key)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	request.Cmd = cmd
	result, err := database.ExecuteCommand(request.Cmd)
	if err != nil {
		return nil, err
	}
	return &models.QueryResponse{Result: result}, nil
}

func (s *ApplicationService) RestoreLocalDatabase() error {
	txs, err := utils.ReadTxFromDBLogFile(2, int(LogStoreApp.State.GetCurrentHeight()))
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	database.RestartRedisServer()
	logger.Log.Info("Restore Local Database Begin ... ")
	for _, c := range txs {
		database.ExecuteCommand(c)
		logger.Log.Info("Restore Execute command: " + c)
	}
	logger.Log.Info("Restore Local Database Finished... ")
	return nil
}

func (s *ApplicationService) GetTransaction(hash string) (*models.Transaction, error) {
	byteHash := utils.HexToByte(hash)
	tx, err := GetTx(byteHash)
	if err != nil {
		return nil, err
	}
	data := &models.TxCommitBody{}
	utils.JsonToStruct(tx.Tx, data)

	transaction := &models.Transaction{
		Hash:          utils.ByteToHex(tx.Hash),
		Height:        tx.Height,
		Index:         tx.Index,
		Data:          data,
		ExecuteResult: string(tx.TxResult.Data),
		Proof:         nil,
	}
	transaction.Proof = &models.TxProof{
		RootHash: utils.ByteToHex(tx.Proof.RootHash),
		Proof: &models.ProofDetail{
			Total:    tx.Proof.Proof.Total,
			Index:    tx.Proof.Proof.Index,
			LeafHash: "",
			Aunts:    nil,
		},
	}

	if tx.Proof.Proof.LeafHash != nil {
		transaction.Proof.Proof.LeafHash = utils.ByteToHex(tx.Proof.Proof.LeafHash)
	}

	if len(tx.Proof.Proof.Aunts) != 0 {
		transaction.Proof.Proof.Aunts = make([]string, 0)
		for _, v := range tx.Proof.Proof.Aunts {
			transaction.Proof.Proof.Aunts = append(transaction.Proof.Proof.Aunts, utils.ByteToHex(v))
		}
	}
	return transaction, nil
}

func (s *ApplicationService) GetCommittedTxList(beginHeight int, endHeight int) (*models.TransactionCommittedList, error) {

	txList := &models.TransactionCommittedList{
		Total: 0,
		Data:  make([]*models.CommittedTx, 0),
	}
	for i := beginHeight; i <= endHeight; i++ {
		originBlock, err := GetBlockFromHeight(int64(i))
		if err != nil {
			return nil, err
		}
		for _, tx := range originBlock.Block.Txs {
			data := &models.CommittedTx{}
			utils.JsonToStruct(tx, data)
			data.Height = int64(i)
			txList.Data = append(txList.Data, data)
		}
	}
	txList.Total = int64(len(txList.Data))
	return txList, nil
}

func (s *ApplicationService) GetBlock(height int) (*models.Block, error) {
	originBlock, err := GetBlockFromHeight(int64(height))
	if err != nil {
		return nil, err
	}
	return s.ConvertBlock(originBlock), nil
}

func (s *ApplicationService) GetChainInfo(min int, max int) (*models.ChainInfo, error) {
	info, err := GetChainInfo(min, max)
	if err != nil {
		return nil, err
	}
	res := &models.ChainInfo{
		LastHeight: info.LastHeight,
		BlockMetas: make([]*models.BlockMeta, 0),
	}
	metas := info.BlockMetas
	for _, v := range metas {
		res.BlockMetas = append(res.BlockMetas, &models.BlockMeta{
			BlockID: *s.ConvertBlockID(&v.BlockID),
			Header:  *s.ConvertBlockHeader(&v.Header),
		})
	}
	return res, nil
}

func (s *ApplicationService) GetChainState() (*models.ChainState, error) {
	originState, err := GetChainState()
	if err != nil {
		return nil, err
	}
	state := &models.ChainState{}
	utils.JsonToStruct(utils.StructToJson(originState), state)
	state.ValidatorInfo.PubKey = originState.ValidatorInfo.PubKey.Bytes()
	state.ValidatorInfo.VotingPower = originState.ValidatorInfo.VotingPower
	return state, nil
}

func (s *ApplicationService) GetGenesis() (*models.Genesis, error) {
	resultGenesis, err := GetGenesis()
	if err != nil {
		return nil, err
	}

	o := resultGenesis.Genesis
	consensusParams := &models.ConsensusParams{}
	utils.JsonToStruct(utils.StructToJson(o.ConsensusParams), consensusParams)

	validators := make([]models.GenesisValidator, 0)
	for _, v := range o.Validators {
		validators = append(validators, models.GenesisValidator{
			Address: v.Address.Bytes(),
			PubKey:  v.PubKey.Bytes(),
			Power:   v.Power,
			Name:    v.Name,
		})
	}

	return &models.Genesis{
		GenesisTime:     o.GenesisTime,
		ChainID:         o.ChainID,
		ConsensusParams: consensusParams,
		Validators:      validators,
		AppHash:         []byte(o.AppHash),
		AppState:        o.AppState,
	}, nil
}

func (s *ApplicationService) GetNetInfo() (*models.NetInfo, error) {
	resultNetInfo, err := GetNetInfo()
	if err != nil {
		return nil, err
	}
	netInfo := &models.NetInfo{}
	utils.JsonToStruct(utils.StructToJson(resultNetInfo), netInfo)
	return netInfo, nil
}

func (s *ApplicationService) GetUnconfirmedTxs() (*models.UnConfirmedTxs, error) {
	resultUnconfirmedTxs, err := GetUnconfirmedTxs()
	if err != nil {
		return nil, err
	}
	txs := models.UnConfirmedTxs{}
	utils.JsonToStruct(utils.StructToJson(resultUnconfirmedTxs), txs)

	txList := make([]models.TxCommitBody, 0)

	for _, tx := range resultUnconfirmedTxs.Txs {
		data := &models.TxCommitBody{}
		utils.JsonToStruct(tx, data)
		txList = append(txList, *data)
	}
	txs.Txs = txList
	return &txs, nil
}

func (s *ApplicationService) GetKeyLog(key string) (*models.OperationKeyLog, error) {
	return database.GetKeyWriteLog(key)
}

func (s *ApplicationService) QueryVotingValidators() *Vote {
	return LogStoreApp.QueryVotingValidators()
}

func (s *ApplicationService) UpdateValidators(update *models.ValidatorUpdateData) *VoteCount {
	
	//Sequence
	u := uuid.NewV1()
	sequence := utils.ByteToHex(u.Bytes())
	
	updateBody := &models.ValidatorUpdateBody{
		ValidatorUpdate: update,
		Signature:       "",
		Address:         "",
		Sequence:	sequence,
		
	}

	updateBody.Signature = utils.ValidatorStringSign(utils.StructToJson(updateBody.ValidatorUpdate))
	updateBody.Address = utils.ValidatorKey.Address.String()

	commit, err := UpdateValidator(updateBody)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	vote := &VoteCount{}
	utils.JsonToStruct(commit.DeliverTx.Data, vote)
	return vote
}

func (s *ApplicationService) StartCommandLogWriter() {
	var lock sync.Mutex
	out, err := SubScribeEvent("NewBlock")
	if err != nil {
		logger.Log.Errorf("Subscribe event NewBlock failed : ", err)
	}
	logger.Log.Info("Subscribe tendermint event : NewBlock")

	state, err := s.GetChainState()
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}

	height := state.SyncInfo.LatestBlockHeight
	if height != 1 {
		for i := 2; i <= int(height); i++ {
			queryBlock, err := GetBlockFromHeight(int64(i))
			if err != nil {
				logger.Log.Error(err)
				break
			}
			strList := s.ConvertTransactionsToLogString(queryBlock.Block.Txs, queryBlock.Block.Height, queryBlock.Block.Time.String())
			utils.AppendToDBLogFile(strList)
		}
	}

	for {
		select {
		case resultEvent := <-out:
			logger.Log.Info("NewBlock event")
			block := resultEvent.Data.(types.EventDataNewBlock).Block
			strList := s.ConvertTransactionsToLogString(block.Txs, block.Height, block.Time.String())
			logList := s.ConvertTransactionsToCommandLog(block.Txs, block.Height, block.Time.String())
			lock.Lock()
			utils.AppendToDBLogFile(strList)
			database.UpdateKeyWriteLog(logList)
			LogStoreApp.plugin.CustomNewBlockEventMethod(s.ConvertBlockFromTypesBlock(block))
			lock.Unlock()
		}
	}
}

func (s *ApplicationService) ConvertTransactionsToCommandLog(txs types.Txs, h int64, time string) []models.OperationLog {
	operationLog := make([]models.OperationLog, 0)
	for _, v := range txs {
		data := &models.TxCommitBody{}
		utils.JsonToStruct(v, data)
		operationLog = append(operationLog,
			models.OperationLog{
				Operation: data.Data.Operation,
				Address:   data.Address,
				Signature: data.Signature,
				Sequence:  data.Data.Sequence,
				Height:    strconv.Itoa(int(h)),
				Time:      time,
			})
	}
	return operationLog
}

func (s *ApplicationService) ConvertTransactionsToLogString(txs types.Txs, h int64, time string) []string {
	strList := make([]string, 0)
	for _, v := range txs {
		data := &models.TxCommitBody{}
		utils.JsonToStruct(v, data)
		strList = append(strList, fmt.Sprintf("%s | %s | %s | %s | %d | %s ", data.Data.Operation, data.Address, data.Signature, data.Data.Sequence, h, time))
	}
	return strList
}

func (s *ApplicationService) ConvertBlockID(b *types.BlockID) *models.BlockID {
	blockID := models.BlockID{}
	blockID.Hash = utils.ByteToHex(b.Hash)
	blockID.PartsHeader = models.PartSetHeader{
		Total: b.PartsHeader.Total,
		Hash:  utils.ByteToHex(b.PartsHeader.Hash),
	}
	return &blockID
}

func (s *ApplicationService) ConvertBlockHeader(b *types.Header) *models.Header {
	return &models.Header{
		Version:            b.Version,
		ChainID:            b.ChainID,
		Height:             b.Height,
		Time:               time.Time{},
		NumTxs:             b.NumTxs,
		TotalTxs:           b.TotalTxs,
		LastBlockID:        *s.ConvertBlockID(&b.LastBlockID),
		LastCommitHash:     utils.ByteToHex(b.LastCommitHash),
		DataHash:           utils.ByteToHex(b.DataHash),
		ValidatorsHash:     utils.ByteToHex(b.ValidatorsHash),
		NextValidatorsHash: utils.ByteToHex(b.NextValidatorsHash),
		ConsensusHash:      utils.ByteToHex(b.ConsensusHash),
		AppHash:            utils.ByteToHex(b.AppHash),
		LastResultsHash:    utils.ByteToHex(b.LastCommitHash),
		EvidenceHash:       utils.ByteToHex(b.EvidenceHash),
		ProposerAddress:    utils.ByteToHex(b.ProposerAddress),
	}
}

func (s *ApplicationService) ConvertBlockData(b *types.Data) *models.Data {
	data := &models.Data{
		Txs:  make([]string, 0),
		Hash: utils.ByteToHex(b.Hash()),
	}
	for _, v := range b.Txs {
		data.Txs = append(data.Txs, fmt.Sprintf("%x", v.Hash()))
	}
	return data
}

func (s *ApplicationService) ConvertCommitSign(b *types.CommitSig) *models.CommitSig {
	return &models.CommitSig{
		Type:             b.Type,
		Height:           b.Height,
		Round:            b.Round,
		Timestamp:        b.Timestamp,
		ValidatorAddress: utils.ByteToHex(b.ValidatorAddress),
		ValidatorIndex:   b.ValidatorIndex,
		Signature:        utils.ByteToHex(b.Signature),
	}
}

func (s *ApplicationService) ConvertBlockFromTypesBlock(b *types.Block) *models.Block {
	header := s.ConvertBlockHeader(&(b.Header))

	data := s.ConvertBlockData(&(b.Data))

	lastCommit := make([]*models.CommitSig, 0)

	for _, v := range b.LastCommit.Precommits {
		lastCommit = append(lastCommit, s.ConvertCommitSign(v))
	}

	block := &models.Block{
		BlockID:    models.BlockID{},
		Header:     *header,
		Data:       *data,
		Evidence:   b.Evidence,
		LastCommit: lastCommit,
	}
	return block
}

func (s *ApplicationService) ConvertBlock(b *ctypes.ResultBlock) *models.Block {

	blockID := s.ConvertBlockID(&(b.BlockMeta.BlockID))

	header := s.ConvertBlockHeader(&(b.Block.Header))

	data := s.ConvertBlockData(&(b.Block.Data))

	lastCommit := make([]*models.CommitSig, 0)

	for _, v := range b.Block.LastCommit.Precommits {
		lastCommit = append(lastCommit, s.ConvertCommitSign(v))
	}

	block := &models.Block{
		BlockID:    *blockID,
		Header:     *header,
		Data:       *data,
		Evidence:   b.Block.Evidence,
		LastCommit: lastCommit,
	}

	return block
}
