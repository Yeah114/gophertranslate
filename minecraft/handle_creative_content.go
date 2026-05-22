package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertCreativeContent converts item block runtime IDs inside a CreativeContent packet.
func (c *MinecraftConverter) ConvertCreativeContent(pk *packet.CreativeContent) (*packet.CreativeContent, error) {
	groups, err := utils.ConvertSliceWithError(pk.Groups, c.bc.ConvertCreativeGroup)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert groups: %w", err)
	}
	items, err := utils.ConvertSliceWithError(pk.Items, c.bc.ConvertCreativeItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert items: %w", err)
	}
	dst := *pk
	dst.Groups = groups
	dst.Items = items
	return &dst, nil
}
