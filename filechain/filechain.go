package filechain

/*
	在区块链上构建链表，实现链上搜索
*/

import (
	"encoding/json"
	"encoding/binary"
	"fmt"
	"time"
	"path"

	"github.com/tendermint/tendermint/abci/types"
	//"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
	//"github.com/spf13/viper"
)


var (
	stateKey        = []byte("stateKey")
	customerPrefixKey = []byte("customerKey:")

	ProtocolVersion version.Protocol = 0x1
)

type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}

func loadState(db dbm.DB) State {
	var state State
	state.db = db
	stateBytes, err := db.Get(stateKey)
	if err != nil {
		panic(err)
	}
	if len(stateBytes) == 0 {
		return state
	}
	err = json.Unmarshal(stateBytes, &state)
	if err != nil {
		panic(err)
	}
	return state
}

func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}

func prefixKey(key []byte) []byte {
	return append(customerPrefixKey, key...)
}


// ------------------------------------------------------


type App struct {
	types.BaseApplication

	state State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)
}

// 交易请求数据
type TxReq struct {
	UserId    string `json:"user_id"` // 文件主的用户id
	FileId    string `json:"file_id"` // 文件id(文件hash)
	OldFileId string `json:"old_file_id"` // 旧文件的id (如果 action==23 需提供)
	ReaderId  string `json:"reader_id"` // 浏览文件的用户id （如果 action==2 需提供）
	Action    byte   `json:"action"` // 0x01 文件建立， 0x02 文件浏览， 0x03 文件修改， 0x04 文件删除
}

// 链上交易数据
type TxData struct {
	ReqData TxReq
	LastHeight int64  // 链表指针 (height)
	Created time.Time // 交易数据建立时间
}

func NewApp(rootDir string) *App {
	// 生成数据文件路径, 放在 --home 目录下的 data 下
	dbDir := path.Join(rootDir, "data")
	fmt.Println("mloab.db path: ", dbDir)

	// 初始化数据库
	db, err := dbm.NewCLevelDB("mloab", dbDir)  
	if err != nil {
		panic(err)
	}

	state := loadState(db)

	return &App{state: state}
}


func (app *App) Info(req types.RequestInfo) (resInfo types.ResponseInfo) {
	return types.ResponseInfo{
		Data:             fmt.Sprintf("{\"size\":%v}", app.state.Size),
		Version:          version.ABCIVersion,
		AppVersion:       ProtocolVersion.Uint64(),
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.AppHash,
	}
}

func (app *App) isValid(tx []byte) (code uint32) {
	var m TxReq

	err := json.Unmarshal(tx, &m)
	if err != nil {
		return 4 // json 格式问题
	}

	// 检查参数
	if len(m.UserId)==0 || len(m.FileId)==0 || m.Action==0 { 
		return 1 // 参数问题
	} else if m.Action==2 && len(m.ReaderId)==0 {
		return 2 // 浏览文件，少reader_id
	} else if m.Action==3 && len(m.OldFileId)==0 {
		return 3 // 修改文件，少 旧文件 id
	}

	return
}

func (app *App) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	fmt.Println("CheckTx()")

	code := app.isValid(req.Tx)
	if code!=0{
		rsp.Log = "bad request data"
	}
	rsp.Code = code
	rsp.GasWanted = 1

	return 
}

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

func (app *App) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	fmt.Println("Query()")

	rsp.Log = fmt.Sprintf("query")
	return
}


func (app *App) Commit() (rsp types.ResponseCommit) {
	fmt.Println("Commit()")

	// Using a db - just return the big endian size of the db
	appHash := make([]byte, 8)
	binary.PutVarint(appHash, app.state.Size)
	app.state.AppHash = appHash
	app.state.Height++
	saveState(app.state)

	resp := types.ResponseCommit{Data: appHash}
	if app.RetainBlocks > 0 && app.state.Height >= app.RetainBlocks {
		resp.RetainHeight = app.state.Height - app.RetainBlocks + 1
	}
	return resp

	// 生成数据文件路径, 放在 --home 目录下的 data 下
	//home := viper.GetString(cli.HomeFlag)
	//stat_path := path.Join(home, "data", "counter.state")

	//fmt.Println("path: ", stat_path)

	// 保存缓存数据
	//bz, err := json.Marshal(app)
	//if err != nil {
	//	panic(err)
	//}
	//ioutil.WriteFile(stat_path, bz, 0644)
	//return
}


