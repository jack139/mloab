package filechain

/*
	cleveldb 相关操作
*/


import (
	"encoding/json"
	"encoding/binary"
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
func filePrefixKey(fileHash string) []byte {
	return append(fileLinkPrefixKey, []byte(fileHash)...)
}


// 用户文件表前缀
func userPrefixKey(userId, fileHash string) []byte {
	tmp := append(userFilePrefixKey, []byte(userId)...)
	tmp = append(tmp, ':')
	return append(tmp, []byte(fileHash)...)
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


// 给点起始点，获取所有符合条件kv
func SearchKeys(db dbm.DB, start, end []byte) int {
	// 循环获取
	itr, err := db.Iterator(start, end)
	if err != nil {
		panic(err)
	}

	count := 0
	for ; itr.Valid(); itr.Next() {
		fmt.Println(string(itr.Key()), "=", string(itr.Value()))
		count += 1
	}

	return count
}


// 获取数据: 未找到返回 nil
func FindKey(db dbm.DB, key []byte) []byte {
	value2, err := db.Get(key)
	if err != nil {
		panic(err)
	}

	return value2
}


// 添加key 成功返回 nil
func AddKV(db dbm.DB, key []byte, value []byte) error {
	err := db.Set(key, value)
	if err != nil {
		panic(err)
	}

	return nil
}

// 检查文件hash是否已存在
func FileHashExisted(db dbm.DB, fileHash string) bool {
	if FindKey(db, []byte(fileHash))!=nil {
		return true
	}

	return false
}


/*
	int64 <---> []byte 

	i := int64(-123456789)

	fmt.Println(i)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))

	fmt.Println(b)

	i = int64(binary.LittleEndian.Uint64(b))
	fmt.Println(i)
*/
func Int64ToByteArray(a int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(a))
	return b
}

func ByteArrayToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b))
}
