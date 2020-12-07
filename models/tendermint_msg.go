package models

type TxCommitBody struct {
	Data      *TxCommitData `json:"data"`
	Signature string        `json:"signature"`
	Address   string        `json:"address"`
}

type ValidatorUpdateBody struct {
	ValidatorUpdate *ValidatorUpdateData `json:"validator_update"`
	Signature       string               `json:"signature"`
	Address         string               `json:"address"`
	Sequence        string               `json:"sequence"`
}
