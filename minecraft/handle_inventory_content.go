package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertInventoryContent converts item block runtime IDs inside an InventoryContent packet.
func (c *MinecraftConverter) ConvertInventoryContent(pk *packet.InventoryContent) (*packet.InventoryContent, error) {
	content, err := utils.ConvertSliceWithError(pk.Content, c.bc.ConvertItemInstance)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryContent: failed to convert content: %w", err)
	}
	storageItem, err := c.bc.ConvertItemInstance(pk.StorageItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryContent: failed to convert storage item: %w", err)
	}
	dst := *pk
	dst.Content = content
	dst.StorageItem = storageItem
	return &dst, nil
}
