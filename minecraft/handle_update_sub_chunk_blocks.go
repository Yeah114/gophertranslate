package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleUpdateSubChunkBlocks converts and writes an UpdateSubChunkBlocks packet.
func (c *MinecraftConverter) HandleUpdateSubChunkBlocks(pk *packet.UpdateSubChunkBlocks) error {
	blocks, err := utils.ConvertSliceWithError(pk.Blocks, c.bc.ConvertServerBlockChangeEntry)
	if err != nil {
		return fmt.Errorf("HandleUpdateSubChunkBlocks: failed to convert blocks: %w", err)
	}
	extra, err := utils.ConvertSliceWithError(pk.Extra, c.bc.ConvertServerBlockChangeEntry)
	if err != nil {
		return fmt.Errorf("HandleUpdateSubChunkBlocks: failed to convert extra blocks: %w", err)
	}
	dst := *pk
	dst.Blocks = blocks
	dst.Extra = extra
	return c.clientConnEcho.WritePacket(&dst)
}
