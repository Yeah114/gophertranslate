package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleCreativeContent converts and writes item block runtime IDs inside a CreativeContent packet.
func (c *MinecraftConverter) HandleCreativeContent(pk *packet.CreativeContent) error {
	groups, err := utils.ConvertSliceWithError(pk.Groups, c.ic.ConvertServerCreativeGroup)
	if err != nil {
		return fmt.Errorf("HandleCreativeContent: failed to convert groups: %w", err)
	}
	items, err := utils.ConvertSliceWithError(pk.Items, c.ic.ConvertServerCreativeItem)
	if err != nil {
		return fmt.Errorf("HandleCreativeContent: failed to convert items: %w", err)
	}
	dst := *pk
	dst.Groups = groups
	dst.Items = items
	return c.clientConnEcho.WritePacket(&dst)
}
