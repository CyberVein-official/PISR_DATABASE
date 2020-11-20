package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"cybervein.org/CyberveinDB/logger"
	"math"
	"strconv"
)

func StructToJson(v interface{}) []byte {
	json, err := json.Marshal(v)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	return json
}

func JsonToStruct(j []byte, s interface{}) {
	err := json.Unmarshal(j, s)
	if err != nil {
		logger.Log.Error(err)
		return
	}
}

func ByteToHex(b []byte) string {
	return fmt.Sprintf("%X", b)
}

func HexToByte(s string) []byte {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		logger.Log.Error(err)
		return nil
	}
	return bytes
}

func SignToHex(b []byte) string {
	return fmt.Sprintf("%X", b)
}

//percent
type Percent uint32

func round(x float64) int64 {
	if _, frac := math.Modf(x); frac >= 0.5 {
		return int64(math.Ceil(x))
	}
	return int64(math.Floor(x))
}

func percentOf(x, total float64) Percent {
	if x < 0 || total <= 0 {
		return 0
	} else if p := round(x / total * 1e5); p <= math.MaxUint32 {
		return Percent(p)
	}
	return Percent(math.MaxUint32)
}

func (p Percent) Float() float64 {
	return float64(p) * 1e-3
}

func (p Percent) String() string {
	var buf [12]byte
	b := strconv.AppendUint(buf[:0], uint64(p)/1000, 10)
	n := len(b)
	b = strconv.AppendUint(b, 1000+uint64(p)%1000, 10)
	b[n] = '.'
	return string(append(b, '%'))
}
