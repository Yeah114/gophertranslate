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
		return c.cc.ConvertSubChunk(pkt)
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

// HandlePacket converts a supported packet and writes it to the destination connection.
func (c *MinecraftConverter) HandlePacket(pk packet.Packet) error {
	dst, err := c.ConvertPacket(pk)
	if err != nil {
		return err
	}
	return c.dstConn.WritePacket(dst)
}
