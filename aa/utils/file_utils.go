package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	DbTxLogFilePath = "../log/db_transactions"
	AppLogFilePath  = "../log/application.log"
)

func InitFiles() {
	os.Remove(DbTxLogFilePath)
}

//command | addr | sign | seq | height | time
func AppendToDBLogFile(lines []string) error {
	dbTxLogFile, err := os.OpenFile(DbTxLogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer dbTxLogFile.Close()

	w := bufio.NewWriter(dbTxLogFile)
	for _, v := range lines {
		w.WriteString(v + "\n")
	}
	return w.Flush()
}

//command | addr | sign | seq | height | time
func ReadTxFromDBLogFile(beginHeight int, endHeight int) ([]string, error) {
	f, err := os.Open(DbTxLogFilePath)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	strList := make([]string, 0)
	for {
		line, err := buf.ReadString('\n')
		split := strings.Split(line, " | ")
		if len(split) < 6 {
			break
		}
		height, err := strconv.Atoi(split[4])
		if err != nil {
			return nil, err
		}
		if height > endHeight {
			break
		} else if height < beginHeight {
			continue
		}
		strList = append(strList, split[0])
	}
	return strList, nil
}

func ReadAll(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func ReadPIDFile(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	br := bufio.NewReader(f)

	defer f.Close()

	line, _, err := br.ReadLine()

	if err != nil {
		return ""
	}
	return string(line)
}
 
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return false
	}
}

func DeleteFile(path string) {
	if Exists(path) {
		os.RemoveAll(path)
	}
}

func CreateDir(path string) {
	os.MkdirAll(path, os.ModePerm)
}

//existing file will be overwritten
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
