package block

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// BlockConverter converts protocol fields that hold block runtime IDs.
type BlockConverter struct {
	bc *world_block.BlockConverter
}

// NewBlockConverter creates a new protocol block converter.
func NewBlockConverter(bc *world_block.BlockConverter) *BlockConverter {
	return &BlockConverter{bc: bc}
}

// BlockConverter returns the underlying block converter.
func (c *BlockConverter) BlockConverter() *world_block.BlockConverter {
	return c.bc
}

// ConvertClientBlockRuntimeID converts a block runtime ID from the source protocol to the destination protocol.
func (c *BlockConverter) ConvertClientBlockRuntimeID(clientBlockRuntimeID uint32) (uint32, error) {
	serverBlockRuntimeID, ok := c.bc.ConvertClientBlockRuntimeID(clientBlockRuntimeID)
	if !ok {
		return 0, fmt.Errorf("ConvertClientBlockRuntimeID: unknown source block runtime ID %d", clientBlockRuntimeID)
	}
	return serverBlockRuntimeID, nil
}

// ConvertClientBlockRuntimeIDInt32 converts an int32 block runtime ID.
func (c *BlockConverter) ConvertClientBlockRuntimeIDInt32(clientBlockRuntimeID int32) (int32, error) {
	if clientBlockRuntimeID < 0 {
		return clientBlockRuntimeID, nil
	}
	serverBlockRuntimeID, err := c.ConvertClientBlockRuntimeID(uint32(clientBlockRuntimeID))
	if err != nil {
		return 0, err
	}
	return int32(serverBlockRuntimeID), nil
}

// ConvertClientItemStack converts the block runtime ID inside an ItemStack.
func (c *BlockConverter) ConvertClientItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	serverItemStack := clientItemStack
	if clientItemStack.NetworkID == 0 {
		return serverItemStack, nil
	}
	serverBlockRuntimeID, err := c.ConvertClientBlockRuntimeIDInt32(clientItemStack.BlockRuntimeID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertClientItemStack: failed to convert block runtime ID: %w", err)
	}
	serverItemStack.BlockRuntimeID = serverBlockRuntimeID
	return serverItemStack, nil
}

// ConvertClientItemInstance converts the block runtime ID inside an ItemInstance.
func (c *BlockConverter) ConvertClientItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	serverStack, err := c.ConvertClientItemStack(clientItemInstance.Stack)
	if err != nil {
		return protocol.ItemInstance{}, fmt.Errorf("ConvertClientItemInstance: failed to convert stack: %w", err)
	}
	serverItemInstance := clientItemInstance
	serverItemInstance.Stack = serverStack
	return serverItemInstance, nil
}

// ConvertClientCreativeGroup converts the block runtime ID inside a CreativeGroup icon.
func (c *BlockConverter) ConvertClientCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	icon, err := c.ConvertClientItemStack(clientCreativeGroup.Icon)
	if err != nil {
		return protocol.CreativeGroup{}, fmt.Errorf("ConvertClientCreativeGroup: failed to convert icon: %w", err)
	}
	serverCreativeGroup := clientCreativeGroup
	serverCreativeGroup.Icon = icon
	return serverCreativeGroup, nil
}

// ConvertClientCreativeItem converts the block runtime ID inside a CreativeItem.
func (c *BlockConverter) ConvertClientCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	item, err := c.ConvertClientItemStack(clientCreativeItem.Item)
	if err != nil {
		return protocol.CreativeItem{}, fmt.Errorf("ConvertClientCreativeItem: failed to convert item: %w", err)
	}
	serverCreativeItem := clientCreativeItem
	serverCreativeItem.Item = item
	return serverCreativeItem, nil
}

// ConvertClientInventoryAction converts item block runtime IDs inside an InventoryAction.
func (c *BlockConverter) ConvertClientInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
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

// ConvertClientBlockChangeEntry converts the block runtime ID inside a BlockChangeEntry.
func (c *BlockConverter) ConvertClientBlockChangeEntry(clientBlockChangeEntry protocol.BlockChangeEntry) (protocol.BlockChangeEntry, error) {
	serverBlockRuntimeID, err := c.ConvertClientBlockRuntimeID(clientBlockChangeEntry.BlockRuntimeID)
	if err != nil {
		return protocol.BlockChangeEntry{}, fmt.Errorf("ConvertClientBlockChangeEntry: failed to convert block runtime ID: %w", err)
	}
	serverBlockChangeEntry := clientBlockChangeEntry
	serverBlockChangeEntry.BlockRuntimeID = serverBlockRuntimeID
	return serverBlockChangeEntry, nil
}

// ConvertClientUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData.
func (c *BlockConverter) ConvertClientUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	actions, err := utils.ConvertSliceWithError(clientData.Actions, c.ConvertClientInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert actions: %w", err)
	}
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert held item: %w", err)
	}
	blockRuntimeID, err := c.ConvertClientBlockRuntimeID(clientData.BlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemTransactionData: failed to convert block runtime ID: %w", err)
	}
	serverData := *clientData
	serverData.Actions = actions
	serverData.HeldItem = heldItem
	serverData.BlockRuntimeID = blockRuntimeID
	return &serverData, nil
}

// ConvertClientUseItemOnEntityTransactionData converts item block runtime IDs inside UseItemOnEntityTransactionData.
func (c *BlockConverter) ConvertClientUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientUseItemOnEntityTransactionData: failed to convert held item: %w", err)
	}
	serverData := *clientData
	serverData.HeldItem = heldItem
	return &serverData, nil
}

// ConvertClientReleaseItemTransactionData converts item block runtime IDs inside ReleaseItemTransactionData.
func (c *BlockConverter) ConvertClientReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	heldItem, err := c.ConvertClientItemInstance(clientData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientReleaseItemTransactionData: failed to convert held item: %w", err)
	}
	serverData := *clientData
	serverData.HeldItem = heldItem
	return &serverData, nil
}

// ConvertClientWaxedOrUnwaxedCopperEvent converts the copper block runtime ID inside WaxedOrUnwaxedCopperEvent.
func (c *BlockConverter) ConvertClientWaxedOrUnwaxedCopperEvent(clientEvent *protocol.WaxedOrUnwaxedCopperEvent) (*protocol.WaxedOrUnwaxedCopperEvent, error) {
	serverCopperBlockID, err := c.ConvertClientBlockRuntimeIDInt32(clientEvent.CopperBlockID)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientWaxedOrUnwaxedCopperEvent: failed to convert copper block ID: %w", err)
	}
	serverEvent := *clientEvent
	serverEvent.CopperBlockID = serverCopperBlockID
	return &serverEvent, nil
}

// ConvertClientEvent converts block runtime IDs inside event data.
func (c *BlockConverter) ConvertClientEvent(clientEvent protocol.Event) (protocol.Event, error) {
	switch event := clientEvent.(type) {
	case nil:
		return nil, nil
	case *protocol.WaxedOrUnwaxedCopperEvent:
		return c.ConvertClientWaxedOrUnwaxedCopperEvent(event)
	default:
		return clientEvent, nil
	}
}

// ConvertClientBiomeDefinition converts block runtime IDs inside a BiomeDefinition.
func (c *BlockConverter) ConvertClientBiomeDefinition(clientBiomeDefinition protocol.BiomeDefinition) (protocol.BiomeDefinition, error) {
	chunkGeneration, ok := clientBiomeDefinition.ChunkGeneration.Value()
	if !ok {
		return clientBiomeDefinition, nil
	}
	serverChunkGeneration, err := c.ConvertClientBiomeChunkGeneration(chunkGeneration)
	if err != nil {
		return protocol.BiomeDefinition{}, fmt.Errorf("ConvertClientBiomeDefinition: failed to convert chunk generation: %w", err)
	}
	serverBiomeDefinition := clientBiomeDefinition
	serverBiomeDefinition.ChunkGeneration = protocol.Option(serverChunkGeneration)
	return serverBiomeDefinition, nil
}

// ConvertClientBiomeChunkGeneration converts block runtime IDs inside BiomeChunkGeneration.
func (c *BlockConverter) ConvertClientBiomeChunkGeneration(clientChunkGeneration protocol.BiomeChunkGeneration) (protocol.BiomeChunkGeneration, error) {
	serverChunkGeneration := clientChunkGeneration

	if mountainParameters, ok := clientChunkGeneration.MountainParameters.Value(); ok {
		serverMountainParameters, err := c.ConvertClientBiomeMountainParameters(mountainParameters)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertClientBiomeChunkGeneration: failed to convert mountain parameters: %w", err)
		}
		serverChunkGeneration.MountainParameters = protocol.Option(serverMountainParameters)
	}
	if surfaceMaterialAdjustments, ok := clientChunkGeneration.SurfaceMaterialAdjustments.Value(); ok {
		serverSurfaceMaterialAdjustments, err := utils.ConvertSliceWithError(surfaceMaterialAdjustments, c.ConvertClientBiomeElementData)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertClientBiomeChunkGeneration: failed to convert surface material adjustments: %w", err)
		}
		serverChunkGeneration.SurfaceMaterialAdjustments = protocol.Option(serverSurfaceMaterialAdjustments)
	}
	if surfaceMaterials, ok := clientChunkGeneration.SurfaceMaterials.Value(); ok {
		serverSurfaceMaterials, err := c.ConvertClientBiomeSurfaceMaterial(surfaceMaterials)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertClientBiomeChunkGeneration: failed to convert surface materials: %w", err)
		}
		serverChunkGeneration.SurfaceMaterials = protocol.Option(serverSurfaceMaterials)
	}
	if mesaSurface, ok := clientChunkGeneration.MesaSurface.Value(); ok {
		serverMesaSurface, err := c.ConvertClientBiomeMesaSurface(mesaSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertClientBiomeChunkGeneration: failed to convert mesa surface: %w", err)
		}
		serverChunkGeneration.MesaSurface = protocol.Option(serverMesaSurface)
	}
	if cappedSurface, ok := clientChunkGeneration.CappedSurface.Value(); ok {
		serverCappedSurface, err := c.ConvertClientBiomeCappedSurface(cappedSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertClientBiomeChunkGeneration: failed to convert capped surface: %w", err)
		}
		serverChunkGeneration.CappedSurface = protocol.Option(serverCappedSurface)
	}

	return serverChunkGeneration, nil
}

// ConvertClientBiomeMountainParameters converts block runtime IDs inside BiomeMountainParameters.
func (c *BlockConverter) ConvertClientBiomeMountainParameters(clientParameters protocol.BiomeMountainParameters) (protocol.BiomeMountainParameters, error) {
	serverSteepBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientParameters.SteepBlock)
	if err != nil {
		return protocol.BiomeMountainParameters{}, fmt.Errorf("ConvertClientBiomeMountainParameters: failed to convert steep block: %w", err)
	}
	serverParameters := clientParameters
	serverParameters.SteepBlock = serverSteepBlock
	return serverParameters, nil
}

// ConvertClientBiomeElementData converts block runtime IDs inside BiomeElementData.
func (c *BlockConverter) ConvertClientBiomeElementData(clientElementData protocol.BiomeElementData) (protocol.BiomeElementData, error) {
	serverAdjustedMaterials, err := c.ConvertClientBiomeSurfaceMaterial(clientElementData.AdjustedMaterials)
	if err != nil {
		return protocol.BiomeElementData{}, fmt.Errorf("ConvertClientBiomeElementData: failed to convert adjusted materials: %w", err)
	}
	serverElementData := clientElementData
	serverElementData.AdjustedMaterials = serverAdjustedMaterials
	return serverElementData, nil
}

// ConvertClientBiomeSurfaceMaterial converts block runtime IDs inside BiomeSurfaceMaterial.
func (c *BlockConverter) ConvertClientBiomeSurfaceMaterial(clientSurfaceMaterial protocol.BiomeSurfaceMaterial) (protocol.BiomeSurfaceMaterial, error) {
	topBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientSurfaceMaterial.TopBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertClientBiomeSurfaceMaterial: failed to convert top block: %w", err)
	}
	midBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientSurfaceMaterial.MidBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertClientBiomeSurfaceMaterial: failed to convert mid block: %w", err)
	}
	seaFloorBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientSurfaceMaterial.SeaFloorBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertClientBiomeSurfaceMaterial: failed to convert sea floor block: %w", err)
	}
	foundationBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientSurfaceMaterial.FoundationBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertClientBiomeSurfaceMaterial: failed to convert foundation block: %w", err)
	}
	seaBlock, err := c.ConvertClientBlockRuntimeIDInt32(clientSurfaceMaterial.SeaBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertClientBiomeSurfaceMaterial: failed to convert sea block: %w", err)
	}
	serverSurfaceMaterial := clientSurfaceMaterial
	serverSurfaceMaterial.TopBlock = topBlock
	serverSurfaceMaterial.MidBlock = midBlock
	serverSurfaceMaterial.SeaFloorBlock = seaFloorBlock
	serverSurfaceMaterial.FoundationBlock = foundationBlock
	serverSurfaceMaterial.SeaBlock = seaBlock
	return serverSurfaceMaterial, nil
}

// ConvertClientBiomeMesaSurface converts block runtime IDs inside BiomeMesaSurface.
func (c *BlockConverter) ConvertClientBiomeMesaSurface(clientMesaSurface protocol.BiomeMesaSurface) (protocol.BiomeMesaSurface, error) {
	clayMaterial, err := c.ConvertClientBlockRuntimeID(clientMesaSurface.ClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertClientBiomeMesaSurface: failed to convert clay material: %w", err)
	}
	hardClayMaterial, err := c.ConvertClientBlockRuntimeID(clientMesaSurface.HardClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertClientBiomeMesaSurface: failed to convert hard clay material: %w", err)
	}
	serverMesaSurface := clientMesaSurface
	serverMesaSurface.ClayMaterial = clayMaterial
	serverMesaSurface.HardClayMaterial = hardClayMaterial
	return serverMesaSurface, nil
}

// ConvertClientBiomeCappedSurface converts block runtime IDs inside BiomeCappedSurface.
func (c *BlockConverter) ConvertClientBiomeCappedSurface(clientCappedSurface protocol.BiomeCappedSurface) (protocol.BiomeCappedSurface, error) {
	floorBlocks, err := utils.ConvertSliceWithError(clientCappedSurface.FloorBlocks, c.ConvertClientBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertClientBiomeCappedSurface: failed to convert floor blocks: %w", err)
	}
	ceilingBlocks, err := utils.ConvertSliceWithError(clientCappedSurface.CeilingBlocks, c.ConvertClientBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertClientBiomeCappedSurface: failed to convert ceiling blocks: %w", err)
	}
	serverCappedSurface := clientCappedSurface
	serverCappedSurface.FloorBlocks = floorBlocks
	serverCappedSurface.CeilingBlocks = ceilingBlocks

	if seaBlock, ok := clientCappedSurface.SeaBlock.Value(); ok {
		serverSeaBlock, err := c.ConvertClientBlockRuntimeID(seaBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertClientBiomeCappedSurface: failed to convert sea block: %w", err)
		}
		serverCappedSurface.SeaBlock = protocol.Option(serverSeaBlock)
	}
	if foundationBlock, ok := clientCappedSurface.FoundationBlock.Value(); ok {
		serverFoundationBlock, err := c.ConvertClientBlockRuntimeID(foundationBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertClientBiomeCappedSurface: failed to convert foundation block: %w", err)
		}
		serverCappedSurface.FoundationBlock = protocol.Option(serverFoundationBlock)
	}
	if beachBlock, ok := clientCappedSurface.BeachBlock.Value(); ok {
		serverBeachBlock, err := c.ConvertClientBlockRuntimeID(beachBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertClientBiomeCappedSurface: failed to convert beach block: %w", err)
		}
		serverCappedSurface.BeachBlock = protocol.Option(serverBeachBlock)
	}

	return serverCappedSurface, nil
}

// ConvertServerBlockRuntimeID converts a block runtime ID from the server protocol to the client protocol.
func (c *BlockConverter) ConvertServerBlockRuntimeID(serverBlockRuntimeID uint32) (uint32, error) {
	clientBlockRuntimeID, ok := c.bc.ConvertServerBlockRuntimeID(serverBlockRuntimeID)
	if !ok {
		return 0, fmt.Errorf("ConvertServerBlockRuntimeID: unknown server block runtime ID %d", serverBlockRuntimeID)
	}
	return clientBlockRuntimeID, nil
}

// ConvertServerBlockRuntimeIDInt32 converts an int32 block runtime ID (server→client).
func (c *BlockConverter) ConvertServerBlockRuntimeIDInt32(serverBlockRuntimeID int32) (int32, error) {
	if serverBlockRuntimeID < 0 {
		return serverBlockRuntimeID, nil
	}
	clientBlockRuntimeID, err := c.ConvertServerBlockRuntimeID(uint32(serverBlockRuntimeID))
	if err != nil {
		return 0, err
	}
	return int32(clientBlockRuntimeID), nil
}

// ConvertServerItemStack converts the block runtime ID inside an ItemStack (server→client).
func (c *BlockConverter) ConvertServerItemStack(serverItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	clientItemStack := serverItemStack
	if serverItemStack.NetworkID == 0 {
		return clientItemStack, nil
	}
	clientBlockRuntimeID, err := c.ConvertServerBlockRuntimeIDInt32(serverItemStack.BlockRuntimeID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertServerItemStack: failed to convert block runtime ID: %w", err)
	}
	clientItemStack.BlockRuntimeID = clientBlockRuntimeID
	return clientItemStack, nil
}

// ConvertServerItemInstance converts the block runtime ID inside an ItemInstance (server→client).
func (c *BlockConverter) ConvertServerItemInstance(serverItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	clientStack, err := c.ConvertServerItemStack(serverItemInstance.Stack)
	if err != nil {
		return protocol.ItemInstance{}, fmt.Errorf("ConvertServerItemInstance: failed to convert stack: %w", err)
	}
	clientItemInstance := serverItemInstance
	clientItemInstance.Stack = clientStack
	return clientItemInstance, nil
}

// ConvertServerCreativeGroup converts the block runtime ID inside a CreativeGroup icon (server→client).
func (c *BlockConverter) ConvertServerCreativeGroup(serverCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	icon, err := c.ConvertServerItemStack(serverCreativeGroup.Icon)
	if err != nil {
		return protocol.CreativeGroup{}, fmt.Errorf("ConvertServerCreativeGroup: failed to convert icon: %w", err)
	}
	clientCreativeGroup := serverCreativeGroup
	clientCreativeGroup.Icon = icon
	return clientCreativeGroup, nil
}

// ConvertServerCreativeItem converts the block runtime ID inside a CreativeItem (server→client).
func (c *BlockConverter) ConvertServerCreativeItem(serverCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	item, err := c.ConvertServerItemStack(serverCreativeItem.Item)
	if err != nil {
		return protocol.CreativeItem{}, fmt.Errorf("ConvertServerCreativeItem: failed to convert item: %w", err)
	}
	clientCreativeItem := serverCreativeItem
	clientCreativeItem.Item = item
	return clientCreativeItem, nil
}

// ConvertServerInventoryAction converts item block runtime IDs inside an InventoryAction (server→client).
func (c *BlockConverter) ConvertServerInventoryAction(serverInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
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

// ConvertServerBlockChangeEntry converts the block runtime ID inside a BlockChangeEntry (server→client).
func (c *BlockConverter) ConvertServerBlockChangeEntry(serverBlockChangeEntry protocol.BlockChangeEntry) (protocol.BlockChangeEntry, error) {
	clientBlockRuntimeID, err := c.ConvertServerBlockRuntimeID(serverBlockChangeEntry.BlockRuntimeID)
	if err != nil {
		return protocol.BlockChangeEntry{}, fmt.Errorf("ConvertServerBlockChangeEntry: failed to convert block runtime ID: %w", err)
	}
	clientBlockChangeEntry := serverBlockChangeEntry
	clientBlockChangeEntry.BlockRuntimeID = clientBlockRuntimeID
	return clientBlockChangeEntry, nil
}

// ConvertServerUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData (server→client).
func (c *BlockConverter) ConvertServerUseItemTransactionData(serverData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	actions, err := utils.ConvertSliceWithError(serverData.Actions, c.ConvertServerInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert actions: %w", err)
	}
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert held item: %w", err)
	}
	blockRuntimeID, err := c.ConvertServerBlockRuntimeID(serverData.BlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemTransactionData: failed to convert block runtime ID: %w", err)
	}
	clientData := *serverData
	clientData.Actions = actions
	clientData.HeldItem = heldItem
	clientData.BlockRuntimeID = blockRuntimeID
	return &clientData, nil
}

// ConvertServerUseItemOnEntityTransactionData converts item block runtime IDs (server→client).
func (c *BlockConverter) ConvertServerUseItemOnEntityTransactionData(serverData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerUseItemOnEntityTransactionData: failed to convert held item: %w", err)
	}
	clientData := *serverData
	clientData.HeldItem = heldItem
	return &clientData, nil
}

// ConvertServerReleaseItemTransactionData converts item block runtime IDs (server→client).
func (c *BlockConverter) ConvertServerReleaseItemTransactionData(serverData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	heldItem, err := c.ConvertServerItemInstance(serverData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerReleaseItemTransactionData: failed to convert held item: %w", err)
	}
	clientData := *serverData
	clientData.HeldItem = heldItem
	return &clientData, nil
}

// ConvertServerWaxedOrUnwaxedCopperEvent converts the copper block runtime ID (server→client).
func (c *BlockConverter) ConvertServerWaxedOrUnwaxedCopperEvent(serverEvent *protocol.WaxedOrUnwaxedCopperEvent) (*protocol.WaxedOrUnwaxedCopperEvent, error) {
	clientCopperBlockID, err := c.ConvertServerBlockRuntimeIDInt32(serverEvent.CopperBlockID)
	if err != nil {
		return nil, fmt.Errorf("ConvertServerWaxedOrUnwaxedCopperEvent: failed to convert copper block ID: %w", err)
	}
	clientEvent := *serverEvent
	clientEvent.CopperBlockID = clientCopperBlockID
	return &clientEvent, nil
}

// ConvertServerEvent converts block runtime IDs inside event data (server→client).
func (c *BlockConverter) ConvertServerEvent(serverEvent protocol.Event) (protocol.Event, error) {
	switch event := serverEvent.(type) {
	case nil:
		return nil, nil
	case *protocol.WaxedOrUnwaxedCopperEvent:
		return c.ConvertServerWaxedOrUnwaxedCopperEvent(event)
	default:
		return serverEvent, nil
	}
}

// ConvertServerBiomeDefinition converts block runtime IDs inside a BiomeDefinition (server→client).
func (c *BlockConverter) ConvertServerBiomeDefinition(serverBiomeDefinition protocol.BiomeDefinition) (protocol.BiomeDefinition, error) {
	chunkGeneration, ok := serverBiomeDefinition.ChunkGeneration.Value()
	if !ok {
		return serverBiomeDefinition, nil
	}
	clientChunkGeneration, err := c.ConvertServerBiomeChunkGeneration(chunkGeneration)
	if err != nil {
		return protocol.BiomeDefinition{}, fmt.Errorf("ConvertServerBiomeDefinition: failed to convert chunk generation: %w", err)
	}
	clientBiomeDefinition := serverBiomeDefinition
	clientBiomeDefinition.ChunkGeneration = protocol.Option(clientChunkGeneration)
	return clientBiomeDefinition, nil
}

// ConvertServerBiomeChunkGeneration converts block runtime IDs inside BiomeChunkGeneration (server→client).
func (c *BlockConverter) ConvertServerBiomeChunkGeneration(serverChunkGeneration protocol.BiomeChunkGeneration) (protocol.BiomeChunkGeneration, error) {
	clientChunkGeneration := serverChunkGeneration

	if mountainParameters, ok := serverChunkGeneration.MountainParameters.Value(); ok {
		clientMountainParameters, err := c.ConvertServerBiomeMountainParameters(mountainParameters)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertServerBiomeChunkGeneration: failed to convert mountain parameters: %w", err)
		}
		clientChunkGeneration.MountainParameters = protocol.Option(clientMountainParameters)
	}
	if surfaceMaterialAdjustments, ok := serverChunkGeneration.SurfaceMaterialAdjustments.Value(); ok {
		clientSurfaceMaterialAdjustments, err := utils.ConvertSliceWithError(surfaceMaterialAdjustments, c.ConvertServerBiomeElementData)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertServerBiomeChunkGeneration: failed to convert surface material adjustments: %w", err)
		}
		clientChunkGeneration.SurfaceMaterialAdjustments = protocol.Option(clientSurfaceMaterialAdjustments)
	}
	if surfaceMaterials, ok := serverChunkGeneration.SurfaceMaterials.Value(); ok {
		clientSurfaceMaterials, err := c.ConvertServerBiomeSurfaceMaterial(surfaceMaterials)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertServerBiomeChunkGeneration: failed to convert surface materials: %w", err)
		}
		clientChunkGeneration.SurfaceMaterials = protocol.Option(clientSurfaceMaterials)
	}
	if mesaSurface, ok := serverChunkGeneration.MesaSurface.Value(); ok {
		clientMesaSurface, err := c.ConvertServerBiomeMesaSurface(mesaSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertServerBiomeChunkGeneration: failed to convert mesa surface: %w", err)
		}
		clientChunkGeneration.MesaSurface = protocol.Option(clientMesaSurface)
	}
	if cappedSurface, ok := serverChunkGeneration.CappedSurface.Value(); ok {
		clientCappedSurface, err := c.ConvertServerBiomeCappedSurface(cappedSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertServerBiomeChunkGeneration: failed to convert capped surface: %w", err)
		}
		clientChunkGeneration.CappedSurface = protocol.Option(clientCappedSurface)
	}

	return clientChunkGeneration, nil
}

// ConvertServerBiomeMountainParameters converts block runtime IDs inside BiomeMountainParameters (server→client).
func (c *BlockConverter) ConvertServerBiomeMountainParameters(serverParameters protocol.BiomeMountainParameters) (protocol.BiomeMountainParameters, error) {
	clientSteepBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverParameters.SteepBlock)
	if err != nil {
		return protocol.BiomeMountainParameters{}, fmt.Errorf("ConvertServerBiomeMountainParameters: failed to convert steep block: %w", err)
	}
	clientParameters := serverParameters
	clientParameters.SteepBlock = clientSteepBlock
	return clientParameters, nil
}

// ConvertServerBiomeElementData converts block runtime IDs inside BiomeElementData (server→client).
func (c *BlockConverter) ConvertServerBiomeElementData(serverElementData protocol.BiomeElementData) (protocol.BiomeElementData, error) {
	clientAdjustedMaterials, err := c.ConvertServerBiomeSurfaceMaterial(serverElementData.AdjustedMaterials)
	if err != nil {
		return protocol.BiomeElementData{}, fmt.Errorf("ConvertServerBiomeElementData: failed to convert adjusted materials: %w", err)
	}
	clientElementData := serverElementData
	clientElementData.AdjustedMaterials = clientAdjustedMaterials
	return clientElementData, nil
}

// ConvertServerBiomeSurfaceMaterial converts block runtime IDs inside BiomeSurfaceMaterial (server→client).
func (c *BlockConverter) ConvertServerBiomeSurfaceMaterial(serverSurfaceMaterial protocol.BiomeSurfaceMaterial) (protocol.BiomeSurfaceMaterial, error) {
	topBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverSurfaceMaterial.TopBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertServerBiomeSurfaceMaterial: failed to convert top block: %w", err)
	}
	midBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverSurfaceMaterial.MidBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertServerBiomeSurfaceMaterial: failed to convert mid block: %w", err)
	}
	seaFloorBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverSurfaceMaterial.SeaFloorBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertServerBiomeSurfaceMaterial: failed to convert sea floor block: %w", err)
	}
	foundationBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverSurfaceMaterial.FoundationBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertServerBiomeSurfaceMaterial: failed to convert foundation block: %w", err)
	}
	seaBlock, err := c.ConvertServerBlockRuntimeIDInt32(serverSurfaceMaterial.SeaBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertServerBiomeSurfaceMaterial: failed to convert sea block: %w", err)
	}
	clientSurfaceMaterial := serverSurfaceMaterial
	clientSurfaceMaterial.TopBlock = topBlock
	clientSurfaceMaterial.MidBlock = midBlock
	clientSurfaceMaterial.SeaFloorBlock = seaFloorBlock
	clientSurfaceMaterial.FoundationBlock = foundationBlock
	clientSurfaceMaterial.SeaBlock = seaBlock
	return clientSurfaceMaterial, nil
}

// ConvertServerBiomeMesaSurface converts block runtime IDs inside BiomeMesaSurface (server→client).
func (c *BlockConverter) ConvertServerBiomeMesaSurface(serverMesaSurface protocol.BiomeMesaSurface) (protocol.BiomeMesaSurface, error) {
	clayMaterial, err := c.ConvertServerBlockRuntimeID(serverMesaSurface.ClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertServerBiomeMesaSurface: failed to convert clay material: %w", err)
	}
	hardClayMaterial, err := c.ConvertServerBlockRuntimeID(serverMesaSurface.HardClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertServerBiomeMesaSurface: failed to convert hard clay material: %w", err)
	}
	clientMesaSurface := serverMesaSurface
	clientMesaSurface.ClayMaterial = clayMaterial
	clientMesaSurface.HardClayMaterial = hardClayMaterial
	return clientMesaSurface, nil
}

// ConvertServerBiomeCappedSurface converts block runtime IDs inside BiomeCappedSurface (server→client).
func (c *BlockConverter) ConvertServerBiomeCappedSurface(serverCappedSurface protocol.BiomeCappedSurface) (protocol.BiomeCappedSurface, error) {
	floorBlocks, err := utils.ConvertSliceWithError(serverCappedSurface.FloorBlocks, c.ConvertServerBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertServerBiomeCappedSurface: failed to convert floor blocks: %w", err)
	}
	ceilingBlocks, err := utils.ConvertSliceWithError(serverCappedSurface.CeilingBlocks, c.ConvertServerBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertServerBiomeCappedSurface: failed to convert ceiling blocks: %w", err)
	}
	clientCappedSurface := serverCappedSurface
	clientCappedSurface.FloorBlocks = floorBlocks
	clientCappedSurface.CeilingBlocks = ceilingBlocks

	if seaBlock, ok := serverCappedSurface.SeaBlock.Value(); ok {
		clientSeaBlock, err := c.ConvertServerBlockRuntimeID(seaBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertServerBiomeCappedSurface: failed to convert sea block: %w", err)
		}
		clientCappedSurface.SeaBlock = protocol.Option(clientSeaBlock)
	}
	if foundationBlock, ok := serverCappedSurface.FoundationBlock.Value(); ok {
		clientFoundationBlock, err := c.ConvertServerBlockRuntimeID(foundationBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertServerBiomeCappedSurface: failed to convert foundation block: %w", err)
		}
		clientCappedSurface.FoundationBlock = protocol.Option(clientFoundationBlock)
	}
	if beachBlock, ok := serverCappedSurface.BeachBlock.Value(); ok {
		clientBeachBlock, err := c.ConvertServerBlockRuntimeID(beachBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertServerBiomeCappedSurface: failed to convert beach block: %w", err)
		}
		clientCappedSurface.BeachBlock = protocol.Option(clientBeachBlock)
	}

	return clientCappedSurface, nil
}
