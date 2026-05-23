package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleAddPlayer converts and writes item block runtime IDs inside an AddPlayer packet.
func (c *MinecraftConverter) HandleAddPlayer(pk *packet.AddPlayer) error {
	heldItem, err := c.ic.ConvertServerItemInstance(pk.HeldItem)
	if err != nil {
		return fmt.Errorf("HandleAddPlayer: failed to convert held item: %w", err)
	}
	dst := *pk
	dst.HeldItem = heldItem
	return c.clientConnEcho.WritePacket(&dst)
}
