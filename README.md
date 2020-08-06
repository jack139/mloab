## MLOAB - Multiple Linked-list On A Blockchain
在区块链上构建链表，实现链上搜索



### 1. 本地多节点测试

编译
```shell
make test
```


初始化

```shell
./test init --home n1
./test init --home n2
```

复制创世块
```shell
cp n1/config/genesis.json n2/config/
```

获取n1节点id
```shell
./test show_node_id --home n1
```

修改n2/config/config.toml
```toml
proxy_app = "tcp://127.0.0.1:36658"
laddr = "tcp://127.0.0.1:36657"
laddr = "tcp://0.0.0.0:36656"
persistent_peers = "b2c82964b2c67236f94a84aa19b0fda6e91869a0@127.0.0.1:26656"
```

启动节点
```shell
./test node --home n1
./test node --home n2
```

提交交易
```shell
curl localhost:26657/broadcast_tx_commit?tx=0x0101
curl localhost:36657/broadcast_tx_commit?tx=0x0101
```

查询交易
```shell
curl localhost:26657/tx?hash=0x...
```

查询验证节点信息
```shell
curl localhost:26657/validators
```

查询网络信息
```shell
curl localhost:26657/net_info
```



### 2. 文件操作日志数据上链

编译
```shell
make build
```

启动

```shell
./mloab init --home 1
./mloab node --home 1
```

tx提交的json格式

```json
{
	"user_id": "abc",  // 文件主的用户id
	"file_hash": "...", // 文件hash
	"old_file_hash": "...", // 旧文件的hash (如果 action==修改 需提供)
	"filename": "file.txt", // 文件名，可为空
	"reader_id": "def", // 浏览文件的用户id （如果 action==浏览 需提供）
	"action": 1,  // 0x01 文件建立， 0x02 文件浏览， 0x03 文件修改， 0x04 文件删除
}
```
新建

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"1234\",\"user_id\":\"abc\",\"action\":1}"'
```

浏览

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"1234\",\"user_id\":\"abc\",\"action\":2,\"reader_id\":\"xyz\"}"'
```

修改

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"file_hash\":\"5678\",\"user_id\":\"abc\",\"action\":3,\"old_file_hash\":\"1234\"}"'
```