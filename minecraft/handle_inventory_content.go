package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleInventoryContent converts and writes item block runtime IDs inside an InventoryContent packet.
func (c *MinecraftConverter) HandleInventoryContent(pk *packet.InventoryContent) error {
	content, err := utils.ConvertSliceWithError(pk.Content, c.ic.ConvertServerItemInstance)
	if err != nil {
		return fmt.Errorf("HandleInventoryContent: failed to convert content: %w", err)
	}
	storageItem, err := c.ic.ConvertServerItemInstance(pk.StorageItem)
	if err != nil {
		return fmt.Errorf("HandleInventoryContent: failed to convert storage item: %w", err)
	}
	dst := *pk
	dst.Content = content
	dst.StorageItem = storageItem
	return c.clientConnEcho.WritePacket(&dst)
}
