package filechain

/*
	交易检查
*/



import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/types"
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

/*
新建
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"1234\",\"user_id\":\"abc\",\"action\":1}"'

浏览
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"1234\",\"user_id\":\"abc\",\"action\":2,\"reader_id\":\"xyz\"}"'

修改
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"5678\",\"user_id\":\"abc\",\"action\":3,\"old_file_hash\":\"1234\"}"'
*/

// 检查参数
func (app *App) isValid(tx []byte) error {
	db := app.state.db
	var m TxReq

	err := json.Unmarshal(tx, &m)
	if err != nil {
		return err // json 格式问题
	}

	// 检查参数
	if len(m.UserId)==0 || len(m.FileHash)==0 || m.Action==0 { 
		return fmt.Errorf("bad parameters") // 参数问题
	}

	switch m.Action {
	case 0x01: // 新建文件
		if FileHashExisted(db, m.FileHash) {
			return fmt.Errorf("file_id existed")
		}
	case 0x02: // 浏览文件
		if len(m.ReaderId)==0 {
			return fmt.Errorf("reader_id needed")
		}
	case 0x03: // 修改文件
		if len(m.OldFileHash)==0 {
			return fmt.Errorf("old_file_id needed")
		}
		if !FileHashExisted(db, m.OldFileHash) {
			return fmt.Errorf("old_file_id not existed")
		}
		if FileHashExisted(db, m.FileHash) {
			return fmt.Errorf("new file_id existed")
		}
	//case 0x04: // 删除文件
	//	rsp.Log = "remove file"
	default:
		return fmt.Errorf("weird command")
	}

	return nil
}

func (app *App) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	fmt.Println("CheckTx()")

	fmt.Println(string(req.Tx))

	err := app.isValid(req.Tx)
	if err!=nil {
		rsp.Log = err.Error()
		rsp.Code = 1
	} else {
		rsp.GasWanted = 1
	}

	return 
}


