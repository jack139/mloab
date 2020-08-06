package filechain

/*
	链上查询
*/



import (
	//"encoding/json"
	"fmt"
	//"time"

	"github.com/tendermint/tendermint/abci/types"
	//dbm "github.com/tendermint/tm-db"
)


func (app *App) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	fmt.Println("Query()")

	rsp.Log = fmt.Sprintf("query")
	return
}

