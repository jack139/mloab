package main

/*
	cleveldb 性能测试

	带 cleveldb 编译
	go build -tags 'cleveldb' levstress.go

	2,000,000 key-value的 FindKey() 耗时
	time elapsed:  166.251µs
	time elapsed:  88.108µs
	time elapsed:  132.681µs
	time elapsed:  87.427µs
	time elapsed:  89.507µs
	time elapsed:  89.192µs
	time elapsed:  117.98µs
	time elapsed:  81.904µs

	10,000,000 key-value的 FindKey() 耗时
	time elapsed:  78.628µs
	time elapsed:  177.877µs
	time elapsed:  179.734µs
	time elapsed:  115.636µs
	time elapsed:  136.898µs
	time elapsed:  95.616µs
	time elapsed:  97.119µs
	time elapsed:  149.53µs

*/

import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"time"
	"math/rand"

	dbm "github.com/tendermint/tm-db"
)

func init() {
    rand.Seed(time.Now().UnixNano())
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
    }
    return string(b)
}

func GenKeys(db dbm.DB, n int) int {
	for i:=0;i<n;i++ {
		key := []byte(RandString(64))
		value := []byte(RandString(64))

		// 存储数据
		db.Set(key, value)
	}

	return 0
}

func CountKeys(db dbm.DB, show int) int {
	// 循环获取
	itr, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	count := 0
	for ; itr.Valid(); itr.Next() {
		if show==1 {
			fmt.Println(string(itr.Key()), "=", string(itr.Value()))
		}
		count += 1
	}

	return count	
}

func FindKey(db dbm.DB, key []byte) []byte {
	// 查询数据
	hasKey, err := db.Has(key)
	if err != nil {
		panic(err)
	}
	if !hasKey {
		return []byte("not found")
	}

	// 获取数据
	value2, err := db.Get(key)
	if err != nil {
		panic(err)
	}

	return value2
}

func AddKV(db dbm.DB, key []byte, value []byte) int {
	db.Set(key, value)
	return 0
}

func main() {
	var db dbm.DB
	name := "mloab"
	dbDir := "n1/data"

	// 初始化数据库
	//db, err := dbm.NewGoLevelDB(name, dbDir) 
	db, err := dbm.NewCLevelDB(name, dbDir)  
	if err != nil {
		panic(err)
	}

	start := time.Now()

	//GenKeys(db, 10000000)
	//fmt.Println("time elapsed: ", time.Now().Sub(start))

	fmt.Println("count=", CountKeys(db, 1))

	//AddKV(db, []byte("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"), []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"))

	//fmt.Println("key=", string(FindKey(db, []byte("zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"))))

	fmt.Println("time elapsed: ", time.Now().Sub(start))
}