package utils

import (
	"fmt"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	bwo_define "github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/full"
	block_utils "github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
	"github.com/Yeah114/gopherconvert/minecraft/world/item"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// BlockRuntimeIDTableFromGameDataAndVersion creates a block runtime ID table from the given game data and protocol version.
func BlockRuntimeIDTableFromGameDataAndVersion(data minecraft.GameData, version string) (*bwo_block.BlockRuntimeIDTable, error) {
	if data.UseBlockNetworkIDHashes {
		return FullBlockRuntimeIDTableFromGameData(data)
	}
	protocolID, ok := protocol.GetProtocol(version)
	if !ok {
		return nil, fmt.Errorf("BlockRuntimeIDTableFromGameDataAndVersion: no protocol found for version %v", version)
	}
	info, _ := protocol.GetProfile(protocolID)
	customBlocks := append([]protocol.BlockEntry{}, data.CustomBlocks...)
	tableFunc, found := block.Pool[info.ID()]
	if !found {
		return nil, fmt.Errorf("BlockRuntimeIDTableFromGameDataAndVersion: no block runtime ID table found for protocol version %v", info.ID())
	}

	table := tableFunc(data.UseBlockNetworkIDHashes)
	existHashSet := make(map[uint32]struct{})

	for _, blockEntry := range customBlocks {
		blockHash := bwo_block.ComputeBlockHash(blockEntry.Name, blockEntry.Properties)
		if _, exists := existHashSet[blockHash]; exists {
			continue
		}
		existHashSet[blockHash] = struct{}{}

		err := table.RegisterCustomBlock(bwo_define.BlockState{
			Name:       blockEntry.Name,
			Properties: block_utils.CloneProperties(blockEntry.Properties),
			Version:    info.BlockStateVersion(),
		})
		if err != nil {
			return nil, fmt.Errorf("BlockRuntimeIDTableFromGameDataAndVersion: failed to register custom block %s: %w", blockEntry.Name, err)
		}
	}

	return table, nil
}

// FullBlockRuntimeIDTableFromGameData creates a block runtime ID table containing all known downgraded block states.
func FullBlockRuntimeIDTableFromGameData(data minecraft.GameData) (*bwo_block.BlockRuntimeIDTable, error) {
	full.Init()
	customBlocks := append([]protocol.BlockEntry{}, data.CustomBlocks...)
	table := bwo_block.NewBlockRuntimeIDTableFromStates(block_utils.CloneBlockStates(full.BlockStates), data.UseBlockNetworkIDHashes)
	existHashSet := make(map[uint32]struct{})

	for _, state := range full.BlockStates {
		blockHash := bwo_block.ComputeBlockHash(state.Name, state.Properties)
		existHashSet[blockHash] = struct{}{}
	}

	for _, blockEntry := range customBlocks {
		blockHash := bwo_block.ComputeBlockHash(blockEntry.Name, blockEntry.Properties)
		if _, exists := existHashSet[blockHash]; exists {
			continue
		}
		existHashSet[blockHash] = struct{}{}

		err := table.RegisterCustomBlock(bwo_define.BlockState{
			Name:       blockEntry.Name,
			Properties: block_utils.CloneProperties(blockEntry.Properties),
			Version:    protocol.CurrentProfile.BlockStateVersion(),
		})
		if err != nil {
			return nil, fmt.Errorf("FullBlockRuntimeIDTableFromGameData: failed to register custom block %s: %w", blockEntry.Name, err)
		}
	}

	return table, nil
}

// ItemRuntimeIDTableFromGameData creates an item runtime ID table from the given game data.
func ItemRuntimeIDTableFromGameData(data minecraft.GameData) (*item.ItemRuntimeIDTable, error) {
	if len(data.Items) == 0 {
		return nil, fmt.Errorf("ItemRuntimeIDTableFromGameData: no item entries found")
	}
	table := item.NewItemRuntimeIDTable()
	table.RegisterItems(data.Items)
	return table, nil
}
