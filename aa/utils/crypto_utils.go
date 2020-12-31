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
package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"cybervein.org/CyberveinDB/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"path/filepath"
	"time"
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

//#################### JWT ####################

type NodeKeySignMethod struct {
}

func (n NodeKeySignMethod) Verify(signingString, signature string, key interface{}) error {
	k := key.(crypto.PrivKey)
	if !k.PubKey().VerifyBytes([]byte(signingString), HexToByte(signature)) {
		err := fmt.Errorf("JWT Signature invalid")
		logger.Log.Error(err)
		return err
	}
	return nil
}
func (n NodeKeySignMethod) Sign(signingString string, key interface{}) (string, error) {
	k := key.(crypto.PrivKey)
	bytes, err := k.Sign([]byte(signingString))
	if err != nil {
		logger.Log.Error(err)
		return "", err
	}
	return SignToHex(bytes), nil
}
func (n NodeKeySignMethod) Alg() string {
	return "NodeKeySignMethod"
}

type Claims struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func InitJWTMethod() {
	jwt.RegisterSigningMethod("NodeKeySignMethod", func() jwt.SigningMethod {
		return &NodeKeySignMethod{}
	})
}

// GenerateToken generate tokens used for auth
func GenerateToken(address string, username string, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		EncodeMD5(address),
		EncodeMD5(username),
		EncodeMD5(password),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    Config.App.Name,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.GetSigningMethod("NodeKeySignMethod"), claims)
	token, err := tokenClaims.SignedString(NodeKey.PrivKey)
	if err != nil {
		logger.Log.Errorf("GenerateToken error : ", err)
		return "", err
	}
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return NodeKey.PrivKey, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}
