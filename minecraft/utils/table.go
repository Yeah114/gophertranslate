package utils

import (
	"fmt"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	bwo_define "github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gopherconvert/minecraft/world/item"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// BlockRuntimeIDTableFromGameDataAndVersion creates a block runtime ID table from the given game data and protocol version.
func BlockRuntimeIDTableFromGameDataAndVersion(data minecraft.GameData, version string) (*bwo_block.BlockRuntimeIDTable, error) {
	info := protocol.NewInfoByVersion(version)
	tableFunc, found := block.Pool[info.ID()]
	if !found {
		return nil, fmt.Errorf("BlockRuntimeIDTableFromGameDataAndVersion: no block runtime ID table found for protocol version %v", info.ID())
	}

	table := tableFunc(data.UseBlockNetworkIDHashes)
	for _, blockEntry := range data.CustomBlocks {
		err := table.RegisterCustomBlock(bwo_define.BlockState{
			Name:       blockEntry.Name,
			Properties: blockEntry.Properties,
			Version:    info.Version(),
		})
		if err != nil {
			return nil, fmt.Errorf("BlockRuntimeIDTableFromGameDataAndVersion: failed to register custom block %s: %w", blockEntry.Name, err)
		}
	}
	table.FinaliseRegister()

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
