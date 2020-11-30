package code

const (
	CodeTypeOK = iota
	CodeTypeInvalidSign
	CodeTypeInvalidValidator
	CodeTypeEncodingError
	CodeTypecyberveinQueryError
	CodeTypecyberveinExecuteError
	CodeTypeGetChainInfoError
	CodeTypeInvalidExecuteMode
	CodeTypeInvalidTx
	CodeTypePermissionDenied
	CodeTypeDBPasswordIncorrectError
	CodeTypeTokenTimeoutError
	CodeTypeTokenInvalidError
	CodeTypeIncorrectPassword
	CodeTypeInternalError
)