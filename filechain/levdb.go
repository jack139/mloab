package filechain

/*
	cleveldb 相关操作
*/


import (
	"encoding/json"
	"fmt"
	"path"

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


// 文件哈希表前缀
func filePrefixKey(fileHash []byte) []byte {
	return append(fileLinkPrefixKey, fileHash...)
}


// 用户文件表前缀
func userPrefixKey(userId, fileHash []byte) []byte {
	tmp := append(fileLinkPrefixKey, userId...)
	tmp = append(tmp, ':')
	return append(tmp, fileHash...)
}


// 初始化/链接db
func InitDB(rootDir string) dbm.DB {
	// 生成数据文件路径, 放在 --home 目录下的 data 下
	dbDir := path.Join(rootDir, "data")
	fmt.Println("mloab.db path: ", dbDir)

	// 初始化数据库
	db, err := dbm.NewCLevelDB("mloab", dbDir)  
	if err != nil {
		panic(err)
	}

	return db
}

// 从db转入应用状态
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

// 保存应用状态
func saveState(state State) {
	stateBytes, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	state.db.Set(stateKey, stateBytes)
}
