package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertBiomeDefinitionList converts block runtime IDs inside a BiomeDefinitionList packet.
func (c *MinecraftConverter) ConvertBiomeDefinitionList(pk *packet.BiomeDefinitionList) (*packet.BiomeDefinitionList, error) {
	biomeDefinitions, err := utils.ConvertSliceWithError(pk.BiomeDefinitions, c.bc.ConvertBiomeDefinition)
	if err != nil {
		return nil, fmt.Errorf("ConvertBiomeDefinitionList: failed to convert biome definitions: %w", err)
	}
	dst := *pk
	dst.BiomeDefinitions = biomeDefinitions
	return &dst, nil
}
