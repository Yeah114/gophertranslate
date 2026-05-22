package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertUpdateBlock converts an UpdateBlock packet.
func (c *MinecraftConverter) ConvertUpdateBlock(pk *packet.UpdateBlock) (*packet.UpdateBlock, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlock: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return &dst, nil
}
