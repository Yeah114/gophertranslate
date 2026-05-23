package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleAddItemActor converts and writes item block runtime IDs inside an AddItemActor packet.
func (c *MinecraftConverter) HandleAddItemActor(pk *packet.AddItemActor) error {
	item, err := c.ic.ConvertServerItemInstance(pk.Item)
	if err != nil {
		return fmt.Errorf("HandleAddItemActor: failed to convert item: %w", err)
	}
	dst := *pk
	dst.Item = item
	return c.clientConnEcho.WritePacket(&dst)
}
