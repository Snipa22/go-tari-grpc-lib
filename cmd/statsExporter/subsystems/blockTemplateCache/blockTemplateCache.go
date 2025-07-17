package blockTemplateCache

import (
	"encoding/binary"
	"github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"math/rand"
	"sync"
)

// poolID is a random byte string used to ID the pool in the coinbase txn
var poolID *[]byte

// PoolStringID is a 9 character string used to ID the pool in the coinbase txn, this is a list of valid ones.
var PoolStringID *[]byte

// Yes, I'm writing a function to turn strings into []bytes

// BlockTemplateCacheStruct contains the response from a GetBlockTemplate, this is used to generate the
// GetNewBlockTemplateWithCoinbasesRequest that is sent to the Tari daemon, this is needed to build the base coinbase
// transaction, which needs to contain the pool ID and other data related to the miner unique nonce.
type BlockTemplateCacheStruct struct {
	BlockTemplate *tari_generated.NewBlockTemplateResponse
	reward        uint64
	mutex         sync.RWMutex
}

var sha3xBTCache *BlockTemplateCacheStruct = nil
var rxtBTCache *BlockTemplateCacheStruct = nil
var rxmBTCache *BlockTemplateCacheStruct = nil
var running = false

// When a miner requests a block template, we need to update the coinbase txn with their unique hash and the pool hash

// UpdateBlockTemplateCache keeps the main block template system spinning and updating to keep things clean
func UpdateBlockTemplateCache(core *milieu.Milieu) {
	if running {
		return
	}
	running = true
	defer func() {
		running = false
	}()
	if poolID == nil {
		buf := make([]byte, 8)
		binary.LittleEndian.PutUint64(buf, rand.Uint64())
		poolID = &buf
	}

	blockTemplateResponse, err := nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_SHA3X})
	if err != nil {
		core.CaptureException(err)
		return
	}
	if sha3xBTCache == nil {
		sha3xBTCache = &BlockTemplateCacheStruct{}
	}

	sha3xBTCache.mutex.Lock()
	sha3xBTCache.BlockTemplate = blockTemplateResponse
	sha3xBTCache.reward = blockTemplateResponse.MinerData.Reward
	sha3xBTCache.mutex.Unlock()

	blockTemplateResponse, err = nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_RANDOMXT})
	if err != nil {
		core.CaptureException(err)
		return
	}
	if rxtBTCache == nil {
		rxtBTCache = &BlockTemplateCacheStruct{}
	}
	rxtBTCache.mutex.Lock()
	rxtBTCache.BlockTemplate = blockTemplateResponse
	rxtBTCache.reward = blockTemplateResponse.MinerData.Reward
	rxtBTCache.mutex.Unlock()

	blockTemplateResponse, err = nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_RANDOMXM})
	if err != nil {
		core.CaptureException(err)
		return
	}
	if rxmBTCache == nil {
		rxmBTCache = &BlockTemplateCacheStruct{}
	}
	rxmBTCache.mutex.Lock()
	rxmBTCache.BlockTemplate = blockTemplateResponse
	rxmBTCache.reward = blockTemplateResponse.MinerData.Reward
	rxmBTCache.mutex.Unlock()
}

func GetBlockTemplateCache(algo tari_generated.PowAlgo_PowAlgos) *tari_generated.NewBlockTemplateResponse {
	switch algo {
	case tari_generated.PowAlgo_POW_ALGOS_RANDOMXM:
		rxmBTCache.mutex.RLock()
		defer rxmBTCache.mutex.RUnlock()
		return rxmBTCache.BlockTemplate
	case tari_generated.PowAlgo_POW_ALGOS_SHA3X:
		sha3xBTCache.mutex.RLock()
		defer sha3xBTCache.mutex.RUnlock()
		return sha3xBTCache.BlockTemplate
	case tari_generated.PowAlgo_POW_ALGOS_RANDOMXT:
		rxtBTCache.mutex.RLock()
		defer rxtBTCache.mutex.RUnlock()
		return rxtBTCache.BlockTemplate
	}
	return nil
}
