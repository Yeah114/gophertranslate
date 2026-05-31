package utils

import (
	"github.com/Yeah114/bedrock-world-operator/define"
	block_utils "github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func BlockStateToBlockEntry(state define.BlockState) protocol.BlockEntry {
	return protocol.BlockEntry{
		Name:       state.Name,
		Properties: block_utils.CloneProperties(state.Properties),
	}
}

func BlockStatesToBlockEntries(states []define.BlockState) []protocol.BlockEntry {
	entries := make([]protocol.BlockEntry, len(states))
	for i, state := range states {
		entries[i] = BlockStateToBlockEntry(state)
	}
	return entries
}
