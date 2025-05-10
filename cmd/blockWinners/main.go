package main

import (
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/nodeGRPC"
	"log"
)

func makeRange(min uint64, max uint64) []uint64 {
	a := make([]uint64, max-min+1)
	for i := range a {
		a[i] = min + uint64(i)
	}
	return a
}

func main() {
	tipData, err := nodeGRPC.GetTipInfo()
	if err != nil {
		log.Fatal(err)
	}
	blocks, err := nodeGRPC.GetBlockByHeight(makeRange(tipData.Metadata.BestBlockHeight-501, tipData.Metadata.BestBlockHeight-1))
	if err != nil {
		log.Fatal(err)
	}
	results := make(map[string][]uint64)
	for i := range blocks {
		block := blocks[i]
		if block.Header.Pow.GetPowAlgo() == 0 {
			results["RandomX"] = append(results["RandomX"], block.Header.Height)
			continue
		}
		outputs := block.Body.GetOutputs()
		if len(outputs) > 0 {
			features := outputs[0].GetFeatures()
			if features != nil {
				txExtra := features.GetCoinbaseExtra()
				if txExtra != nil {
					results[string(txExtra[3:12])] = append(results[string(txExtra[3:12])], block.Header.Height)
				} else {
					results["unknown_no_tx_extra"] = append(results["unknown_no_tx_extra"], block.Header.Height)
					continue
				}
			} else {
				results["unknown_no_features"] = append(results["unknown_no_features"], block.Header.Height)
				continue
			}
		} else {
			results["unknown_no_output"] = append(results["unknown_no_output"], block.Header.Height)
			continue
		}
	}
	for pool, blockIds := range results {
		fmt.Println(pool, "has", len(blockIds), "blocks:", blockIds)
	}
}
