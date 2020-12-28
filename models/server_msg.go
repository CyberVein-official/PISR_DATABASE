package models

type ExecuteRequest struct {
	Cmd  string `json:"cmd"`
	Mode string `json:"mode"`
}

type TxHashRequest struct {
	Hash string `json:"hash"`
}

type BlockHeightRequest struct {
	Height int `json:"height"`
}

type CommittedTxListRequest struct {
	Begin int `json:"begin_height"`
	End   int `json:"end_height"`
}

type ChainInfoRequest struct {
	Min int `json:"min"`
	Max int `json:"max"`
}


type LoginRequest struct {
	Name     string `json:"db_name"`
	Password string `json:"db_password"`
}

type TokenResponse struct {
	Code    uint32 `json:"code"`
	CodeMsg string `json:"code_info"`
	Token   string `json:"token"`
}
