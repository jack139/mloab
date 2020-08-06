package filechain

/*
	在区块链上构建链表，实现链上搜索

	区块链主要定义
*/



import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"
)


// 一些参数
var (
	stateKey        = []byte("stateKey")
	fileLinkPrefixKey = []byte("fileLink:")
	blockLinkPrefixKey = []byte("blockLink:")
	userFilePrefixKey = []byte("userFile:")

	ProtocolVersion version.Protocol = 0x1
)


// 交易请求数据
type TxReq struct {
	UserId      string `json:"user_id"` // 文件主的用户id
	FileHash    string `json:"file_hash"` // 文件hash
	OldFileHash string `json:"old_file_hash"` // 旧文件的hash (如果 action==修改 需提供)
	FileName    string `json:"filename"` // 文件名，可为空
	ReaderId    string `json:"reader_id"` // 浏览文件的用户id （如果 action==浏览 需提供）
	Action      byte   `json:"action"` // 0x01 文件建立， 0x02 文件浏览， 0x03 文件修改， 0x04 文件删除
}

// 链上交易数据
type TxData struct {
	ReqData TxReq
	LastHeight int64  // 链表指针 (height)
	Created time.Time // 交易数据建立时间
}


type App struct {
	types.BaseApplication

	state State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)
}


func NewApp(rootDir string) *App {
	state := loadState(InitDB(rootDir))
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
}


