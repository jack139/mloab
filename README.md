## MLOAB - Multiple Linked-list On A Blockchain
在区块链上构建链表，实现链上搜索。（以文件操作日志为例）



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
	"uid": "abc",  // 文件主的用户id
	"fhash": "...", // 文件hash
	"ofhash": "...", // 旧文件的hash (如果 act==修改 需提供)
	"fn": "file.txt", // 文件名，可为空
	"rid": "def", // 浏览文件的用户id （如果 act==浏览 需提供）
	"act": 1,  // 1 文件建立， 2 文件浏览， 3 文件修改， 4 文件删除
    "nonce" : "...", // 随机字符串（如果能保证tx提交内容不重复，此字段可不用）
}
```
新建

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"1234\",\"uid\":\"abc\",\"act\":1}"'
```

浏览

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"1234\",\"uid\":\"abc\",\"act\":2,\"rid\":\"xyz\",\"nonce\":123}"'
```

修改

```shell
curl -g 'http://localhost:26657/broadcast_tx_commit?tx="{\"fhash\":\"5678\",\"uid\":\"abc\",\"act\":3,\"ofhash\":\"1234\"}"'
```



### 3. 链上数据检索

query提交的json格式

```json
{
    "user_id": "abc", // 文件主的用户id，action==2时提供
    "file_hash": "1234",  // 文件hash，action==1时提供
    "action": 1, // 0x01 查询文件历史, 0x02 查询用户的文件列表
}
```

文件浏览/修改历史

```shell
curl -g 'http://localhost:26657/abci_query?data="{\"fhash\":\"5678\",\"act\":1}"'
```

用户文件列表

```shell
curl -g 'http://localhost:26657/abci_query?data="{\"uid\":\"abc\",\"act\":2}"'
```

