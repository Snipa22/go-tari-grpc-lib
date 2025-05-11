package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/nodeGRPC"
	"log"
	"strings"
	"unicode"
)

func makeRange(min uint64, max uint64) []uint64 {
	a := make([]uint64, max-min+1)
	for i := range a {
		a[i] = min + uint64(i)
	}
	return a
}

func main() {
	depthPtr := flag.Int("depth", 100, "an int")
	tipData, err := nodeGRPC.GetTipInfo()
	if err != nil {
		log.Fatal(err)
	}
	blocks, err := nodeGRPC.GetBlockByHeight(makeRange(tipData.Metadata.BestBlockHeight-uint64(*depthPtr), tipData.Metadata.BestBlockHeight))
	if err != nil {
		log.Fatal(err)
	}
	results := make(map[string][]uint64)
	for i := range blocks {
		block := blocks[i]
		if block.Header.Pow.GetPowAlgo() == 0 {
			outputs := block.Body.GetOutputs()
			if len(outputs) > 0 {
				features := outputs[0].GetFeatures()
				if features != nil {
					txExtra := features.GetCoinbaseExtra()
					if txExtra != nil {
						poolID := strings.Map(func(r rune) rune {
							if unicode.IsPrint(r) {
								return r
							}
							return -1
						}, "RandomX_"+string(txExtra))
						results[poolID] = append(results[poolID], block.Header.Height)
					} else {
						results["RandomX_unknown_no_tx_extra"] = append(results["RandomX_unknown_no_tx_extra"], block.Header.Height)
						continue
					}
				} else {
					results["RandomX_unknown_no_features"] = append(results["RandomX_unknown_no_features"], block.Header.Height)
					continue
				}
			} else {
				results["RandomX_unknown_no_output"] = append(results["RandomX_unknown_no_output"], block.Header.Height)
				continue
			}
			continue
		}
		outputs := block.Body.GetOutputs()
		if len(outputs) > 0 {
			features := outputs[0].GetFeatures()
			if features != nil {
				txExtra := features.GetCoinbaseExtra()
				if txExtra != nil {
					poolID := strings.Map(func(r rune) rune {
						if unicode.IsPrint(r) && r < 129 {
							return r
						}
						return -1
					}, "SHA3X_"+string(txExtra[0:12]))

					results[poolID] = append(results[poolID], block.Header.Height)
				} else {
					results["SHA3X_unknown_no_tx_extra"] = append(results["unknown_no_tx_extra"], block.Header.Height)
					continue
				}
			} else {
				results["SHA3X_unknown_no_features"] = append(results["unknown_no_features"], block.Header.Height)
				continue
			}
		} else {
			results["SHA3X_unknown_no_output"] = append(results["unknown_no_output"], block.Header.Height)
			continue
		}
	}
	for pool, blockIds := range results {
		fmt.Println(pool, "has", len(blockIds), "blocks:", blockIds)
	}
}
