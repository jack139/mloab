package main

/*
	json 测试

*/

import (
	"encoding/json"
	"fmt"
)


// 交易请求数据
type TxReq struct {
	UserId      string `json:"user_id"`
	FileHash    string `json:"file_hash"`
	OldFileHash string `json:"old_file_hash"` 
	FileName    string `json:"filename"`
	ReaderId    string `json:"reader_id"` 
	Action      byte   `json:"action"`
	BoolT       bool   `json:"bool_test"`
}


func main() {
	var m, n TxReq
	var txs []TxReq

	tx := []byte(`
		{
			"user_id"       : "abc",
			"file_hash"     : "1234",
			"old_file_hash" : "5678",
			"filename"      : "xxxx",
			"reader_id"     : "xyz",
			"action"        : 1,
			"others"        : "not show",
			"bool_test"     : true
		}
	`)
	err := json.Unmarshal(tx, &m)
	if err != nil {
		panic(err)
	}

	txs = append(txs, m)
	txs = append(txs, m)

	fmt.Println(txs) 

}
