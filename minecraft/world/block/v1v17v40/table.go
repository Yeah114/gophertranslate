package v1v17v40

import (
	_ "embed"
	"sync"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
)

var (
	//go:embed block_states.nbt
	blockStatesBytes []byte
	blockStates      []define.BlockState
	initOnce         sync.Once
)

func Init() {
	initOnce.Do(func() {
		blockStates = utils.DecodeBlockStates(blockStatesBytes)
	})
}

func NewBlockRuntimeIDTable(useNetworkIDHashes bool) *block.BlockRuntimeIDTable {
	Init()
	return block.NewBlockRuntimeIDTableFromStates(utils.CloneBlockStates(blockStates), useNetworkIDHashes)
}
