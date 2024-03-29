package main

/*
	在区块链上构建链表，实现链上搜索

	主程序：

	编译：
	go build -tags 'cleveldb' main.go
*/

import (
	"mloab/filechain"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

func main() {
	root := commands.RootCmd
	root.AddCommand(commands.GenNodeKeyCmd)
	root.AddCommand(commands.GenValidatorCmd)
	root.AddCommand(commands.InitFilesCmd)
	root.AddCommand(commands.ResetAllCmd)
	root.AddCommand(commands.ShowNodeIDCmd)
	root.AddCommand(commands.TestnetFilesCmd)

	nodeProvider := makeNodeProvider()
	root.AddCommand(commands.NewRunNodeCmd(nodeProvider))

	exec := cli.PrepareBaseCmd(root, "wiz", ".")
	exec.Execute()
}

func makeNodeProvider() node.Provider {
	return func(config *cfg.Config, logger log.Logger) (*node.Node, error) {
		nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
		if err != nil {
			return nil, err
		}

		// instance app
		app := filechain.NewApp(config.RootDir)

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
