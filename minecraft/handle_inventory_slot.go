package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleInventorySlot converts and writes item block runtime IDs inside an InventorySlot packet.
func (c *MinecraftConverter) HandleInventorySlot(pk *packet.InventorySlot) error {
	newItem, err := c.ic.ConvertServerItemInstance(pk.NewItem)
	if err != nil {
		return fmt.Errorf("HandleInventorySlot: failed to convert new item: %w", err)
	}
	dst := *pk
	dst.NewItem = newItem
	if storageItem, ok := pk.StorageItem.Value(); ok {
		serverStorageItem, err := c.ic.ConvertServerItemInstance(storageItem)
		if err != nil {
			return fmt.Errorf("HandleInventorySlot: failed to convert storage item: %w", err)
		}
		dst.StorageItem = protocol.Option(serverStorageItem)
	}
	return c.clientConnEcho.WritePacket(&dst)
}
