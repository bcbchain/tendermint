package rpctest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bcbchain/bclib/tendermint/tmlibs/log"

	abci "github.com/bcbchain/bclib/tendermint/abci/types"
	cmn "github.com/bcbchain/bclib/tendermint/tmlibs/common"

	rpcclient "github.com/bcbchain/bclib/rpc/lib/client"
	cfg "github.com/bcbchain/tendermint/config"
	nm "github.com/bcbchain/tendermint/node"
	"github.com/bcbchain/tendermint/proxy"
	ctypes "github.com/bcbchain/tendermint/rpc/core/types"
	"github.com/bcbchain/tendermint/rpc/grpc"
	pvm "github.com/bcbchain/tendermint/types/priv_validator"
)

var globalConfig *cfg.Config

func waitForRPC() {
	laddr := GetConfig().RPC.ListenAddress
	client := rpcclient.NewJSONRPCClient(laddr)
	ctypes.RegisterAmino(client.Codec())
	result := new(ctypes.ResultStatus)
	for {
		_, err := client.Call("status", map[string]interface{}{}, result)
		if err == nil {
			return
		} else {
			fmt.Println("error", err)
			time.Sleep(time.Millisecond)
		}
	}
}

func waitForGRPC() {
	client := GetGRPCClient()
	for {
		_, err := client.Ping(context.Background(), &core_grpc.RequestPing{})
		if err == nil {
			return
		}
	}
}

// f**ing long, but unique for each test
func makePathname() string {
	// get path
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// fmt.Println(p)
	sep := string(filepath.Separator)
	return strings.Replace(p, sep, "_", -1)
}

func randPort() int {
	return int(cmn.RandUint16()/2 + 10000)
}

func makeAddrs() (string, string, string) {
	start := randPort()
	return fmt.Sprintf("tcp://0.0.0.0:%d", start),
		fmt.Sprintf("tcp://0.0.0.0:%d", start+1),
		fmt.Sprintf("tcp://0.0.0.0:%d", start+2)
}

// GetConfig returns a config for the test cases as a singleton
func GetConfig() *cfg.Config {
	if globalConfig == nil {
		pathname := makePathname()
		globalConfig = cfg.ResetTestRoot(pathname)

		// and we use random ports to run in parallel
		tm, rpc, grpc := makeAddrs()
		globalConfig.P2P.ListenAddress = tm
		globalConfig.RPC.ListenAddress = rpc
		globalConfig.RPC.GRPCListenAddress = grpc
		globalConfig.TxIndex.IndexTags = "app.creator" // see kvstore application
	}
	return globalConfig
}

func GetGRPCClient() core_grpc.BroadcastAPIClient {
	grpcAddr := globalConfig.RPC.GRPCListenAddress
	return core_grpc.StartGRPCClient(grpcAddr)
}

// StartTendermint starts a test tendermint server in a go routine and returns when it is initialized
func StartTendermint(app abci.Application) *nm.Node {
	node := NewTendermint(app)
	err := node.Start()
	if err != nil {
		panic(err)
	}

	// wait for rpc
	waitForRPC()
	waitForGRPC()

	fmt.Println("Tendermint running!")

	return node
}

// NewTendermint creates a new tendermint server and sleeps forever
func NewTendermint(app abci.Application) *nm.Node {
	// Create & start node
	config := GetConfig()
	logger := log.NewTMLogger("/tmp/111", "")
	//logger = log.NewFilter(logger, log.AllowError())
	pvFile := config.PrivValidatorFile()
	pv := pvm.LoadOrGenFilePV(pvFile)
	papp := proxy.NewLocalClientCreator(app)
	node, err := nm.NewNode(config, pv, papp,
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider, logger)
	if err != nil {
		panic(err)
	}
	return node
}
