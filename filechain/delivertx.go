package filechain

/*
	交易上链处理
*/



import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/abci/types"
	//dbm "github.com/tendermint/tm-db"
)


/*
	type TxReq struct {
		UserId      string
		FileHash    string
		OldFileHash string
		FileName    string
		ReaderId    string
		Action      byte  
	}
*/
func (app *App) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	fmt.Println("DeliverTx()")

	db := app.state.db
	var m TxReq

	err := json.Unmarshal(req.Tx, &m)
	if err != nil {
		rsp.Log = "bad json format"
		rsp.Code = 1
		return
	}

	switch m.Action {
	case 0x01: // 新建文件
		rsp.Log = "new file"

		// 生成文件key，添加到db
		fileLinkKey := filePrefixKey(m.FileHash)

		// 生成用户key
		userFileKey := userPrefixKey(m.UserId, m.FileHash)

		// 生成file_data
		fileData := FileData{
			FileName: "",
			Modified: false,
		}
		fileBytes, err := json.Marshal(fileData)
		if err != nil {
			panic(err)
		}

		// 添加到 db
		AddKV(db, fileLinkKey, Int64ToByteArray(app.state.Height+1)) 
		AddKV(db, userFileKey, fileBytes)

		txData := TxData{
			ReqData : m,
			LastHeight : 0,
			Created : time.Now(),
		}
		txDataBytes, err := json.Marshal(txData)
		if err != nil {
			panic(err)
		}

		rsp.Data = txDataBytes

	case 0x02: // 浏览文件
		rsp.Log = "view file"
	case 0x03: // 修改文件
		rsp.Log = "modify file"
	//case 0x04: // 删除文件
	//	rsp.Log = "remove file"
	default:
		rsp.Log = "weird command"
		rsp.Code = 2
	}

	fmt.Println("=================>", rsp.Log)

	return
}

