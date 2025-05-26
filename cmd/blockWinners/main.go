package main

import (
	"flag"
	"fmt"
	"github.com/Snipa22/go-tari-grpc-lib/v2/nodeGRPC"
	"github.com/Snipa22/go-tari-grpc-lib/v2/tari_generated"
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
	nodeGRPCPtr := flag.String("base-node-grpc-address", "node-pool.tari.jagtech.io:18102", "Address for the base-node, defaults to Impala's public pool")
	flag.Parse()
	nodeGRPC.InitNodeGRPC(*nodeGRPCPtr)
	tipData, err := nodeGRPC.GetTipInfo()
	if err != nil {
		log.Fatal(err)
	}

	start := uint64(*depthPtr)
	end := tipData.Metadata.BestBlockHeight

	if start != 0 {
		start = end - start
	}

	results := make(map[string][]uint64)
	for {
		blocks, err := nodeGRPC.GetBlockByHeight(makeRange(start, end))
		if err != nil {
			log.Fatal(err)
		}
		var block *tari_generated.Block
		for i := range blocks {
			block = blocks[i]
			if block.Header.Height > start {
				start = block.Header.Height
			}
			outputs := block.Body.GetOutputs()
			if block.Header.Pow.GetPowAlgo() == 0 {
				// RandomX merge mine
				if len(outputs) == 0 {
					results["RXM_unknown_no_output"] = append(results["RXM_unknown_no_output"], block.Header.Height)
					continue
				}
				for _, output := range outputs {
					features := output.GetFeatures()
					if features != nil {
						if features.OutputType != 1 {
							continue
						}
						txExtra := features.GetCoinbaseExtra()
						if txExtra != nil {
							if strings.Contains(string(txExtra), "WUFJagtech") {
								results["RXM_"+string(txExtra[0:12])] = append(results["RXM_"+string(txExtra[0:12])], block.Header.Height)
								break
							}
							poolID := strings.Map(func(r rune) rune {
								if unicode.IsPrint(r) {
									return r
								}
								return -1
							}, "RXM_"+string(txExtra))
							results[poolID] = append(results[poolID], block.Header.Height)
							break
						} else {
							results["RXM_unknown_no_tx_extra"] = append(results["RXM_unknown_no_tx_extra"], block.Header.Height)
							continue
						}
					} else {
						results["RXM_unknown_no_features"] = append(results["RXM_unknown_no_features"], block.Header.Height)
						continue
					}
				}
				continue
			} else if block.Header.Pow.GetPowAlgo() == 2 {
				// RandomX Tari
				if len(outputs) == 0 {
					results["RXT_unknown_no_output"] = append(results["RXT_unknown_no_output"], block.Header.Height)
					continue
				}
				for _, output := range outputs {
					features := output.GetFeatures()
					if features != nil {
						if features.OutputType != 1 {
							continue
						}
						txExtra := features.GetCoinbaseExtra()
						if txExtra != nil {
							if strings.Contains(string(txExtra), "WUFJagtech") {
								results["RXT_"+string(txExtra[0:12])] = append(results["RXT_"+string(txExtra[0:12])], block.Header.Height)
								break
							}
							poolID := strings.Map(func(r rune) rune {
								if unicode.IsPrint(r) {
									return r
								}
								return -1
							}, "RXT_"+string(txExtra))
							results[poolID] = append(results[poolID], block.Header.Height)
							break
						} else {
							results["RXT_unknown_no_tx_extra"] = append(results["RXT_unknown_no_tx_extra"], block.Header.Height)
							continue
						}
					} else {
						results["RXT_unknown_no_features"] = append(results["RXT_unknown_no_features"], block.Header.Height)
						continue
					}
				}
				continue
			} else {
				// Sha3x is ID 1, but using it as a catch here.
				if len(outputs) == 0 {
					results["SHA3X_unknown_no_output"] = append(results["unknown_no_output"], block.Header.Height)
					continue
				}
				for _, output := range outputs {
					features := output.GetFeatures()
					if features != nil {
						if features.OutputType != 1 {
							continue
						}
						txExtra := features.GetCoinbaseExtra()
						if txExtra != nil {
							poolID := strings.Map(func(r rune) rune {
								if unicode.IsPrint(r) && r < 129 {
									return r
								}
								return -1
							}, "SHA3X_"+string(txExtra[0:12]))
							results[poolID] = append(results[poolID], block.Header.Height)
							break
						} else {
							results["SHA3X_unknown_no_tx_extra"] = append(results["unknown_no_tx_extra"], block.Header.Height)
							continue
						}
					} else {
						results["SHA3X_unknown_no_features"] = append(results["unknown_no_features"], block.Header.Height)
						continue
					}
				}
			}
		}
		if start >= end {
			break
		}
	}
	for pool, blockIds := range results {
		fmt.Println(pool, "has", len(blockIds), "blocks:", blockIds)
	}
}
