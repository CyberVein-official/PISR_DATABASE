package models

import (
	"encoding/json"
	"time"

	"cybervein.org/CyberveinDB/utils"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

//Execution
type CommandRequest struct {
	Cmd string `json:"cmd"`
}

type QueryPrivateWithAddrRequest struct {
	Cmd     string `json:"command"`
	Address string `json:"address"`
}

type ExecuteResponse struct {
	Cmd           string `json:"command"`
	ExecuteResult string `json:"execute_result"`
	Signature     string `json:"signature"`
	Sequence      string `json:"sequence"`
	TimeStamp     string `json:"time_stamp"`
	Hash          string `json:"hash"`
	Height        int64  `json:"height"`
}

type ExecuteAsyncResponse struct {
	Cmd       string `json:"command"`
	Signature string `json:"signature"`
	Sequence  string `json:"sequence"`
	TimeStamp string `json:"time_stamp"`
	Hash      string `json:"hash"`
}

type QueryResponse struct {
	Result string `json:"result"`
}

//Transaction
type Transaction struct {
	Hash          string        `json:"hash"`
	Height        int64         `json:"height"`
	Index         uint32        `json:"index"`
	Data          *TxCommitBody `json:"data"`
	ExecuteResult string        `json:"execute_result"`
	Proof         *TxProof      `json:"proof,omitempty"`
}

type TransactionList struct {
	Height int64         `json:"height"`
	Total  int64         `json:"total"`
	Txs    []Transaction `json:"txs"`
}

type TxProof struct {
	RootHash string       `json:"root_hash"`
	Proof    *ProofDetail `json:"proof"`
}

type ProofDetail struct {
	Total    int      `json:"total"`
	Index    int      `json:"index"`
	LeafHash string   `json:"leaf_hash"`
	Aunts    []string `json:"aunts"`
}

type TransactionCommittedList struct {
	Total int64          `json:"total"`
	Data  []*CommittedTx `json:"data"`
}

type CommittedTx struct {
	Data      *TxCommitData `json:"data"`
	Signature string        `json:"signature"`
	Address   string        `json:"address"`
	Height    int64         `json:"height"`
}

//Block
type BlockMeta struct {
	BlockID `json:"block_id"`
	Header  `json:"header"`
}

type Block struct {
	BlockID    `json:"block_id"`
	Header     `json:"header"`
	Data       `json:"data"`
	Evidence   types.EvidenceData `json:"evidence"`
	LastCommit []*CommitSig       `json:"last_commit"`
}

type BlockID struct {
	Hash        string        `json:"hash"`
	PartsHeader PartSetHeader `json:"parts"`
}

type PartSetHeader struct {
	Total int    `json:"total"`
	Hash  string `json:"hash"`
}

type Header struct {
	Version  version.Consensus `json:"version"`
	ChainID  string            `json:"chain_id"`
	Height   int64             `json:"height"`
	Time     time.Time         `json:"time"`
	NumTxs   int64             `json:"num_txs"`
	TotalTxs int64             `json:"total_txs"`

	LastBlockID BlockID `json:"last_block_id"`

	LastCommitHash string `json:"last_commit_hash"`
	DataHash       string `json:"data_hash"`

	ValidatorsHash     string `json:"validators_hash"`
	NextValidatorsHash string `json:"next_validators_hash"`
	ConsensusHash      string `json:"consensus_hash"`
	AppHash            string `json:"app_hash"`
	LastResultsHash    string `json:"last_results_hash"`

	EvidenceHash    string `json:"evidence_hash"`
	ProposerAddress string `json:"proposer_address"`
}

type Data struct {
	Txs  []string `json:"txs"`
	Hash string   `json:"hash"`
}

type CommitSig struct {
	Type             types.SignedMsgType `json:"type"`
	Height           int64               `json:"height"`
	Round            int                 `json:"round"`
	Timestamp        time.Time           `json:"timestamp"`
	ValidatorAddress string              `json:"validator_address"`
	ValidatorIndex   int                 `json:"validator_index"`
	Signature        string              `json:"signature"`
}

//chain
type ChainInfo struct {
	LastHeight int64        `json:"last_height"`
	BlockMetas []*BlockMeta `json:"block_metas"`
}

type NodeInfo struct {
	ProtocolVersion p2p.ProtocolVersion `json:"protocol_version"`
	ID              p2p.ID              `json:"id"`
	ListenAddr      string              `json:"listen_addr"`
	Network         string              `json:"network"`
	Version         string              `json:"version"`
	Channels        HexBytes            `json:"channels"`
	Moniker         string              `json:"moniker"`
}

type SyncInfo struct {
	LatestBlockHash   HexBytes  `json:"latest_block_hash"`
	LatestAppHash     HexBytes  `json:"latest_app_hash"`
	LatestBlockHeight int64     `json:"latest_block_height"`
	LatestBlockTime   time.Time `json:"latest_block_time"`
	CatchingUp        bool      `json:"catching_up"`
}

type ValidatorInfo struct {
	Address     HexBytes `json:"address"`
	PubKey      HexBytes `json:"pub_key"`
	VotingPower int64    `json:"voting_power"`
}

type ChainState struct {
	NodeInfo      NodeInfo      `json:"node_info"`
	SyncInfo      SyncInfo      `json:"sync_info"`
	ValidatorInfo ValidatorInfo `json:"validator_info"`
}

type GenesisValidator struct {
	Address HexBytes `json:"address"`
	PubKey  HexBytes `json:"pub_key"`
	Power   int64    `json:"power"`
	Name    string   `json:"name"`
}

type Genesis struct {
	GenesisTime     time.Time          `json:"genesis_time"`
	ChainID         string             `json:"chain_id"`
	ConsensusParams *ConsensusParams   `json:"consensus_params,omitempty"`
	Validators      []GenesisValidator `json:"validators,omitempty"`
	AppHash         HexBytes           `json:"app_hash"`
	AppState        json.RawMessage    `json:"app_state,omitempty"`
}

type ConsensusParams struct {
	Block     BlockParams     `json:"block"`
	Evidence  EvidenceParams  `json:"evidence"`
	Validator ValidatorParams `json:"validator"`
}

type BlockParams struct {
	MaxBytes   int64 `json:"max_bytes"`
	MaxGas     int64 `json:"max_gas"`
	TimeIotaMs int64 `json:"time_iota_ms"`
}

type EvidenceParams struct {
	MaxAge int64 `json:"max_age"`
}

type ValidatorParams struct {
	PubKeyTypes []string `json:"pub_key_types"`
}

type UnConfirmedTxs struct {
	Count      int            `json:"n_txs"`
	Total      int            `json:"total"`
	TotalBytes int64          `json:"total_bytes"`
	Txs        []TxCommitBody `json:"txs"`
}

//net info
type NetInfo struct {
	Listening bool     `json:"listening"`
	Listeners []string `json:"listeners"`
	NPeers    int      `json:"n_peers"`
	Peers     []Peer   `json:"peers"`
}

type Peer struct {
	NodeInfo         NodeInfo         `json:"node_info"`
	IsOutbound       bool             `json:"is_outbound"`
	ConnectionStatus ConnectionStatus `json:"connection_status"`
	RemoteIP         string           `json:"remote_ip"`
}

type ConnectionStatus struct {
	Duration    time.Duration
	SendMonitor Status
	RecvMonitor Status
	Channels    []ChannelStatus
}

type Status struct {
	Active   bool          // Flag indicating an active transfer
	Start    time.Time     // Transfer start time
	Duration time.Duration // Time period covered by the statistics
	Idle     time.Duration // Time since the last transfer of at least 1 byte
	Bytes    int64         // Total number of bytes transferred
	Samples  int64         // Total number of samples taken
	InstRate int64         // Instantaneous transfer rate
	CurRate  int64         // Current transfer rate (EMA of InstRate)
	AvgRate  int64         // Average transfer rate (Bytes / Duration)
	PeakRate int64         // Maximum instantaneous transfer rate
	BytesRem int64         // Number of bytes remaining in the transfer
	TimeRem  time.Duration // Estimated time to completion
	Progress utils.Percent // Overall transfer progress
}

type ChannelStatus struct {
	ID                byte
	SendQueueCapacity int
	SendQueueSize     int
	Priority          int
	RecentlySent      int64
}

//bench mark
type BenchMarkRequest struct {
	TxNums       int    `json:"tx_nums"`
	TxSendPerSec int    `json:"tx_send_per_sec"`
	Connections  int    `json:"connections"`
	Mode         string `json:"mode"`
	Cmd          string `json:"cmd"`
}

type BenchMarkResponse struct {
	Latency *BenchMarkDetail
	Tps     *BenchMarkDetail
}

type BenchMarkDetail struct {
	Avg   string
	Max   string
	Stdev string
}

//badger
type OperationKeyLog struct {
	Key          string         `json:"key"`
	OperationLog []OperationLog `json:"command_log"`
}

type OperationLog struct {
	Operation string `json:"operation"`
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Time      string `json:"time"`
	Height    string `json:"height"`
	Sequence  string `json:"sequence"`
}
