package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertAddItemActor converts item block runtime IDs inside an AddItemActor packet.
func (c *MinecraftConverter) ConvertAddItemActor(pk *packet.AddItemActor) (*packet.AddItemActor, error) {
	item, err := c.bc.ConvertItemInstance(pk.Item)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddItemActor: failed to convert item: %w", err)
	}
	dst := *pk
	dst.Item = item
	return &dst, nil
}
