package v1

import (
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v3/cmd/statsExporter/subsystems/blockTemplateCache"
	"github.com/Snipa22/go-tari-grpc-lib/v3/cmd/statsExporter/subsystems/tipDataCache"
	"github.com/Snipa22/go-tari-grpc-lib/v3/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v3/tari_generated"
	"github.com/gin-gonic/gin"
)

type NetworkStats struct {
	BestBlockHeight    uint64 `json:"best_block_height"`
	BestBlockHash      string `json:"best_block_hash"`
	RXMDiff            uint64 `json:"rxm_diff"`
	RXTDiff            uint64 `json:"rxt_diff"`
	Sha3XDiff          uint64 `json:"sha3x_diff"`
	CurBlockReward     uint64 `json:"cur_block_reward"`
	CurBlockRootReward uint64 `json:"cur_block_root_reward"`
}

func GetNetworkStats(c *gin.Context) {
	tipData := tipDataCache.GetTipData()
	returnStruct := NetworkStats{
		BestBlockHeight: tipData.Metadata.BestBlockHeight,
		BestBlockHash:   fmt.Sprintf("%x", tipData.Metadata.BestBlockHash),
	}
	// Get the root reward
	if netState, _ := nodeGRPC.GetNetworkState(); netState != nil {
		returnStruct.CurBlockRootReward = netState.Reward
	}
	// Get chain diffs
	if rxmBT := blockTemplateCache.GetBlockTemplateCache(tari_generated.PowAlgo_POW_ALGOS_RANDOMXM); rxmBT != nil {
		returnStruct.RXMDiff = rxmBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = rxmBT.MinerData.Reward
	}
	if rxtBT := blockTemplateCache.GetBlockTemplateCache(tari_generated.PowAlgo_POW_ALGOS_RANDOMXT); rxtBT != nil {
		returnStruct.RXTDiff = rxtBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = rxtBT.MinerData.Reward
	}
	if sha3xBT := blockTemplateCache.GetBlockTemplateCache(tari_generated.PowAlgo_POW_ALGOS_SHA3X); sha3xBT != nil {
		returnStruct.Sha3XDiff = sha3xBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = sha3xBT.MinerData.Reward
	}
	c.JSON(200, returnStruct)
}
