package utils

import (
	"fmt"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/block"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// BlockRuntimeIDTableFromGameData creates a block runtime ID table from the provided game data.
// It looks up the correct block runtime table generator for the protocol version
// and registers any custom blocks present in the game data.
func BlockRuntimeIDTableFromGameData(data minecraft.GameData) (*bwo_block.BlockRuntimeIDTable, error) {
	info := protocol.NewInfoByVersion(data.BaseGameVersion)
	tableFunc, found := block.Pool[info.ID()]
	if !found {
		return nil, fmt.Errorf("BlockRuntimeIDTableFromGameData: no block runtime ID table found for protocol version %v", info.ID())
	}

	table := tableFunc(data.UseBlockNetworkIDHashes)
	for _, blockEntry := range data.CustomBlocks {
		err := table.RegisterCustomBlock(define.BlockState{
			Name:       blockEntry.Name,
			Properties: blockEntry.Properties,
			Version:    info.Version(),
		})
		if err != nil {
			return nil, fmt.Errorf("BlockRuntimeIDTableFromGameData: failed to register custom block %s: %w", blockEntry.Name, err)
		}
	}
	table.FinaliseRegister()

	return table, nil
}
