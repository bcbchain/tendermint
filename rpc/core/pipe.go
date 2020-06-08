package core

import (
	"strings"
	"time"

	"github.com/bcbchain/bclib/tendermint/go-crypto"
	dbm "github.com/bcbchain/bclib/tendermint/tmlibs/db"
	"github.com/bcbchain/bclib/tendermint/tmlibs/log"
	"github.com/bcbchain/tendermint/consensus"
	"github.com/bcbchain/tendermint/p2p"
	"github.com/bcbchain/tendermint/proxy"
	sm "github.com/bcbchain/tendermint/state"
	"github.com/bcbchain/tendermint/state/txindex"
	"github.com/bcbchain/tendermint/types"
)

var subscribeTimeout = 5 * time.Second

//----------------------------------------------
// These interfaces are used by RPC and must be thread safe

type Consensus interface {
	GetState() sm.State
	GetValidators() (int64, []*types.Validator)
	GetRoundStateJSON() ([]byte, error)
}

type P2P interface {
	Listeners() []p2p.Listener
	Peers() p2p.IPeerSet
	NumPeers() (outbound, inbound, dialig int)
	NodeInfo() p2p.NodeInfo
	IsListening() bool
	DialPeersAsync(p2p.AddrBook, []string, bool) error
}

//----------------------------------------------
// These package level globals come with setters
// that are expected to be called only once, on startup

var (
	// external, thread safe interfaces
	proxyAppQuery proxy.AppConnQuery
	proxyApp      proxy.AppConns

	// interfaces defined in types and above
	stateDB        dbm.DB
	blockStore     types.BlockStore
	mempool        types.Mempool
	evidencePool   types.EvidencePool
	consensusState Consensus
	p2pSwitch      P2P

	// objects
	pubKey           crypto.PubKey
	genDoc           *types.GenesisDoc // cache the genesis structure
	addrBook         p2p.AddrBook
	txIndexer        txindex.TxIndexer
	consensusReactor *consensus.ConsensusReactor
	eventBus         *types.EventBus // thread safe

	logger log.Logger

	privatePeerIDs []string

	completeStarted bool // it's true if application complete started
)

func SetPrivatePeerIDs(pids string) {
	privatePeerIDs = strings.Split(pids, ",")
}

func SetStateDB(db dbm.DB) {
	stateDB = db
}

func SetBlockStore(bs types.BlockStore) {
	blockStore = bs
}

func SetMempool(mem types.Mempool) {
	mempool = mem
}

func SetEvidencePool(evpool types.EvidencePool) {
	evidencePool = evpool
}

func SetConsensusState(cs Consensus) {
	consensusState = cs
}

func SetSwitch(sw P2P) {
	p2pSwitch = sw
}

func SetPubKey(pk crypto.PubKey) {
	pubKey = pk
}

func SetGenesisDoc(doc *types.GenesisDoc) {
	genDoc = doc
}

func SetAddrBook(book p2p.AddrBook) {
	addrBook = book
}

func SetProxyAppQuery(appConn proxy.AppConnQuery) {
	proxyAppQuery = appConn
}

func SetAppConns(appConns proxy.AppConns) {
	proxyApp = appConns
}

func SetTxIndexer(indexer txindex.TxIndexer) {
	txIndexer = indexer
}

func SetConsensusReactor(conR *consensus.ConsensusReactor) {
	consensusReactor = conR
}

func SetLogger(l log.Logger) {
	logger = l
}

func SetEventBus(b *types.EventBus) {
	eventBus = b
}

func SetCompleteStarted(b bool) {
	completeStarted = b
}
