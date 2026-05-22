package protocol

import (
	"fmt"

	"github.com/Yeah114/bedrock-world-operator/define"
	minecraft_block "github.com/Yeah114/gophertranslate/minecraft/block"
	minecraft_chunk "github.com/Yeah114/gophertranslate/minecraft/chunk"
	protocol_block "github.com/Yeah114/gophertranslate/minecraft/protocol/block"
	protocol_chunk "github.com/Yeah114/gophertranslate/minecraft/protocol/chunk"
	"github.com/Yeah114/gophertranslate/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// Converter converts protocol packets and protocol fields between two Minecraft protocol versions.
type Converter struct {
	bc *protocol_block.Converter
	cc *protocol_chunk.ChunkConverter
}

// NewConverter creates a new protocol converter.
func NewConverter(
	bc *minecraft_block.BlockConverter,
	ranges map[int32]define.Range,
	currentDimension int32,
) *Converter {
	return &Converter{
		bc: protocol_block.NewConverter(bc),
		cc: protocol_chunk.NewChunkConverter(minecraft_chunk.NewChunkConverter(bc), ranges, currentDimension),
	}
}

// BlockConverter returns the block protocol converter used by the converter.
func (c *Converter) BlockConverter() *protocol_block.Converter {
	return c.bc
}

// ChunkConverter returns the chunk protocol converter used by the converter.
func (c *Converter) ChunkConverter() *protocol_chunk.ChunkConverter {
	return c.cc
}

// SetCurrentDimension sets the dimension used for dimension-less cache blobs.
func (c *Converter) SetCurrentDimension(dimension int32) {
	c.cc.CurrentDimension = dimension
}

// ConvertPacket converts block runtime IDs in a supported packet.
func (c *Converter) ConvertPacket(srcPacket packet.Packet) (packet.Packet, error) {
	switch pk := srcPacket.(type) {
	case *packet.UpdateBlock:
		return c.ConvertUpdateBlock(pk)
	case *packet.UpdateBlockSynced:
		return c.ConvertUpdateBlockSynced(pk)
	case *packet.UpdateSubChunkBlocks:
		return c.ConvertUpdateSubChunkBlocks(pk)
	case *packet.LevelChunk:
		return c.ConvertLevelChunk(pk)
	case *packet.SubChunk:
		return c.ConvertSubChunk(pk)
	case *packet.ClientCacheMissResponse:
		return c.ConvertClientCacheMissResponse(pk)
	case *packet.AddItemActor:
		return c.ConvertAddItemActor(pk)
	case *packet.AddPlayer:
		return c.ConvertAddPlayer(pk)
	case *packet.CreativeContent:
		return c.ConvertCreativeContent(pk)
	case *packet.InventoryContent:
		return c.ConvertInventoryContent(pk)
	case *packet.InventorySlot:
		return c.ConvertInventorySlot(pk)
	case *packet.InventoryTransaction:
		return c.ConvertInventoryTransaction(pk)
	case *packet.MobArmourEquipment:
		return c.ConvertMobArmourEquipment(pk)
	case *packet.MobEquipment:
		return c.ConvertMobEquipment(pk)
	case *packet.BiomeDefinitionList:
		return c.ConvertBiomeDefinitionList(pk)
	case *packet.Event:
		return c.ConvertEvent(pk)
	case *packet.StartGame:
		return c.ConvertStartGame(pk)
	default:
		return srcPacket, nil
	}
}

// ConvertUpdateBlock converts an UpdateBlock packet.
func (c *Converter) ConvertUpdateBlock(srcUpdateBlock *packet.UpdateBlock) (*packet.UpdateBlock, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(srcUpdateBlock.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlock: failed to convert new block runtime ID: %w", err)
	}
	dstUpdateBlock := *srcUpdateBlock
	dstUpdateBlock.NewBlockRuntimeID = newBlockRuntimeID
	return &dstUpdateBlock, nil
}

// ConvertUpdateBlockSynced converts an UpdateBlockSynced packet.
func (c *Converter) ConvertUpdateBlockSynced(srcUpdateBlockSynced *packet.UpdateBlockSynced) (*packet.UpdateBlockSynced, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(srcUpdateBlockSynced.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlockSynced: failed to convert new block runtime ID: %w", err)
	}
	dstUpdateBlockSynced := *srcUpdateBlockSynced
	dstUpdateBlockSynced.NewBlockRuntimeID = newBlockRuntimeID
	return &dstUpdateBlockSynced, nil
}

// ConvertUpdateSubChunkBlocks converts an UpdateSubChunkBlocks packet.
func (c *Converter) ConvertUpdateSubChunkBlocks(srcUpdateSubChunkBlocks *packet.UpdateSubChunkBlocks) (*packet.UpdateSubChunkBlocks, error) {
	blocks, err := utils.ConvertSliceWithError(srcUpdateSubChunkBlocks.Blocks, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert blocks: %w", err)
	}
	extra, err := utils.ConvertSliceWithError(srcUpdateSubChunkBlocks.Extra, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert extra blocks: %w", err)
	}
	dstUpdateSubChunkBlocks := *srcUpdateSubChunkBlocks
	dstUpdateSubChunkBlocks.Blocks = blocks
	dstUpdateSubChunkBlocks.Extra = extra
	return &dstUpdateSubChunkBlocks, nil
}

// ConvertLevelChunk converts a LevelChunk packet.
func (c *Converter) ConvertLevelChunk(srcLevelChunk *packet.LevelChunk) (*packet.LevelChunk, error) {
	return c.cc.ConvertLevelChunk(srcLevelChunk)
}

// ConvertSubChunk converts a SubChunk packet.
func (c *Converter) ConvertSubChunk(srcSubChunk *packet.SubChunk) (*packet.SubChunk, error) {
	return c.cc.ConvertSubChunk(srcSubChunk)
}

// ConvertClientCacheMissResponse converts a ClientCacheMissResponse packet.
func (c *Converter) ConvertClientCacheMissResponse(srcClientCacheMissResponse *packet.ClientCacheMissResponse) (*packet.ClientCacheMissResponse, error) {
	return c.cc.ConvertClientCacheMissResponse(srcClientCacheMissResponse)
}

// ConvertAddItemActor converts item block runtime IDs inside an AddItemActor packet.
func (c *Converter) ConvertAddItemActor(srcAddItemActor *packet.AddItemActor) (*packet.AddItemActor, error) {
	item, err := c.bc.ConvertItemInstance(srcAddItemActor.Item)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddItemActor: failed to convert item: %w", err)
	}
	dstAddItemActor := *srcAddItemActor
	dstAddItemActor.Item = item
	return &dstAddItemActor, nil
}

// ConvertAddPlayer converts item block runtime IDs inside an AddPlayer packet.
func (c *Converter) ConvertAddPlayer(srcAddPlayer *packet.AddPlayer) (*packet.AddPlayer, error) {
	heldItem, err := c.bc.ConvertItemInstance(srcAddPlayer.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddPlayer: failed to convert held item: %w", err)
	}
	dstAddPlayer := *srcAddPlayer
	dstAddPlayer.HeldItem = heldItem
	dstAddPlayer.EntityLinks = append([]protocol.EntityLink{}, srcAddPlayer.EntityLinks...)
	return &dstAddPlayer, nil
}

// ConvertCreativeContent converts item block runtime IDs inside a CreativeContent packet.
func (c *Converter) ConvertCreativeContent(srcCreativeContent *packet.CreativeContent) (*packet.CreativeContent, error) {
	groups, err := utils.ConvertSliceWithError(srcCreativeContent.Groups, c.bc.ConvertCreativeGroup)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert groups: %w", err)
	}
	items, err := utils.ConvertSliceWithError(srcCreativeContent.Items, c.bc.ConvertCreativeItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert items: %w", err)
	}
	dstCreativeContent := *srcCreativeContent
	dstCreativeContent.Groups = groups
	dstCreativeContent.Items = items
	return &dstCreativeContent, nil
}

// ConvertInventoryContent converts item block runtime IDs inside an InventoryContent packet.
func (c *Converter) ConvertInventoryContent(srcInventoryContent *packet.InventoryContent) (*packet.InventoryContent, error) {
	content, err := utils.ConvertSliceWithError(srcInventoryContent.Content, c.bc.ConvertItemInstance)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryContent: failed to convert content: %w", err)
	}
	storageItem, err := c.bc.ConvertItemInstance(srcInventoryContent.StorageItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryContent: failed to convert storage item: %w", err)
	}
	dstInventoryContent := *srcInventoryContent
	dstInventoryContent.Content = content
	dstInventoryContent.StorageItem = storageItem
	return &dstInventoryContent, nil
}

// ConvertInventorySlot converts item block runtime IDs inside an InventorySlot packet.
func (c *Converter) ConvertInventorySlot(srcInventorySlot *packet.InventorySlot) (*packet.InventorySlot, error) {
	newItem, err := c.bc.ConvertItemInstance(srcInventorySlot.NewItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventorySlot: failed to convert new item: %w", err)
	}
	dstInventorySlot := *srcInventorySlot
	dstInventorySlot.NewItem = newItem
	if storageItem, ok := srcInventorySlot.StorageItem.Value(); ok {
		dstStorageItem, err := c.bc.ConvertItemInstance(storageItem)
		if err != nil {
			return nil, fmt.Errorf("ConvertInventorySlot: failed to convert storage item: %w", err)
		}
		dstInventorySlot.StorageItem = protocol.Option(dstStorageItem)
	}
	return &dstInventorySlot, nil
}

// ConvertInventoryTransaction converts block runtime IDs inside an InventoryTransaction packet.
func (c *Converter) ConvertInventoryTransaction(srcInventoryTransaction *packet.InventoryTransaction) (*packet.InventoryTransaction, error) {
	actions, err := utils.ConvertSliceWithError(srcInventoryTransaction.Actions, c.bc.ConvertInventoryAction)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryTransaction: failed to convert actions: %w", err)
	}
	transactionData, err := c.ConvertInventoryTransactionData(srcInventoryTransaction.TransactionData)
	if err != nil {
		return nil, fmt.Errorf("ConvertInventoryTransaction: failed to convert transaction data: %w", err)
	}
	dstInventoryTransaction := *srcInventoryTransaction
	dstInventoryTransaction.Actions = actions
	dstInventoryTransaction.TransactionData = transactionData
	return &dstInventoryTransaction, nil
}

// ConvertMobArmourEquipment converts item block runtime IDs inside a MobArmourEquipment packet.
func (c *Converter) ConvertMobArmourEquipment(srcMobArmourEquipment *packet.MobArmourEquipment) (*packet.MobArmourEquipment, error) {
	helmet, err := c.bc.ConvertItemInstance(srcMobArmourEquipment.Helmet)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert helmet: %w", err)
	}
	chestplate, err := c.bc.ConvertItemInstance(srcMobArmourEquipment.Chestplate)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert chestplate: %w", err)
	}
	leggings, err := c.bc.ConvertItemInstance(srcMobArmourEquipment.Leggings)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert leggings: %w", err)
	}
	boots, err := c.bc.ConvertItemInstance(srcMobArmourEquipment.Boots)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert boots: %w", err)
	}
	body, err := c.bc.ConvertItemInstance(srcMobArmourEquipment.Body)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert body: %w", err)
	}
	dstMobArmourEquipment := *srcMobArmourEquipment
	dstMobArmourEquipment.Helmet = helmet
	dstMobArmourEquipment.Chestplate = chestplate
	dstMobArmourEquipment.Leggings = leggings
	dstMobArmourEquipment.Boots = boots
	dstMobArmourEquipment.Body = body
	return &dstMobArmourEquipment, nil
}

// ConvertMobEquipment converts item block runtime IDs inside a MobEquipment packet.
func (c *Converter) ConvertMobEquipment(srcMobEquipment *packet.MobEquipment) (*packet.MobEquipment, error) {
	newItem, err := c.bc.ConvertItemInstance(srcMobEquipment.NewItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobEquipment: failed to convert new item: %w", err)
	}
	dstMobEquipment := *srcMobEquipment
	dstMobEquipment.NewItem = newItem
	return &dstMobEquipment, nil
}

// ConvertInventoryTransactionData converts block runtime IDs inside inventory transaction data.
func (c *Converter) ConvertInventoryTransactionData(srcData protocol.InventoryTransactionData) (protocol.InventoryTransactionData, error) {
	switch data := srcData.(type) {
	case nil:
		return &protocol.NormalTransactionData{}, nil
	case *protocol.NormalTransactionData:
		dstData := *data
		return &dstData, nil
	case *protocol.MismatchTransactionData:
		dstData := *data
		return &dstData, nil
	case *protocol.UseItemTransactionData:
		return c.bc.ConvertUseItemTransactionData(data)
	case *protocol.UseItemOnEntityTransactionData:
		return c.bc.ConvertUseItemOnEntityTransactionData(data)
	case *protocol.ReleaseItemTransactionData:
		return c.bc.ConvertReleaseItemTransactionData(data)
	default:
		return srcData, nil
	}
}

// ConvertBiomeDefinitionList converts block runtime IDs inside a BiomeDefinitionList packet.
func (c *Converter) ConvertBiomeDefinitionList(srcBiomeDefinitionList *packet.BiomeDefinitionList) (*packet.BiomeDefinitionList, error) {
	biomeDefinitions, err := utils.ConvertSliceWithError(srcBiomeDefinitionList.BiomeDefinitions, c.bc.ConvertBiomeDefinition)
	if err != nil {
		return nil, fmt.Errorf("ConvertBiomeDefinitionList: failed to convert biome definitions: %w", err)
	}
	dstBiomeDefinitionList := *srcBiomeDefinitionList
	dstBiomeDefinitionList.BiomeDefinitions = biomeDefinitions
	dstBiomeDefinitionList.StringList = append([]string{}, srcBiomeDefinitionList.StringList...)
	return &dstBiomeDefinitionList, nil
}

// ConvertEvent converts block runtime IDs inside an Event packet.
func (c *Converter) ConvertEvent(srcEvent *packet.Event) (*packet.Event, error) {
	eventData, err := c.ConvertEventData(srcEvent.Event)
	if err != nil {
		return nil, fmt.Errorf("ConvertEvent: failed to convert event data: %w", err)
	}
	dstEvent := *srcEvent
	dstEvent.Event = eventData
	return &dstEvent, nil
}

// ConvertEventData converts block runtime IDs inside event data.
func (c *Converter) ConvertEventData(srcEvent protocol.Event) (protocol.Event, error) {
	switch event := srcEvent.(type) {
	case nil:
		return nil, nil
	case *protocol.WaxedOrUnwaxedCopperEvent:
		return c.bc.ConvertWaxedOrUnwaxedCopperEvent(event)
	default:
		return srcEvent, nil
	}
}

// ConvertStartGame adjusts StartGame block runtime ID compatibility fields.
func (c *Converter) ConvertStartGame(srcStartGame *packet.StartGame) (*packet.StartGame, error) {
	dstStartGame := *srcStartGame
	dstStartGame.UseBlockNetworkIDHashes = c.bc.BlockConverter().DstTable().UseNetworkIDHashes()
	return &dstStartGame, nil
}
