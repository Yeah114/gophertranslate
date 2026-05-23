package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleInventoryTransaction converts and writes block runtime IDs inside an InventoryTransaction packet.
func (c *MinecraftConverter) HandleInventoryTransaction(pk *packet.InventoryTransaction) error {
	actions, err := utils.ConvertSliceWithError(pk.Actions, c.ic.ConvertClientInventoryAction)
	if err != nil {
		return fmt.Errorf("HandleInventoryTransaction: failed to convert actions: %w", err)
	}
	transactionData, err := c.ic.ConvertClientInventoryTransactionData(pk.TransactionData)
	if err != nil {
		return fmt.Errorf("HandleInventoryTransaction: failed to convert transaction data: %w", err)
	}
	dst := *pk
	dst.Actions = actions
	dst.TransactionData = transactionData
	return c.serverConnEcho.WritePacket(&dst)
}
