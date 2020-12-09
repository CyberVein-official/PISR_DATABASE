package utils

import (
	"path/filepath"

	"cybervein.org/CyberveinDB/logger"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
)

var NodeKey *p2p.NodeKey
var ValidatorKey *privval.FilePVKey

func InitKey() {
	InitNodeKey()
	InitJWTMethod()
	InitValidatorKey()
}

func InitNodeKey() {
	key, err := p2p.LoadNodeKey("../chain/config/node_key.json")
	if err != nil {
		logger.Log.Error(err)
		return
	}
	NodeKey = key
}

func InitValidatorKey() {
	keyFile := filepath.Join("../chain", "config", "priv_validator_key.json")
	stateFile := filepath.Join("../chain", "data", "priv_validator_state.json")
	fpv := privval.LoadFilePV(keyFile, stateFile)
	ValidatorKey = &fpv.Key
}

//#################### Node key ####################
func GetNodeID() string {
	return string(p2p.PubKeyToID(NodeKey.PubKey()))
}

func NodeSign(msg []byte) []byte {
	bytes, err := NodeKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	return bytes
}

func NodeStringSign(msg []byte) string {
	bytes, err := NodeKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Log.Error(err)
		return ""
	}
	return SignToHex(bytes)
}

//#################### Validator key ####################
func ValidatorSign(msg []byte) []byte {
	bytes, err := ValidatorKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	return bytes
}

func ValidatorStringSign(msg []byte) string {
	bytes, err := ValidatorKey.PrivKey.Sign(msg)
	if err != nil {
		logger.Log.Error(err)
		return ""
	}
	return SignToHex(bytes)
}
