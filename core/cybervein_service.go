package core

import "cybervein.org/CyberveinDB/models"

type Service interface {
	//database
	RestoreLocalDatabase() error
	Query(request *models.CommandRequest) (*models.QueryResponse, error)
	QueryPrivateData(request *models.CommandRequest) (*models.QueryResponse, error)
	QueryPrivateDataWithAddress(request *models.QueryPrivateWithAddrRequest) (*models.QueryResponse, error)
	Execute(request *models.CommandRequest) (*models.ExecuteResponse, error)
	ExecuteAsync(request *models.CommandRequest) (*models.ExecuteAsyncResponse, error)
	ExecuteWithPrivateKey(request *models.CommandRequest) (*models.ExecuteResponse, error)

	//chain
	GetBlock(height int) (*models.Block, error)
	GetGenesis() (*models.Genesis, error)
	GetChainInfo(min int, max int) (*models.ChainInfo, error)
	GetChainState() (*models.ChainState, error)
	GetTransaction(hash string) (*models.Transaction, error)
	GetCommittedTxList(beginHeight int, endHeight int) (*models.TransactionCommittedList, error)
	GetUnconfirmedTxs() (*models.UnConfirmedTxs, error)
	GetNetInfo() (*models.NetInfo, error)
	GetKeyLog(key string) (*models.OperationKeyLog, error)

	//validators
	QueryVotingValidators() *Vote
	UpdateValidators(update *models.ValidatorUpdateData) *VoteCount

	//log writer
	StartCommandLogWriter()
}
