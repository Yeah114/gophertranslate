package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleUpdateBlock converts and writes an UpdateBlock packet.
func (c *MinecraftConverter) HandleUpdateBlock(pk *packet.UpdateBlock) error {
	newBlockRuntimeID, err := c.bc.ConvertServerBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return fmt.Errorf("HandleUpdateBlock: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return c.clientConnEcho.WritePacket(&dst)
}
