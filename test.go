package main

import (
	"encoding/json"
	"fmt"
	"path"
	"io/ioutil"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
	"github.com/spf13/viper"
)

func main() {
	root := commands.RootCmd
	root.AddCommand(commands.GenNodeKeyCmd)
	root.AddCommand(commands.GenValidatorCmd)
	root.AddCommand(commands.InitFilesCmd)
	root.AddCommand(commands.ResetAllCmd)
	root.AddCommand(commands.ShowNodeIDCmd)
	root.AddCommand(commands.TestnetFilesCmd)

	app := NewApp()
	nodeProvider := makeNodeProvider(app)
	root.AddCommand(commands.NewRunNodeCmd(nodeProvider))

	exec := cli.PrepareBaseCmd(root, "wiz", ".")
	exec.Execute()
}

type App struct {
	types.BaseApplication
	Value int
}

func NewApp() *App {
	return &App{}
}

func (app *App) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	fmt.Println("CheckTx()")

	tx := req.Tx

	if tx[0] == 0x01 || tx[0] == 0x02 || tx[0] == 0x03 {
		rsp.Log = "tx accepted"
		return
	}
	rsp.Code = 1
	rsp.Log = "bad tx rejected"
	return
}

func (app *App) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	fmt.Println("DeliverTx()")

	tx := req.Tx

	fmt.Println("=================>")
	switch tx[0] {
	case 0x01:
		app.Value += 1
	case 0x02:
		app.Value -= 1
	case 0x03:
		app.Value = 0
	default:
		rsp.Log = "weird command"
	}
	return
}

func (app *App) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	fmt.Println("Query()")

	rsp.Log = fmt.Sprintf("counter: %d", app.Value)
	return
}


func (app *App) Commit() (rsp types.ResponseCommit) {
	fmt.Println("Commit()")

	// 生成数据文件路径, 放在 --home 目录下的 data 下
	home := viper.GetString(cli.HomeFlag)
	stat_path := path.Join(home, "data", "counter.state")

	//fmt.Println("path: ", stat_path)

	// 保存缓存数据
	bz, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(stat_path, bz, 0644)
	return
}


func makeNodeProvider(app types.Application) node.Provider {
	return func(config *cfg.Config, logger log.Logger) (*node.Node, error) {
		nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
		if err != nil {
			return nil, err
		}

		// read private validator
		pv := privval.LoadFilePV(
			config.PrivValidatorKeyFile(),
			config.PrivValidatorStateFile(),
		)

		return node.NewNode(config,
			pv, //privval.LoadOrGenFilePV(config.PrivValidatorFile()),
			nodeKey,
			proxy.NewLocalClientCreator(app),
			node.DefaultGenesisDocProviderFunc(config),
			node.DefaultDBProvider,
			node.DefaultMetricsProvider(config.Instrumentation),
			logger,
		)
	}
}
