package v1

import (
	"fmt"
	core "github.com/Snipa22/core-go-lib/milieu"
	"github.com/Snipa22/go-tari-grpc-lib/v2/cmd/statsExporter/support"
	"github.com/Snipa22/go-tari-grpc-lib/v2/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v2/tari_generated"
	"github.com/gin-gonic/gin"
)

type NetworkStats struct {
	BestBlockHeight uint64 `json:"best_block_height"`
	BestBlockHash   string `json:"best_block_hash"`
	RXMDiff         uint64 `json:"rxm_diff"`
	RXTDiff         uint64 `json:"rxt_diff"`
	Sha3XDiff       uint64 `json:"sha3x_diff"`
	CurBlockReward  uint64 `json:"cur_block_reward"`
}

func GetNetworkStats(c *gin.Context) {
	milieu := c.MustGet("MILIEU").(core.Milieu)
	tipData, err := nodeGRPC.GetTipInfo()
	if err != nil {
		c.JSON(500, support.ErrorResponse{Error: "Unable to get tip data"})
		milieu.CaptureException(err)
		return
	}
	returnStruct := NetworkStats{
		BestBlockHeight: tipData.Metadata.BestBlockHeight,
		BestBlockHash:   fmt.Sprintf("%x", tipData.Metadata.BestBlockHash),
	}
	// Get chain diffs
	if rxmBT, _ := nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_RANDOMXM}); rxmBT != nil {
		returnStruct.RXMDiff = rxmBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = rxmBT.MinerData.Reward
	}
	if rxtBT, _ := nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_RANDOMXT}); rxtBT != nil {
		returnStruct.RXMDiff = rxtBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = rxtBT.MinerData.Reward
	}
	if sha3xBT, _ := nodeGRPC.GetBlockTemplate(&tari_generated.PowAlgo{PowAlgo: tari_generated.PowAlgo_POW_ALGOS_SHA3X}); sha3xBT != nil {
		returnStruct.RXMDiff = sha3xBT.MinerData.TargetDifficulty
		returnStruct.CurBlockReward = sha3xBT.MinerData.Reward
	}
	c.JSON(200, returnStruct)
}
