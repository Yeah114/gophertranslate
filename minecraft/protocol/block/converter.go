package block

import (
	"fmt"

	minecraft_block "github.com/Yeah114/gophertranslate/minecraft/block"
	"github.com/Yeah114/gophertranslate/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// Converter converts protocol fields that hold block runtime IDs.
type Converter struct {
	bc *minecraft_block.BlockConverter
}

// NewConverter creates a new protocol block converter.
func NewConverter(bc *minecraft_block.BlockConverter) *Converter {
	return &Converter{bc: bc}
}

// BlockConverter returns the underlying block converter.
func (c *Converter) BlockConverter() *minecraft_block.BlockConverter {
	return c.bc
}

// ConvertBlockRuntimeID converts a block runtime ID from the source protocol to the destination protocol.
func (c *Converter) ConvertBlockRuntimeID(srcBlockRuntimeID uint32) (uint32, error) {
	dstBlockRuntimeID, ok := c.bc.ConvertBlockRuntimeID(srcBlockRuntimeID)
	if !ok {
		return 0, fmt.Errorf("ConvertBlockRuntimeID: unknown source block runtime ID %d", srcBlockRuntimeID)
	}
	return dstBlockRuntimeID, nil
}

// ConvertBlockRuntimeIDInt32 converts an int32 block runtime ID.
func (c *Converter) ConvertBlockRuntimeIDInt32(srcBlockRuntimeID int32) (int32, error) {
	if srcBlockRuntimeID < 0 {
		return srcBlockRuntimeID, nil
	}
	dstBlockRuntimeID, err := c.ConvertBlockRuntimeID(uint32(srcBlockRuntimeID))
	if err != nil {
		return 0, err
	}
	return int32(dstBlockRuntimeID), nil
}

// ConvertItemStack converts the block runtime ID inside an ItemStack.
func (c *Converter) ConvertItemStack(srcItemStack protocol.ItemStack) (protocol.ItemStack, error) {
	dstItemStack := srcItemStack
	if srcItemStack.NetworkID == 0 {
		return dstItemStack, nil
	}
	dstBlockRuntimeID, err := c.ConvertBlockRuntimeIDInt32(srcItemStack.BlockRuntimeID)
	if err != nil {
		return protocol.ItemStack{}, fmt.Errorf("ConvertItemStack: failed to convert block runtime ID: %w", err)
	}
	dstItemStack.BlockRuntimeID = dstBlockRuntimeID
	return dstItemStack, nil
}

// ConvertItemInstance converts the block runtime ID inside an ItemInstance.
func (c *Converter) ConvertItemInstance(srcItemInstance protocol.ItemInstance) (protocol.ItemInstance, error) {
	dstStack, err := c.ConvertItemStack(srcItemInstance.Stack)
	if err != nil {
		return protocol.ItemInstance{}, fmt.Errorf("ConvertItemInstance: failed to convert stack: %w", err)
	}
	dstItemInstance := srcItemInstance
	dstItemInstance.Stack = dstStack
	return dstItemInstance, nil
}

// ConvertCreativeGroup converts the block runtime ID inside a CreativeGroup icon.
func (c *Converter) ConvertCreativeGroup(srcCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error) {
	icon, err := c.ConvertItemStack(srcCreativeGroup.Icon)
	if err != nil {
		return protocol.CreativeGroup{}, fmt.Errorf("ConvertCreativeGroup: failed to convert icon: %w", err)
	}
	dstCreativeGroup := srcCreativeGroup
	dstCreativeGroup.Icon = icon
	return dstCreativeGroup, nil
}

// ConvertCreativeItem converts the block runtime ID inside a CreativeItem.
func (c *Converter) ConvertCreativeItem(srcCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error) {
	item, err := c.ConvertItemStack(srcCreativeItem.Item)
	if err != nil {
		return protocol.CreativeItem{}, fmt.Errorf("ConvertCreativeItem: failed to convert item: %w", err)
	}
	dstCreativeItem := srcCreativeItem
	dstCreativeItem.Item = item
	return dstCreativeItem, nil
}

// ConvertInventoryAction converts item block runtime IDs inside an InventoryAction.
func (c *Converter) ConvertInventoryAction(srcInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error) {
	dstOldItem, err := c.ConvertItemInstance(srcInventoryAction.OldItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertInventoryAction: failed to convert old item: %w", err)
	}
	dstNewItem, err := c.ConvertItemInstance(srcInventoryAction.NewItem)
	if err != nil {
		return protocol.InventoryAction{}, fmt.Errorf("ConvertInventoryAction: failed to convert new item: %w", err)
	}
	dstInventoryAction := srcInventoryAction
	dstInventoryAction.OldItem = dstOldItem
	dstInventoryAction.NewItem = dstNewItem
	return dstInventoryAction, nil
}

// ConvertBlockChangeEntry converts the block runtime ID inside a BlockChangeEntry.
func (c *Converter) ConvertBlockChangeEntry(srcBlockChangeEntry protocol.BlockChangeEntry) (protocol.BlockChangeEntry, error) {
	dstBlockRuntimeID, err := c.ConvertBlockRuntimeID(srcBlockChangeEntry.BlockRuntimeID)
	if err != nil {
		return protocol.BlockChangeEntry{}, fmt.Errorf("ConvertBlockChangeEntry: failed to convert block runtime ID: %w", err)
	}
	dstBlockChangeEntry := srcBlockChangeEntry
	dstBlockChangeEntry.BlockRuntimeID = dstBlockRuntimeID
	return dstBlockChangeEntry, nil
}

// ConvertUseItemTransactionData converts item and block runtime IDs inside UseItemTransactionData.
func (c *Converter) ConvertUseItemTransactionData(srcData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error) {
	actions, err := utils.ConvertSliceWithError(srcData.Actions, c.ConvertInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertUseItemTransactionData: failed to convert actions: %w", err)
	}
	heldItem, err := c.ConvertItemInstance(srcData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertUseItemTransactionData: failed to convert held item: %w", err)
	}
	blockRuntimeID, err := c.ConvertBlockRuntimeID(srcData.BlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUseItemTransactionData: failed to convert block runtime ID: %w", err)
	}
	dstData := *srcData
	dstData.Actions = actions
	dstData.HeldItem = heldItem
	dstData.BlockRuntimeID = blockRuntimeID
	return &dstData, nil
}

// ConvertUseItemOnEntityTransactionData converts item block runtime IDs inside UseItemOnEntityTransactionData.
func (c *Converter) ConvertUseItemOnEntityTransactionData(srcData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error) {
	heldItem, err := c.ConvertItemInstance(srcData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertUseItemOnEntityTransactionData: failed to convert held item: %w", err)
	}
	dstData := *srcData
	dstData.HeldItem = heldItem
	return &dstData, nil
}

// ConvertReleaseItemTransactionData converts item block runtime IDs inside ReleaseItemTransactionData.
func (c *Converter) ConvertReleaseItemTransactionData(srcData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error) {
	heldItem, err := c.ConvertItemInstance(srcData.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertReleaseItemTransactionData: failed to convert held item: %w", err)
	}
	dstData := *srcData
	dstData.HeldItem = heldItem
	return &dstData, nil
}

// ConvertWaxedOrUnwaxedCopperEvent converts the copper block runtime ID inside WaxedOrUnwaxedCopperEvent.
func (c *Converter) ConvertWaxedOrUnwaxedCopperEvent(srcEvent *protocol.WaxedOrUnwaxedCopperEvent) (*protocol.WaxedOrUnwaxedCopperEvent, error) {
	dstCopperBlockID, err := c.ConvertBlockRuntimeIDInt32(srcEvent.CopperBlockID)
	if err != nil {
		return nil, fmt.Errorf("ConvertWaxedOrUnwaxedCopperEvent: failed to convert copper block ID: %w", err)
	}
	dstEvent := *srcEvent
	dstEvent.CopperBlockID = dstCopperBlockID
	return &dstEvent, nil
}

// ConvertBiomeDefinition converts block runtime IDs inside a BiomeDefinition.
func (c *Converter) ConvertBiomeDefinition(srcBiomeDefinition protocol.BiomeDefinition) (protocol.BiomeDefinition, error) {
	chunkGeneration, ok := srcBiomeDefinition.ChunkGeneration.Value()
	if !ok {
		return srcBiomeDefinition, nil
	}
	dstChunkGeneration, err := c.ConvertBiomeChunkGeneration(chunkGeneration)
	if err != nil {
		return protocol.BiomeDefinition{}, fmt.Errorf("ConvertBiomeDefinition: failed to convert chunk generation: %w", err)
	}
	dstBiomeDefinition := srcBiomeDefinition
	dstBiomeDefinition.ChunkGeneration = protocol.Option(dstChunkGeneration)
	return dstBiomeDefinition, nil
}

// ConvertBiomeChunkGeneration converts block runtime IDs inside BiomeChunkGeneration.
func (c *Converter) ConvertBiomeChunkGeneration(srcChunkGeneration protocol.BiomeChunkGeneration) (protocol.BiomeChunkGeneration, error) {
	dstChunkGeneration := srcChunkGeneration

	if mountainParameters, ok := srcChunkGeneration.MountainParameters.Value(); ok {
		dstMountainParameters, err := c.ConvertBiomeMountainParameters(mountainParameters)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertBiomeChunkGeneration: failed to convert mountain parameters: %w", err)
		}
		dstChunkGeneration.MountainParameters = protocol.Option(dstMountainParameters)
	}
	if surfaceMaterialAdjustments, ok := srcChunkGeneration.SurfaceMaterialAdjustments.Value(); ok {
		dstSurfaceMaterialAdjustments, err := utils.ConvertSliceWithError(surfaceMaterialAdjustments, c.ConvertBiomeElementData)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertBiomeChunkGeneration: failed to convert surface material adjustments: %w", err)
		}
		dstChunkGeneration.SurfaceMaterialAdjustments = protocol.Option(dstSurfaceMaterialAdjustments)
	}
	if surfaceMaterials, ok := srcChunkGeneration.SurfaceMaterials.Value(); ok {
		dstSurfaceMaterials, err := c.ConvertBiomeSurfaceMaterial(surfaceMaterials)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertBiomeChunkGeneration: failed to convert surface materials: %w", err)
		}
		dstChunkGeneration.SurfaceMaterials = protocol.Option(dstSurfaceMaterials)
	}
	if mesaSurface, ok := srcChunkGeneration.MesaSurface.Value(); ok {
		dstMesaSurface, err := c.ConvertBiomeMesaSurface(mesaSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertBiomeChunkGeneration: failed to convert mesa surface: %w", err)
		}
		dstChunkGeneration.MesaSurface = protocol.Option(dstMesaSurface)
	}
	if cappedSurface, ok := srcChunkGeneration.CappedSurface.Value(); ok {
		dstCappedSurface, err := c.ConvertBiomeCappedSurface(cappedSurface)
		if err != nil {
			return protocol.BiomeChunkGeneration{}, fmt.Errorf("ConvertBiomeChunkGeneration: failed to convert capped surface: %w", err)
		}
		dstChunkGeneration.CappedSurface = protocol.Option(dstCappedSurface)
	}

	return dstChunkGeneration, nil
}

// ConvertBiomeMountainParameters converts block runtime IDs inside BiomeMountainParameters.
func (c *Converter) ConvertBiomeMountainParameters(srcParameters protocol.BiomeMountainParameters) (protocol.BiomeMountainParameters, error) {
	dstSteepBlock, err := c.ConvertBlockRuntimeIDInt32(srcParameters.SteepBlock)
	if err != nil {
		return protocol.BiomeMountainParameters{}, fmt.Errorf("ConvertBiomeMountainParameters: failed to convert steep block: %w", err)
	}
	dstParameters := srcParameters
	dstParameters.SteepBlock = dstSteepBlock
	return dstParameters, nil
}

// ConvertBiomeElementData converts block runtime IDs inside BiomeElementData.
func (c *Converter) ConvertBiomeElementData(srcElementData protocol.BiomeElementData) (protocol.BiomeElementData, error) {
	dstAdjustedMaterials, err := c.ConvertBiomeSurfaceMaterial(srcElementData.AdjustedMaterials)
	if err != nil {
		return protocol.BiomeElementData{}, fmt.Errorf("ConvertBiomeElementData: failed to convert adjusted materials: %w", err)
	}
	dstElementData := srcElementData
	dstElementData.AdjustedMaterials = dstAdjustedMaterials
	return dstElementData, nil
}

// ConvertBiomeSurfaceMaterial converts block runtime IDs inside BiomeSurfaceMaterial.
func (c *Converter) ConvertBiomeSurfaceMaterial(srcSurfaceMaterial protocol.BiomeSurfaceMaterial) (protocol.BiomeSurfaceMaterial, error) {
	topBlock, err := c.ConvertBlockRuntimeIDInt32(srcSurfaceMaterial.TopBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertBiomeSurfaceMaterial: failed to convert top block: %w", err)
	}
	midBlock, err := c.ConvertBlockRuntimeIDInt32(srcSurfaceMaterial.MidBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertBiomeSurfaceMaterial: failed to convert mid block: %w", err)
	}
	seaFloorBlock, err := c.ConvertBlockRuntimeIDInt32(srcSurfaceMaterial.SeaFloorBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertBiomeSurfaceMaterial: failed to convert sea floor block: %w", err)
	}
	foundationBlock, err := c.ConvertBlockRuntimeIDInt32(srcSurfaceMaterial.FoundationBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertBiomeSurfaceMaterial: failed to convert foundation block: %w", err)
	}
	seaBlock, err := c.ConvertBlockRuntimeIDInt32(srcSurfaceMaterial.SeaBlock)
	if err != nil {
		return protocol.BiomeSurfaceMaterial{}, fmt.Errorf("ConvertBiomeSurfaceMaterial: failed to convert sea block: %w", err)
	}
	dstSurfaceMaterial := srcSurfaceMaterial
	dstSurfaceMaterial.TopBlock = topBlock
	dstSurfaceMaterial.MidBlock = midBlock
	dstSurfaceMaterial.SeaFloorBlock = seaFloorBlock
	dstSurfaceMaterial.FoundationBlock = foundationBlock
	dstSurfaceMaterial.SeaBlock = seaBlock
	return dstSurfaceMaterial, nil
}

// ConvertBiomeMesaSurface converts block runtime IDs inside BiomeMesaSurface.
func (c *Converter) ConvertBiomeMesaSurface(srcMesaSurface protocol.BiomeMesaSurface) (protocol.BiomeMesaSurface, error) {
	clayMaterial, err := c.ConvertBlockRuntimeID(srcMesaSurface.ClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertBiomeMesaSurface: failed to convert clay material: %w", err)
	}
	hardClayMaterial, err := c.ConvertBlockRuntimeID(srcMesaSurface.HardClayMaterial)
	if err != nil {
		return protocol.BiomeMesaSurface{}, fmt.Errorf("ConvertBiomeMesaSurface: failed to convert hard clay material: %w", err)
	}
	dstMesaSurface := srcMesaSurface
	dstMesaSurface.ClayMaterial = clayMaterial
	dstMesaSurface.HardClayMaterial = hardClayMaterial
	return dstMesaSurface, nil
}

// ConvertBiomeCappedSurface converts block runtime IDs inside BiomeCappedSurface.
func (c *Converter) ConvertBiomeCappedSurface(srcCappedSurface protocol.BiomeCappedSurface) (protocol.BiomeCappedSurface, error) {
	floorBlocks, err := utils.ConvertSliceWithError(srcCappedSurface.FloorBlocks, c.ConvertBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertBiomeCappedSurface: failed to convert floor blocks: %w", err)
	}
	ceilingBlocks, err := utils.ConvertSliceWithError(srcCappedSurface.CeilingBlocks, c.ConvertBlockRuntimeIDInt32)
	if err != nil {
		return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertBiomeCappedSurface: failed to convert ceiling blocks: %w", err)
	}
	dstCappedSurface := srcCappedSurface
	dstCappedSurface.FloorBlocks = floorBlocks
	dstCappedSurface.CeilingBlocks = ceilingBlocks

	if seaBlock, ok := srcCappedSurface.SeaBlock.Value(); ok {
		dstSeaBlock, err := c.ConvertBlockRuntimeID(seaBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertBiomeCappedSurface: failed to convert sea block: %w", err)
		}
		dstCappedSurface.SeaBlock = protocol.Option(dstSeaBlock)
	}
	if foundationBlock, ok := srcCappedSurface.FoundationBlock.Value(); ok {
		dstFoundationBlock, err := c.ConvertBlockRuntimeID(foundationBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertBiomeCappedSurface: failed to convert foundation block: %w", err)
		}
		dstCappedSurface.FoundationBlock = protocol.Option(dstFoundationBlock)
	}
	if beachBlock, ok := srcCappedSurface.BeachBlock.Value(); ok {
		dstBeachBlock, err := c.ConvertBlockRuntimeID(beachBlock)
		if err != nil {
			return protocol.BiomeCappedSurface{}, fmt.Errorf("ConvertBiomeCappedSurface: failed to convert beach block: %w", err)
		}
		dstCappedSurface.BeachBlock = protocol.Option(dstBeachBlock)
	}

	return dstCappedSurface, nil
}
