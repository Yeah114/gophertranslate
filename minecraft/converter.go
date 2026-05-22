package v1v26v10

import (
	"fmt"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gopherconvert/define"
	minecraft_block "github.com/Yeah114/gopherconvert/minecraft/block"
	minecraft_chunk "github.com/Yeah114/gopherconvert/minecraft/chunk"
	protocol_block "github.com/Yeah114/gopherconvert/minecraft/protocol/block"
	"github.com/Yeah114/gopherconvert/minecraft/protocol/chunk"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// MinecraftConverter converts Minecraft packets between two connections.
type MinecraftConverter struct {
	srcConn  define.Conn
	dstConn  define.Conn
	srcTable *bwo_block.BlockRuntimeIDTable
	dstTable *bwo_block.BlockRuntimeIDTable
	bc       *protocol_block.BlockConverter
	cc       *chunk.ChunkConverter
}

// NewMinecraftConverter creates a converter that translates packets from srcConn to dstConn.
func NewMinecraftConverter(srcConn define.Conn, dstConn define.Conn) (*MinecraftConverter, error) {
	srcTable, err := utils.BlockRuntimeIDTableFromGameData(srcConn.GameData())
	if err != nil {
		return nil, fmt.Errorf("NewMinecraftConverter: failed to create block runtime ID table for source connection: %w", err)
	}

	dstTable, err := utils.BlockRuntimeIDTableFromGameData(dstConn.GameData())
	if err != nil {
		return nil, fmt.Errorf("NewMinecraftConverter: failed to create block runtime ID table for destination connection: %w", err)
	}

	srcGameData := srcConn.GameData()
	blockConverter := minecraft_block.NewBlockConverter(
		srcConn.Proto().ID(),
		srcTable,
		dstConn.Proto().ID(),
		dstTable,
	)
	bc := protocol_block.NewBlockConverter(blockConverter)
	cc := chunk.NewChunkConverter(
		minecraft_chunk.NewChunkConverter(blockConverter),
		utils.RangesFromGameData(srcGameData),
		srcGameData.Dimension,
	)

	return &MinecraftConverter{
		srcConn:  srcConn,
		dstConn:  dstConn,
		srcTable: srcTable,
		dstTable: dstTable,
		bc:       bc,
		cc:       cc,
	}, nil
}

// BlockConverter returns the protocol block converter used by the converter.
func (c *MinecraftConverter) BlockConverter() *protocol_block.BlockConverter {
	return c.bc
}

// ChunkConverter returns the protocol chunk converter used by the converter.
func (c *MinecraftConverter) ChunkConverter() *chunk.ChunkConverter {
	return c.cc
}

// ConvertPacket converts block runtime IDs in a supported packet.
func (c *MinecraftConverter) ConvertPacket(pk packet.Packet) (packet.Packet, error) {
	if pk == nil {
		return nil, fmt.Errorf("ConvertPacket: packet is nil")
	}
	switch pkt := pk.(type) {
	case *packet.UpdateBlock:
		return c.ConvertUpdateBlock(pkt)
	case *packet.UpdateBlockSynced:
		return c.ConvertUpdateBlockSynced(pkt)
	case *packet.UpdateSubChunkBlocks:
		return c.ConvertUpdateSubChunkBlocks(pkt)
	case *packet.LevelChunk:
		return c.ConvertLevelChunk(pkt)
	case *packet.SubChunk:
		return c.ConvertSubChunk(pkt)
	case *packet.ClientCacheMissResponse:
		return c.ConvertClientCacheMissResponse(pkt)
	case *packet.AddItemActor:
		return c.ConvertAddItemActor(pkt)
	case *packet.AddPlayer:
		return c.ConvertAddPlayer(pkt)
	case *packet.CreativeContent:
		return c.ConvertCreativeContent(pkt)
	case *packet.InventoryContent:
		return c.ConvertInventoryContent(pkt)
	case *packet.InventorySlot:
		return c.ConvertInventorySlot(pkt)
	case *packet.InventoryTransaction:
		return c.ConvertInventoryTransaction(pkt)
	case *packet.MobArmourEquipment:
		return c.ConvertMobArmourEquipment(pkt)
	case *packet.MobEquipment:
		return c.ConvertMobEquipment(pkt)
	case *packet.BiomeDefinitionList:
		return c.ConvertBiomeDefinitionList(pkt)
	case *packet.Event:
		return c.ConvertEvent(pkt)
	case *packet.StartGame:
		return c.ConvertStartGame(pkt)
	default:
		return pk, nil
	}
}

// ConvertUpdateBlock converts an UpdateBlock packet.
func (c *MinecraftConverter) ConvertUpdateBlock(pk *packet.UpdateBlock) (*packet.UpdateBlock, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlock: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return &dst, nil
}

// ConvertUpdateBlockSynced converts an UpdateBlockSynced packet.
func (c *MinecraftConverter) ConvertUpdateBlockSynced(pk *packet.UpdateBlockSynced) (*packet.UpdateBlockSynced, error) {
	newBlockRuntimeID, err := c.bc.ConvertBlockRuntimeID(pk.NewBlockRuntimeID)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateBlockSynced: failed to convert new block runtime ID: %w", err)
	}
	dst := *pk
	dst.NewBlockRuntimeID = newBlockRuntimeID
	return &dst, nil
}

// ConvertUpdateSubChunkBlocks converts an UpdateSubChunkBlocks packet.
func (c *MinecraftConverter) ConvertUpdateSubChunkBlocks(pk *packet.UpdateSubChunkBlocks) (*packet.UpdateSubChunkBlocks, error) {
	blocks, err := utils.ConvertSliceWithError(pk.Blocks, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert blocks: %w", err)
	}
	extra, err := utils.ConvertSliceWithError(pk.Extra, c.bc.ConvertBlockChangeEntry)
	if err != nil {
		return nil, fmt.Errorf("ConvertUpdateSubChunkBlocks: failed to convert extra blocks: %w", err)
	}
	dst := *pk
	dst.Blocks = blocks
	dst.Extra = extra
	return &dst, nil
}

// ConvertLevelChunk converts a LevelChunk packet.
func (c *MinecraftConverter) ConvertLevelChunk(pk *packet.LevelChunk) (*packet.LevelChunk, error) {
	return c.cc.ConvertLevelChunk(pk)
}

// ConvertSubChunk converts a SubChunk packet.
func (c *MinecraftConverter) ConvertSubChunk(pk *packet.SubChunk) (*packet.SubChunk, error) {
	return c.cc.ConvertSubChunk(pk)
}

// ConvertClientCacheMissResponse converts a ClientCacheMissResponse packet.
func (c *MinecraftConverter) ConvertClientCacheMissResponse(pk *packet.ClientCacheMissResponse) (*packet.ClientCacheMissResponse, error) {
	return c.cc.ConvertClientCacheMissResponse(pk)
}

// ConvertAddItemActor converts item block runtime IDs inside an AddItemActor packet.
func (c *MinecraftConverter) ConvertAddItemActor(pk *packet.AddItemActor) (*packet.AddItemActor, error) {
	item, err := c.bc.ConvertItemInstance(pk.Item)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddItemActor: failed to convert item: %w", err)
	}
	dst := *pk
	dst.Item = item
	return &dst, nil
}

// ConvertAddPlayer converts item block runtime IDs inside an AddPlayer packet.
func (c *MinecraftConverter) ConvertAddPlayer(pk *packet.AddPlayer) (*packet.AddPlayer, error) {
	heldItem, err := c.bc.ConvertItemInstance(pk.HeldItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertAddPlayer: failed to convert held item: %w", err)
	}
	dst := *pk
	dst.HeldItem = heldItem
	return &dst, nil
}

// ConvertCreativeContent converts item block runtime IDs inside a CreativeContent packet.
func (c *MinecraftConverter) ConvertCreativeContent(pk *packet.CreativeContent) (*packet.CreativeContent, error) {
	groups, err := utils.ConvertSliceWithError(pk.Groups, c.bc.ConvertCreativeGroup)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert groups: %w", err)
	}
	items, err := utils.ConvertSliceWithError(pk.Items, c.bc.ConvertCreativeItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertCreativeContent: failed to convert items: %w", err)
	}
	dst := *pk
	dst.Groups = groups
	dst.Items = items
	return &dst, nil
}

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

// ConvertMobArmourEquipment converts item block runtime IDs inside a MobArmourEquipment packet.
func (c *MinecraftConverter) ConvertMobArmourEquipment(pk *packet.MobArmourEquipment) (*packet.MobArmourEquipment, error) {
	helmet, err := c.bc.ConvertItemInstance(pk.Helmet)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert helmet: %w", err)
	}
	chestplate, err := c.bc.ConvertItemInstance(pk.Chestplate)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert chestplate: %w", err)
	}
	leggings, err := c.bc.ConvertItemInstance(pk.Leggings)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert leggings: %w", err)
	}
	boots, err := c.bc.ConvertItemInstance(pk.Boots)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert boots: %w", err)
	}
	body, err := c.bc.ConvertItemInstance(pk.Body)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobArmourEquipment: failed to convert body: %w", err)
	}
	dst := *pk
	dst.Helmet = helmet
	dst.Chestplate = chestplate
	dst.Leggings = leggings
	dst.Boots = boots
	dst.Body = body
	return &dst, nil
}

// ConvertMobEquipment converts item block runtime IDs inside a MobEquipment packet.
func (c *MinecraftConverter) ConvertMobEquipment(pk *packet.MobEquipment) (*packet.MobEquipment, error) {
	newItem, err := c.bc.ConvertItemInstance(pk.NewItem)
	if err != nil {
		return nil, fmt.Errorf("ConvertMobEquipment: failed to convert new item: %w", err)
	}
	dst := *pk
	dst.NewItem = newItem
	return &dst, nil
}

// ConvertBiomeDefinitionList converts block runtime IDs inside a BiomeDefinitionList packet.
func (c *MinecraftConverter) ConvertBiomeDefinitionList(pk *packet.BiomeDefinitionList) (*packet.BiomeDefinitionList, error) {
	biomeDefinitions, err := utils.ConvertSliceWithError(pk.BiomeDefinitions, c.bc.ConvertBiomeDefinition)
	if err != nil {
		return nil, fmt.Errorf("ConvertBiomeDefinitionList: failed to convert biome definitions: %w", err)
	}
	dst := *pk
	dst.BiomeDefinitions = biomeDefinitions
	return &dst, nil
}

// ConvertEvent converts block runtime IDs inside an Event packet.
func (c *MinecraftConverter) ConvertEvent(pk *packet.Event) (*packet.Event, error) {
	eventData, err := c.ConvertEventData(pk.Event)
	if err != nil {
		return nil, fmt.Errorf("ConvertEvent: failed to convert event data: %w", err)
	}
	dst := *pk
	dst.Event = eventData
	return &dst, nil
}

// ConvertEventData converts block runtime IDs inside event data.
func (c *MinecraftConverter) ConvertEventData(event protocol.Event) (protocol.Event, error) {
	switch typedEvent := event.(type) {
	case nil:
		return nil, nil
	case *protocol.WaxedOrUnwaxedCopperEvent:
		return c.bc.ConvertWaxedOrUnwaxedCopperEvent(typedEvent)
	default:
		return event, nil
	}
}

// ConvertStartGame adjusts StartGame block runtime ID compatibility fields.
func (c *MinecraftConverter) ConvertStartGame(pk *packet.StartGame) (*packet.StartGame, error) {
	dst := *pk
	dst.UseBlockNetworkIDHashes = c.dstTable.UseNetworkIDHashes()
	return &dst, nil
}

// HandleSubChunk converts and writes a SubChunk packet.
func (c *MinecraftConverter) HandleSubChunk(pk *packet.SubChunk) error {
	dst, err := c.ConvertSubChunk(pk)
	if err != nil {
		return err
	}
	return c.dstConn.WritePacket(dst)
}

// HandlePacket converts a supported packet and writes it to the destination connection.
func (c *MinecraftConverter) HandlePacket(pk packet.Packet) error {
	dst, err := c.ConvertPacket(pk)
	if err != nil {
		return err
	}
	return c.dstConn.WritePacket(dst)
}
