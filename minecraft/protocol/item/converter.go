package item

import (
	"fmt"

	protocol_block "github.com/Yeah114/gopherconvert/minecraft/protocol/block"
	world_item "github.com/Yeah114/gopherconvert/minecraft/world/item"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// ItemConverter converts protocol fields that hold item runtime IDs.
type ItemConverter struct {
	ic *world_item.ItemConverter
	bc *protocol_block.BlockConverter
}

// NewItemConverter creates a new protocol item converter.
func NewItemConverter(ic *world_item.ItemConverter, bc *protocol_block.BlockConverter) *ItemConverter {
	return &ItemConverter{ic: ic, bc: bc}
}

// ItemConverter returns the underlying item converter.
func (c *ItemConverter) ItemConverter() *world_item.ItemConverter {
	return c.ic
}

// ConvertClientItemRuntimeID converts an item runtime ID from the client protocol to the server protocol.
func (c *ItemConverter) ConvertClientItemRuntimeID(clientItemRuntimeID int32) (int32, error) {
	serverItemRuntimeID, ok := c.ic.ConvertClientItemRuntimeID(clientItemRuntimeID)
	if !ok {
		return 0, fmt.Errorf("ConvertClientItemRuntimeID: unknown client item runtime ID %d", clientItemRuntimeID)
	}
	return serverItemRuntimeID, nil
}

// ConvertServerItemRuntimeID converts an item runtime ID from the server protocol to the client protocol.
func (c *ItemConverter) ConvertServerItemRuntimeID(serverItemRuntimeID int32) (int32, error) {
	clientItemRuntimeID, ok := c.ic.ConvertServerItemRuntimeID(serverItemRuntimeID)
	if !ok {
		return 0, fmt.Errorf("ConvertServerItemRuntimeID: unknown server item runtime ID %d", serverItemRuntimeID)
	}
	return clientItemRuntimeID, nil
}

// ConvertClientItemStack converts the item and block runtime IDs inside an ItemStack.
func (c *ItemConverter) ConvertClientItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	serverItemStack := clientItemStack
	if clientItemStack.NetworkID == 0 {
		return serverItemStack, nil
	}
	serverItemRuntimeID, err := c.ConvertClientItemRuntimeID(clientItemStack.NetworkID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertClientItemStack: failed to convert item runtime ID: %w", err)
	}
	serverBlockRuntimeID, err := c.bc.ConvertClientBlockRuntimeIDInt32(clientItemStack.BlockRuntimeID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertClientItemStack: failed to convert block runtime ID: %w", err)
	}
	serverItemStack.NetworkID = serverItemRuntimeID
	serverItemStack.BlockRuntimeID = serverBlockRuntimeID
	return serverItemStack, nil
}

// ConvertServerItemStack converts the item and block runtime IDs inside an ItemStack.
func (c *ItemConverter) ConvertServerItemStack(serverItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	clientItemStack := serverItemStack
	if serverItemStack.NetworkID == 0 {
		return clientItemStack, nil
	}
	clientItemRuntimeID, err := c.ConvertServerItemRuntimeID(serverItemStack.NetworkID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertServerItemStack: failed to convert item runtime ID: %w", err)
	}
	clientBlockRuntimeID, err := c.bc.ConvertServerBlockRuntimeIDInt32(serverItemStack.BlockRuntimeID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertServerItemStack: failed to convert block runtime ID: %w", err)
	}
	clientItemStack.NetworkID = clientItemRuntimeID
	clientItemStack.BlockRuntimeID = clientBlockRuntimeID
	return clientItemStack, nil
}

// ConvertClientItemInstance converts the item and block runtime IDs inside an ItemInstance.
func (c *ItemConverter) ConvertClientItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	serverStack, err := c.ConvertClientItemStack(clientItemInstance.Stack)
	if err != nil {
		return protocol.ItemInstance{}, fmt.Errorf("ConvertClientItemInstance: failed to convert stack: %w", err)
	}
	serverItemInstance := clientItemInstance
	serverItemInstance.Stack = serverStack
	return serverItemInstance, nil
}

// ConvertServerItemInstance converts the item and block runtime IDs inside an ItemInstance.
func (c *ItemConverter) ConvertServerItemInstance(serverItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	clientStack, err := c.ConvertServerItemStack(serverItemInstance.Stack)
	if err != nil {
		return protocol.ItemInstance{}, fmt.Errorf("ConvertServerItemInstance: failed to convert stack: %w", err)
	}
	clientItemInstance := serverItemInstance
	clientItemInstance.Stack = clientStack
	return clientItemInstance, nil
}

// ConvertClientCreativeGroup converts the item and block runtime IDs inside a CreativeGroup icon.
func (c *ItemConverter) ConvertClientCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	icon, err := c.ConvertClientItemStack(clientCreativeGroup.Icon)
	if err != nil {
		return protocol.CreativeGroup{}, fmt.Errorf("ConvertClientCreativeGroup: failed to convert icon: %w", err)
	}
	serverCreativeGroup := clientCreativeGroup
	serverCreativeGroup.Icon = icon
	return serverCreativeGroup, nil
}

// ConvertServerCreativeGroup converts the item and block runtime IDs inside a CreativeGroup icon.
func (c *ItemConverter) ConvertServerCreativeGroup(serverCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	icon, err := c.ConvertServerItemStack(serverCreativeGroup.Icon)
	if err != nil {
		return protocol.CreativeGroup{}, fmt.Errorf("ConvertServerCreativeGroup: failed to convert icon: %w", err)
	}
	clientCreativeGroup := serverCreativeGroup
	clientCreativeGroup.Icon = icon
	return clientCreativeGroup, nil
}

// ConvertClientCreativeItem converts the item and block runtime IDs inside a CreativeItem.
func (c *ItemConverter) ConvertClientCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	item, err := c.ConvertClientItemStack(clientCreativeItem.Item)
	if err != nil {
		return protocol.CreativeItem{}, fmt.Errorf("ConvertClientCreativeItem: failed to convert item: %w", err)
	}
	serverCreativeItem := clientCreativeItem
	serverCreativeItem.Item = item
	return serverCreativeItem, nil
}

// ConvertServerCreativeItem converts the item and block runtime IDs inside a CreativeItem.
func (c *ItemConverter) ConvertServerCreativeItem(serverCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	item, err := c.ConvertServerItemStack(serverCreativeItem.Item)
	if err != nil {
		return protocol.CreativeItem{}, fmt.Errorf("ConvertServerCreativeItem: failed to convert item: %w", err)
	}
	clientCreativeItem := serverCreativeItem
	clientCreativeItem.Item = item
	return clientCreativeItem, nil
}

// ConvertClientInventoryAction converts item and block runtime IDs inside an InventoryAction.
func (c *ItemConverter) ConvertClientInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
	serverOldItem, err := c.ConvertClientItemInstance(clientInventoryAction.OldItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertClientInventoryAction: failed to convert old item: %w", err)
	}
	serverNewItem, err := c.ConvertClientItemInstance(clientInventoryAction.NewItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertClientInventoryAction: failed to convert new item: %w", err)
	}
	serverInventoryAction := clientInventoryAction
	serverInventoryAction.OldItem = serverOldItem
	serverInventoryAction.NewItem = serverNewItem
	return serverInventoryAction, nil
}

// ConvertServerInventoryAction converts item and block runtime IDs inside an InventoryAction.
func (c *ItemConverter) ConvertServerInventoryAction(serverInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
	clientOldItem, err := c.ConvertServerItemInstance(serverInventoryAction.OldItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertServerInventoryAction: failed to convert old item: %w", err)
	}
	clientNewItem, err := c.ConvertServerItemInstance(serverInventoryAction.NewItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertServerInventoryAction: failed to convert new item: %w", err)
	}
	clientInventoryAction := serverInventoryAction
	clientInventoryAction.OldItem = clientOldItem
	clientInventoryAction.NewItem = clientNewItem
	return clientInventoryAction, nil
}

// ConvertClientUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData.
func (c *ItemConverter) ConvertClientUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	actions, err := convertSliceWithError(clientData.Actions, c.ConvertClientInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert actions: %w", err)
	}
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert held item: %w", err)
	}
	blockRuntimeID, err := c.bc.ConvertClientBlockRuntimeID(clientData.BlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert block runtime ID: %w", err)
	}
	serverData := *clientData
	serverData.Actions = actions
	serverData.HeldItem = heldItem
	serverData.BlockRuntimeID = blockRuntimeID
	return &serverData, nil
}

// ConvertServerUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData.
func (c *ItemConverter) ConvertServerUseItemTransactionData(serverData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	actions, err := convertSliceWithError(serverData.Actions, c.ConvertServerInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert actions: %w", err)
	}
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert held item: %w", err)
	}
	blockRuntimeID, err := c.bc.ConvertServerBlockRuntimeID(serverData.BlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert block runtime ID: %w", err)
	}
	clientData := *serverData
	clientData.Actions = actions
	clientData.HeldItem = heldItem
	clientData.BlockRuntimeID = blockRuntimeID
	return &clientData, nil
}

// ConvertClientUseItemOnEntityTransactionData converts item and block runtime IDs inside UseItemOnEntityTransactionData.
func (c *ItemConverter) ConvertClientUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemOnEntityTransactionData: failed to convert held item: %w", err)
	}
	serverData := *clientData
	serverData.HeldItem = heldItem
	return &serverData, nil
}

// ConvertServerUseItemOnEntityTransactionData converts item and block runtime IDs inside UseItemOnEntityTransactionData.
func (c *ItemConverter) ConvertServerUseItemOnEntityTransactionData(serverData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemOnEntityTransactionData: failed to convert held item: %w", err)
	}
	clientData := *serverData
	clientData.HeldItem = heldItem
	return &clientData, nil
}

// ConvertClientReleaseItemTransactionData converts item and block runtime IDs inside ReleaseItemTransactionData.
func (c *ItemConverter) ConvertClientReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientReleaseItemTransactionData: failed to convert held item: %w", err)
	}
	serverData := *clientData
	serverData.HeldItem = heldItem
	return &serverData, nil
}

// ConvertServerReleaseItemTransactionData converts item and block runtime IDs inside ReleaseItemTransactionData.
func (c *ItemConverter) ConvertServerReleaseItemTransactionData(serverData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerReleaseItemTransactionData: failed to convert held item: %w", err)
	}
	clientData := *serverData
	clientData.HeldItem = heldItem
	return &clientData, nil
}

// ConvertClientInventoryTransactionData converts item and block runtime IDs inside InventoryTransactionData.
func (c *ItemConverter) ConvertClientInventoryTransactionData(clientData protocol.InventoryTransactionData) (protocol.InventoryTransactionData, error) {
	switch typedData := clientData.(type) {
	case nil:
		return nil, nil
	case *protocol.NormalTransactionData:
		client := *typedData
		return &client, nil
	case *protocol.MismatchTransactionData:
		client := *typedData
		return &client, nil
	case *protocol.UseItemTransactionData:
		return c.ConvertClientUseItemTransactionData(typedData)
	case *protocol.UseItemOnEntityTransactionData:
		return c.ConvertClientUseItemOnEntityTransactionData(typedData)
	case *protocol.ReleaseItemTransactionData:
		return c.ConvertClientReleaseItemTransactionData(typedData)
	default:
		return clientData, nil
	}
}

// ConvertItemRuntimeID converts an item runtime ID from the client protocol to the server protocol.
func (c *ItemConverter) ConvertItemRuntimeID(clientItemRuntimeID int32) (int32, error) {
	return c.ConvertClientItemRuntimeID(clientItemRuntimeID)
}

// ConvertItemStack converts the item and block runtime IDs inside an ItemStack.
func (c *ItemConverter) ConvertItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	return c.ConvertClientItemStack(clientItemStack)
}

// ConvertItemInstance converts the item and block runtime IDs inside an ItemInstance.
func (c *ItemConverter) ConvertItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	return c.ConvertClientItemInstance(clientItemInstance)
}

// ConvertCreativeGroup converts the item and block runtime IDs inside a CreativeGroup icon.
func (c *ItemConverter) ConvertCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	return c.ConvertClientCreativeGroup(clientCreativeGroup)
}

// ConvertCreativeItem converts the item and block runtime IDs inside a CreativeItem.
func (c *ItemConverter) ConvertCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	return c.ConvertClientCreativeItem(clientCreativeItem)
}

// ConvertInventoryAction converts item and block runtime IDs inside an InventoryAction.
func (c *ItemConverter) ConvertInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
	return c.ConvertClientInventoryAction(clientInventoryAction)
}

// ConvertUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData.
func (c *ItemConverter) ConvertUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	return c.ConvertClientUseItemTransactionData(clientData)
}

// ConvertUseItemOnEntityTransactionData converts item and block runtime IDs inside UseItemOnEntityTransactionData.
func (c *ItemConverter) ConvertUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	return c.ConvertClientUseItemOnEntityTransactionData(clientData)
}

// ConvertReleaseItemTransactionData converts item and block runtime IDs inside ReleaseItemTransactionData.
func (c *ItemConverter) ConvertReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	return c.ConvertClientReleaseItemTransactionData(clientData)
}

func convertSliceWithError[S ~[]E, E any](client S, convert func(E) (E, error)) (S, error) {
	server := make(S, len(client))
	for i, value := range client {
		converted, err := convert(value)
		if err != nil {
			return nil, err
		}
		server[i] = converted
	}
	return server, nil
}
