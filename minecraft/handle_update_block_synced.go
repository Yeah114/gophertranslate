package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertUpdateBlockSynced converts an UpdateBlockSynced packet.
func (c *MinecraftConverter) ConvertUpdateBlockSynced(pk *packet.UpdateBlockSynced) (*packet.UpdateBlockSynced, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlockSynced: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return &dst, nil
}
