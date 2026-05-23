package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleItemRegistry updates the source item table, converts item network IDs, and writes the packet.
func (c *MinecraftConverter) HandleItemRegistry(pk *packet.ItemRegistry) error {
	if len(pk.Items) == 0 {
		dst := *pk
		return c.clientConnEcho.WritePacket(&dst)
	}
	c.serverItems.ReplaceItems(pk.Items)

	dst := *pk
	dst.Items = make([]protocol.ItemEntry, len(pk.Items))
	for i, entry := range pk.Items {
		serverRuntimeID, err := c.ic.ConvertServerItemRuntimeID(int32(entry.RuntimeID))
		if err != nil {
			return fmt.Errorf("HandleItemRegistry: failed to convert item %s runtime ID: %w", entry.Name, err)
		}
		entry.RuntimeID = int16(serverRuntimeID)
		dst.Items[i] = entry
	}
	return c.clientConnEcho.WritePacket(&dst)
}
