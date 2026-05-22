package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertUpdateSubChunkBlocks converts an UpdateSubChunkBlocks packet.
func (c *MinecraftConverter) ConvertUpdateSubChunkBlocks(pk *packet.UpdateSubChunkBlocks) (*packet.UpdateSubChunkBlocks, error) {
	blocks, err := utils.ConvertSliceWithError(pk.Blocks, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert blocks: %w", err)
	}
	extra, err := utils.ConvertSliceWithError(pk.Extra, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert extra blocks: %w", err)
	}
	dst := *pk
	dst.Blocks = blocks
	dst.Extra = extra
	return &dst, nil
}
