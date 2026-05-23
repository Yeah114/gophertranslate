package v1v21v110

import (
	_ "embed"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
)

var (
	//go:embed block_states.nbt
	blockStatesBytes []byte
	blockStates      []define.BlockState
)

func init() {
	blockStates = utils.DecodeBlockStates(blockStatesBytes)
}

func NewBlockRuntimeIDTable(useNetworkIDHashes bool) *block.BlockRuntimeIDTable {
	return block.NewBlockRuntimeIDTableFromStates(blockStates, useNetworkIDHashes)
}
