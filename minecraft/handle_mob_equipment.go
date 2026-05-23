package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleMobEquipment converts and writes item block runtime IDs inside a MobEquipment packet.
func (c *MinecraftConverter) HandleMobEquipment(pk *packet.MobEquipment) error {
	newItem, err := c.ic.ConvertClientItemInstance(pk.NewItem)
	if err != nil {
		return fmt.Errorf("HandleMobEquipment: failed to convert new item: %w", err)
	}
	dst := *pk
	dst.NewItem = newItem
	return c.serverConnEcho.WritePacket(&dst)
}
