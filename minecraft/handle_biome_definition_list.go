package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleBiomeDefinitionList converts and writes block runtime IDs inside a BiomeDefinitionList packet.
func (c *MinecraftConverter) HandleBiomeDefinitionList(pk *packet.BiomeDefinitionList) error {
	biomeDefinitions, err := utils.ConvertSliceWithError(pk.BiomeDefinitions, c.bc.ConvertServerBiomeDefinition)
	if err != nil {
		return fmt.Errorf("HandleBiomeDefinitionList: failed to convert biome definitions: %w", err)
	}
	dst := *pk
	dst.BiomeDefinitions = biomeDefinitions
	return c.clientConnEcho.WritePacket(&dst)
}
