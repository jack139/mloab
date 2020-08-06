package filechain

/*
	链上查询
*/



import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/types"
	rpctypes "github.com/tendermint/tendermint/rpc/jsonrpc/types"
	rpc "github.com/tendermint/tendermint/rpc/core"
)

/*
	type QueryReq struct {
		UserId      string `json:"user_id"` // 文件主的用户id，action==2时提供
		FileHash    string `json:"file_hash"` // 文件hash，action==1时提供
		Action      byte   `json:"action"` // 0x01 查询文件历史, 0x02 查询用户的文件列表
	}
*/

/*
查询文件历史
curl -g 'http://localhost:26657/abci_query?data="{\"file_hash\":\"5678\",\"action\":1}"'

查询用户文件
curl -g 'http://localhost:26657/abci_query?data="{\"user_id\":\"abc\",\"action\":2}"'

测试入口
curl -g 'http://localhost:26657/abci_query?data="{\"action\":3}"'
*/
func (app *App) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	fmt.Println("Query()")

	fmt.Println(string(req.Data))

	db := app.state.db

	var m QueryReq

	err := json.Unmarshal(req.Data, &m)
	if err != nil {
		rsp.Log = "bad json format"
		rsp.Code = 1
		return
	}

	switch m.Action {
	case 0x01: // 文件历史
		rsp.Log = "file history"

		// 文件key, 找到链头
		fileLinkKey := filePrefixKey(m.FileHash)
		height := FindKey(db, fileLinkKey)  // 这里 height 返回是 []byte
		for ;len(height)!=0; {
			fmt.Printf("--> %s", height)

			// 在blcok链上找下一个
			blockLinkKey := blockPrefixKey(ByteArrayToInt64(height))
			height = FindKey(db, blockLinkKey)
		}
		fmt.Println("")

	case 0x02: // 用户文件
		rsp.Log = "user file list"

		start := fmt.Sprintf("%s%s:", userFilePrefixKey, m.UserId)
		end := fmt.Sprintf("%s\xff", start)

		// 循环获取
		itr, err := db.Iterator([]byte(start), []byte(end))
		if err != nil {
			panic(err)
		}

		count := 0
		for ; itr.Valid(); itr.Next() {
			fmt.Println(string(itr.Key()), "=", string(itr.Value()))
			count += 1
		}

	case 0x03: // 测试
		// func Block(ctx *rpctypes.Context, heightPtr *int64) (*ctypes.ResultBlock, error) 
		var height int64
		var ctx rpctypes.Context
		height = 1
		re, err := rpc.Block(&ctx, &height)
		if err!=nil {
			panic(err)
		}
		fmt.Println(re)

	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	return
}

