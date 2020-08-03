## MLOAB - Multiple Linked-list On A Blockchain

### 本地多节点测试

初始化
```shell
$./test init --home n1
$./test init --home n2
```

复制创世块
```shell
$cp n1/config/genesis.json n2/config/
```

获取n1节点id
```shell
$./test show_node_id --home n1
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
$./test node --home n1
$./test node --home n2
```

提交交易
```shell
$curl localhost:26657/broadcast_tx_commit?tx=0x0101
$curl localhost:36657/broadcast_tx_commit?tx=0x0101
```

查询交易
```shell
$curl localhost:26657/tx?hash=0x...
```

查询验证节点信息
```shell
$curl localhost:26657/validators
```

查询网络信息
```shell
$curl localhost:26657/net_info
```
