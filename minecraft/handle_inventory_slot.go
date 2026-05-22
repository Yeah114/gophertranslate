package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertInventorySlot converts item block runtime IDs inside an InventorySlot packet.
func (c *MinecraftConverter) ConvertInventorySlot(pk *packet.InventorySlot) (*packet.InventorySlot, error) {
	newItem, err := c.bc.ConvertItemInstance(pk.NewItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventorySlot: failed to convert new item: %w", err)
	}
	dst := *pk
	dst.NewItem = newItem
	if storageItem, ok := pk.StorageItem.Value(); ok {
		dstStorageItem, err := c.bc.ConvertItemInstance(storageItem)
		if err != nil {
			return nil, fmt.Errorf("ConvertInventorySlot: failed to convert storage item: %w", err)
		}
		dst.StorageItem = protocol.Option(dstStorageItem)
	}
	return &dst, nil
}
