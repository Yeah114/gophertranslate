package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertMobEquipment converts item block runtime IDs inside a MobEquipment packet.
func (c *MinecraftConverter) ConvertMobEquipment(pk *packet.MobEquipment) (*packet.MobEquipment, error) {
	newItem, err := c.bc.ConvertItemInstance(pk.NewItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobEquipment: failed to convert new item: %w", err)
	}
	dst := *pk
	dst.NewItem = newItem
	return &dst, nil
}
