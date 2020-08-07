## leveldb 逻辑分表

| 前缀       | key             | value     |
| ---------- | --------------- | --------- |
| blockLink: | 区块高度        | 区块高度  |
| fileLink:  | 文件hash        | 区块高度  |
| userFile:  | 用户id:文件hash | file_data |

> 说明：

1. Key中数据段用冒号":"分隔
3. file_data定义：

```json
{
	"filename": "abc.docx", // 文件名
	"is_modified": false, // 是否有修改，如果已修改说明不是最新文件hash
}
```

