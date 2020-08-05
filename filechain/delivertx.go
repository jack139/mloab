package filechain

/*
	交易上链处理
*/



import (
	"encoding/json"
	"fmt"
	//"time"

	"github.com/tendermint/tendermint/abci/types"
	//dbm "github.com/tendermint/tm-db"
)



func (app *App) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	fmt.Println("DeliverTx()")

	var m TxReq

	err := json.Unmarshal(req.Tx, &m)
	if err != nil {
		rsp.Log = "bad json format"
		return
	}

	fmt.Println("=================>")
	switch m.Action {
	case 0x01: // 新建文件
		rsp.Log = "new file"
	case 0x02: // 浏览文件
		rsp.Log = "view file"
	case 0x03: // 修改文件
		rsp.Log = "modify file"
	case 0x04: // 删除文件
		rsp.Log = "remove file"
	default:
		rsp.Log = "weird command"
	}

	fmt.Println(rsp.Log)

	return
}

