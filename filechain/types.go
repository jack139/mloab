package filechain


import (
	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/version"
	dbm "github.com/tendermint/tm-db"
)


// 保存应用状态使用
type State struct {
	db      dbm.DB
	Size    int64  `json:"size"`
	Height  int64  `json:"height"`
	AppHash []byte `json:"app_hash"`
}


// 文件信息 - 用户文件表使用
type FileData struct {
	FileName    string `json:"filename"` // 文件名，可为空
	Modified    bool `json:"modified"` // 文件是否已修改
}


// 应用的结构
type App struct {
	types.BaseApplication

	state State
	RetainBlocks int64 // blocks to retain after commit (via ResponseCommit.RetainHeight)
}


// 交易请求数据
type TxReq struct {
	UserId      string `json:"user_id"` // 文件主的用户id
	FileHash    string `json:"file_hash"` // 文件hash
	OldFileHash string `json:"old_file_hash"` // 旧文件的hash (如果 action==修改 需提供)
	FileName    string `json:"filename"` // 文件名，可为空
	ReaderId    string `json:"reader_id"` // 浏览文件的用户id （如果 action==浏览 需提供）
	Action      byte   `json:"action"` // 0x01 文件建立， 0x02 文件浏览， 0x03 文件修改， 0x04 文件删除
}

// 查询请求数据
type QueryReq struct {
	UserId      string `json:"user_id"` // 文件主的用户id，action==2时提供
	FileHash    string `json:"file_hash"` // 文件hash，action==1时提供
	Action      byte   `json:"action"` // 0x01 查询文件历史, 0x02 查询用户的文件列表
}


// 一些参数
var (
	stateKey        = []byte("stateKey")
	fileLinkPrefixKey = []byte("fileLink:")
	blockLinkPrefixKey = []byte("blockLink:")
	userFilePrefixKey = []byte("userFile:")

	ProtocolVersion version.Protocol = 0x1
)
