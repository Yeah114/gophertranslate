package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleUpdateBlockSynced converts and writes an UpdateBlockSynced packet.
func (c *MinecraftConverter) HandleUpdateBlockSynced(pk *packet.UpdateBlockSynced) error {
	newBlockRuntimeID, err := c.bc.ConvertServerBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return fmt.Errorf("HandleUpdateBlockSynced: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return c.clientConnEcho.WritePacket(&dst)
}
