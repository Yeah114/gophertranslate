package full

import (
	_ "embed"
	"sync"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	block_utils "github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
)

var (
	//go:embed full_block_states.nbt
	blockStatesBytes []byte
	// BlockStates contains every unique block state reachable by downgrading the latest block palette.
	BlockStates []define.BlockState
	initOnce    sync.Once
)

func Init() {
	initOnce.Do(func() {
		BlockStates = block_utils.DecodeBlockStates(blockStatesBytes)
	})
}

func NewBlockRuntimeIDTable(useNetworkIDHashes bool) *block.BlockRuntimeIDTable {
	Init()
	return block.NewBlockRuntimeIDTableFromStates(block_utils.CloneBlockStates(BlockStates), useNetworkIDHashes)
}
