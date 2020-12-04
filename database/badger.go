package database

import (
	"fmt"
	"cybervein.org/CyberveinDB/logger"
	"cybervein.org/CyberveinDB/models"
	"cybervein.org/CyberveinDB/utils"
	"github.com/dgraph-io/badger"
	"os"
)

const BadgerPath string = "../badger"

var BadgerDB *badger.DB

func InitBadgerDB() {
	BadgerDB = NewBadgerDB()
}

func NewBadgerDB() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions(BadgerPath))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	return db
}

func UpdateBadgerVal(key []byte, val []byte) error {
	txn := BadgerDB.NewTransaction(true)
	if err := txn.Set(key, val); err == nil {
		_ = txn.Commit()
	} else {
		return err
	}
	return nil
}

func DeleteBadgerKey(key []byte) error {
	txn := BadgerDB.NewTransaction(true)
	if err := txn.Delete(key); err == nil {
		_ = txn.Commit()
	} else {
		return err
	}
	return nil
}

func GetBadgerVal(key []byte) []byte {
	var value []byte
	BadgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		} else {
			return item.Value(func(val []byte) error {
				value = val
				return nil
			})
		}
		return nil
	})
	return value
}

func UpdateSingleKeyWriteLog(op models.OperationLog) {
	key, err := GetKey(op.Operation)
	if err != nil {
		logger.Log.Error(err)
		panic(err)
	}
	val := GetBadgerVal([]byte(key))
	res := &models.OperationKeyLog{
		Key:          key,
		OperationLog: make([]models.OperationLog, 0),
	}
	if val == nil {
		res.OperationLog = append(res.OperationLog, op)
	} else {
		opLog := &models.OperationKeyLog{}
		utils.JsonToStruct(val, &opLog)
		res.OperationLog = append(opLog.OperationLog, op)
	}
	UpdateBadgerVal([]byte(key), utils.StructToJson(res))
}

func UpdateKeyWriteLog(operations []models.OperationLog) {
	for _, op := range operations {
		key, err := GetKey(op.Operation)
		if err != nil {
			logger.Log.Error(err)
			continue
		}
		val := GetBadgerVal([]byte(key))
		res := &models.OperationKeyLog{
			Key:          key,
			OperationLog: make([]models.OperationLog, 0),
		}
		if val == nil {
			res.OperationLog = append(res.OperationLog, op)
		} else {
			opLog := &models.OperationKeyLog{}
			utils.JsonToStruct(val, &opLog)
			res.OperationLog = append(opLog.OperationLog, op)
		}
		UpdateBadgerVal([]byte(key), utils.StructToJson(res))
	}
}

func GetKeyWriteLog(key string) (*models.OperationKeyLog, error) {
	val := GetBadgerVal([]byte(key))
	if val != nil {
		opLog := &models.OperationKeyLog{}
		utils.JsonToStruct(val, &opLog)
		return opLog, nil
	}
	return nil, nil
}
