package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertAddPlayer converts item block runtime IDs inside an AddPlayer packet.
func (c *MinecraftConverter) ConvertAddPlayer(pk *packet.AddPlayer) (*packet.AddPlayer, error) {
	heldItem, err := c.bc.ConvertItemInstance(pk.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddPlayer: failed to convert held item: %w", err)
	}
	dst := *pk
	dst.HeldItem = heldItem
	return &dst, nil
}
