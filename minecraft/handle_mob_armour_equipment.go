package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleMobArmourEquipment converts and writes item block runtime IDs inside a MobArmourEquipment packet.
func (c *MinecraftConverter) HandleMobArmourEquipment(pk *packet.MobArmourEquipment) error {
	helmet, err := c.ic.ConvertServerItemInstance(pk.Helmet)
	if err != nil {
		return fmt.Errorf("HandleMobArmourEquipment: failed to convert helmet: %w", err)
	}
	chestplate, err := c.ic.ConvertServerItemInstance(pk.Chestplate)
	if err != nil {
		return fmt.Errorf("HandleMobArmourEquipment: failed to convert chestplate: %w", err)
	}
	leggings, err := c.ic.ConvertServerItemInstance(pk.Leggings)
	if err != nil {
		return fmt.Errorf("HandleMobArmourEquipment: failed to convert leggings: %w", err)
	}
	boots, err := c.ic.ConvertServerItemInstance(pk.Boots)
	if err != nil {
		return fmt.Errorf("HandleMobArmourEquipment: failed to convert boots: %w", err)
	}
	body, err := c.ic.ConvertServerItemInstance(pk.Body)
	if err != nil {
		return fmt.Errorf("HandleMobArmourEquipment: failed to convert body: %w", err)
	}
	dst := *pk
	dst.Helmet = helmet
	dst.Chestplate = chestplate
	dst.Leggings = leggings
	dst.Boots = boots
	dst.Body = body
	return c.clientConnEcho.WritePacket(&dst)
}
