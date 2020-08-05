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
}


func main() {
	var m TxReq

	tx := []byte(`
		{
			"user_id"       : "abc",
			"file_hash"     : "1234",
			"old_file_hash" : "5678",
			"filename"      : "xxxx",
			"reader_id"     : "xyz",
			"action"        : 1,
			"others"        : "not show"
		}
	`)
	err := json.Unmarshal(tx, &m)
	if err != nil {
		panic(err)
	}

	fmt.Println(m) 

}
