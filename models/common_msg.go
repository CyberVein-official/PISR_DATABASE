package models

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"cybervein.org/CyberveinDB/models/code"
	"cybervein.org/CyberveinDB/utils"
	"github.com/gin-gonic/gin"
)

type GinMsg struct {
	C *gin.Context
}

type CommonResponse struct {
	Code    uint32      `json:"code"`
	CodeMsg string      `json:"code_info"`
	Data    interface{} `json:"data"`
}

func (g *GinMsg) Response(httpCode int, data interface{}) {
	g.C.JSON(httpCode, data)
	return
}

func (g *GinMsg) SuccessWithData(data interface{}) {
	g.C.JSON(http.StatusOK, &CommonResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.CodeTypeOKMsg,
		Data:    data,
	})
	return
}

func (g *GinMsg) Success() {
	g.C.JSON(http.StatusOK, &CommonResponse{
		Code:    code.CodeTypeOK,
		CodeMsg: code.CodeTypeOKMsg,
	})
	return
}

func (g *GinMsg) Error(httpCode int, code uint32, codeMsg string, err error) {
	g.C.JSON(httpCode, &CommonResponse{
		Code:    code,
		CodeMsg: fmt.Sprintf(codeMsg+" : %s", err),
	})
	return
}

func (g *GinMsg) SimpleErrorMsg(httpCode int, code uint32, codeMsg string) {
	g.C.JSON(httpCode, &CommonResponse{
		Code:    code,
		CodeMsg: codeMsg,
	})
	return
}

func (g *GinMsg) CommonResponse(httpCode int, code uint32, codeMsg string) {
	g.C.JSON(httpCode, &CommonResponse{
		Code:    code,
		CodeMsg: codeMsg,
	})
	return
}

func (g *GinMsg) DecodeRequestBody(data interface{}) {
	body, _ := ioutil.ReadAll(g.C.Request.Body)
	utils.JsonToStruct(body, data)
}

type HexBytes []byte

func (bz HexBytes) Marshal() ([]byte, error) {
	return bz, nil
}

func (bz *HexBytes) Unmarshal(data []byte) error {
	*bz = data
	return nil
}

func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], []byte(s))
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return nil
	}
	bz2, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}

func (bz HexBytes) Bytes() []byte {
	return bz
}

func (bz HexBytes) String() string {
	return strings.ToUpper(hex.EncodeToString(bz))
}
