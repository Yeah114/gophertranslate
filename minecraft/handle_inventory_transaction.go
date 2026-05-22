package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertInventoryTransaction converts block runtime IDs inside an InventoryTransaction packet.
func (c *MinecraftConverter) ConvertInventoryTransaction(pk *packet.InventoryTransaction) (*packet.InventoryTransaction, error) {
	actions, err := utils.ConvertSliceWithError(pk.Actions, c.bc.ConvertInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryTransaction: failed to convert actions: %w", err)
	}
	transactionData, err := c.ConvertInventoryTransactionData(pk.TransactionData)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryTransaction: failed to convert transaction data: %w", err)
	}
	dst := *pk
	dst.Actions = actions
	dst.TransactionData = transactionData
	return &dst, nil
}

// ConvertInventoryTransactionData converts block runtime IDs inside inventory transaction data.
func (c *MinecraftConverter) ConvertInventoryTransactionData(data protocol.InventoryTransactionData) (protocol.InventoryTransactionData, error) {
	switch typedData := data.(type) {
	case nil:
		return nil, nil
	case *protocol.NormalTransactionData:
		dst := *typedData
		return &dst, nil
	case *protocol.MismatchTransactionData:
		dst := *typedData
		return &dst, nil
	case *protocol.UseItemTransactionData:
		return c.bc.ConvertUseItemTransactionData(typedData)
	case *protocol.UseItemOnEntityTransactionData:
		return c.bc.ConvertUseItemOnEntityTransactionData(typedData)
	case *protocol.ReleaseItemTransactionData:
		return c.bc.ConvertReleaseItemTransactionData(typedData)
	default:
		return data, nil
	}
}
