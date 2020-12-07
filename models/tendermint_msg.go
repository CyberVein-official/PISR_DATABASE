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

type ValidatorUpdateData struct {
	PublicKey string `json:"public_key"`
	Power     string `json:"power"`
}

type TxCommitData struct {
	Operation string `json:"operation"`
	Sequence  string `json:"sequence"`
}
