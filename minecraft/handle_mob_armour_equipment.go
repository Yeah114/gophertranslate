package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertMobArmourEquipment converts item block runtime IDs inside a MobArmourEquipment packet.
func (c *MinecraftConverter) ConvertMobArmourEquipment(pk *packet.MobArmourEquipment) (*packet.MobArmourEquipment, error) {
	helmet, err := c.bc.ConvertItemInstance(pk.Helmet)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert helmet: %w", err)
	}
	chestplate, err := c.bc.ConvertItemInstance(pk.Chestplate)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert chestplate: %w", err)
	}
	leggings, err := c.bc.ConvertItemInstance(pk.Leggings)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert leggings: %w", err)
	}
	boots, err := c.bc.ConvertItemInstance(pk.Boots)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert boots: %w", err)
	}
	body, err := c.bc.ConvertItemInstance(pk.Body)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert body: %w", err)
	}
	dst := *pk
	dst.Helmet = helmet
	dst.Chestplate = chestplate
	dst.Leggings = leggings
	dst.Boots = boots
	dst.Body = body
	return &dst, nil
}
