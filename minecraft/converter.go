package minecraft

import (
	"fmt"
	"context"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gopherconvert/minecraft/define"
	protocol_block "github.com/Yeah114/gopherconvert/minecraft/protocol/block"
	"github.com/Yeah114/gopherconvert/minecraft/protocol/chunk"
	protocol_item "github.com/Yeah114/gopherconvert/minecraft/protocol/item"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gopherconvert/minecraft/version"
	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	world_chunk "github.com/Yeah114/gopherconvert/minecraft/world/chunk"
	world_item "github.com/Yeah114/gopherconvert/minecraft/world/item"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// MinecraftConverter converts Minecraft packets between two connections.
type MinecraftConverter struct {
	clientConn     define.Conn
	serverConn     define.Conn
	clientConnEcho define.Conn
	serverConnEcho define.Conn
	clientTable    *bwo_block.BlockRuntimeIDTable
	serverTable    *bwo_block.BlockRuntimeIDTable
	clientItems    *world_item.ItemRuntimeIDTable
	serverItems    *world_item.ItemRuntimeIDTable
	bc             *protocol_block.BlockConverter
	cc             *chunk.ChunkConverter
	ic             *protocol_item.ItemConverter
	vcs            []define.VersionConverter
	dcs            map[string]define.DaemonConverter
}

var _ define.MinecraftConverter = (*MinecraftConverter)(nil)

func NewMinecraftConverter(clientConn define.Conn, serverConn define.Conn) *MinecraftConverter {
	clientConnEcho := define.NewEchoConn(clientConn.(*minecraft.Conn))
	serverConnEcho := define.NewEchoConn(serverConn.(*minecraft.Conn))
	return &MinecraftConverter{
		clientConn:     clientConn,
		serverConn:     serverConn,
		clientConnEcho: clientConnEcho,
		serverConnEcho: serverConnEcho,
		dcs:            make(map[string]define.DaemonConverter),
	}
}

func (c *MinecraftConverter) StartGameContext(ctx context.Context, data *minecraft.GameData) error {
	serverGameData := c.serverConn.GameData()
	clientTable, err := utils.BlockRuntimeIDTableFromGameDataAndVersion(*data, c.clientConn.Proto().Ver())
	if err != nil {
		return fmt.Errorf("NewMinecraftConverter: failed to create block runtime ID table for source connection: %w", err)
	}
	serverTable, err := utils.BlockRuntimeIDTableFromGameDataAndVersion(serverGameData, c.serverConn.Proto().Ver())
	if err != nil {
		return fmt.Errorf("NewMinecraftConverter: failed to create block runtime ID table for destination connection: %w", err)
	}

	clientItems, err := utils.ItemRuntimeIDTableFromGameData(*data)
	if err != nil {
		return fmt.Errorf("NewMinecraftConverter: failed to create item runtime ID table for source connection: %w", err)
	}
	serverItems, err := utils.ItemRuntimeIDTableFromGameData(serverGameData)
	if err != nil {
		return fmt.Errorf("NewMinecraftConverter: failed to create item runtime ID table for destination connection: %w", err)
	}

	blockConverter := world_block.NewBlockConverter(
		c.clientConn.Proto().ID(),
		clientTable,
		c.serverConn.Proto().ID(),
		serverTable,
	)
	bc := protocol_block.NewBlockConverter(blockConverter)
	cc := chunk.NewChunkConverter(
		world_chunk.NewChunkConverter(blockConverter),
		utils.RangesFromGameData(serverGameData),
		serverGameData.Dimension,
	)
	ic := protocol_item.NewItemConverter(
		world_item.NewItemConverter(c.clientConn.Proto().ID(), clientItems, c.serverConn.Proto().ID(), serverItems),
		bc,
	)

	c.clientTable = clientTable
	c.serverTable = serverTable
	c.clientItems = clientItems
	c.serverItems = serverItems
	c.bc = bc
	c.cc = cc
	c.ic = ic
	vcs := version.GetVersionConverters(c.clientConn.Proto().ID(), c.serverConn.Proto().ID())
	c.vcs = make([]define.VersionConverter, len(vcs))
	for i, f := range vcs {
		v := f(c)
		err := v.StartGame(data)
		if err != nil {
			return fmt.Errorf("MinecraftConverter: failed to convert game data: %w", err)
		}
		c.vcs[i] = v
	}
	clientTable.FinaliseRegister()
	serverTable.FinaliseRegister()

	err = c.clientConn.StartGameContext(ctx, *data)
	if err != nil {
		return fmt.Errorf("MinecraftConverter: failed to start game: %w", err)
	}
	return nil
}

// ClientConn returns the source connection.
func (c *MinecraftConverter) ClientConn() define.Conn {
	return c.clientConn
}

// ClientConnEcho returns the echo connection for packets from the source side.
func (c *MinecraftConverter) ClientConnEcho() define.Conn {
	return c.clientConnEcho
}

// ServerConn returns the destination connection.
func (c *MinecraftConverter) ServerConn() define.Conn {
	return c.serverConn
}

// ServerConnEcho returns the echo connection for packets from the destination side.
func (c *MinecraftConverter) ServerConnEcho() define.Conn {
	return c.serverConnEcho
}

// BlockConverter returns the protocol block converter used by the converter.
func (c *MinecraftConverter) BlockConverter() define.BlockConverter {
	return c.bc
}

// ChunkConverter returns the protocol chunk converter used by the converter.
func (c *MinecraftConverter) ChunkConverter() define.ChunkConverter {
	return c.cc
}

// ItemConverter returns the protocol item converter used by the converter.
func (c *MinecraftConverter) ItemConverter() define.ItemConverter {
	return c.ic
}

func (c *MinecraftConverter) RunDaemonConverter(dc define.DaemonConverter) bool {
	if dc == nil {
		return false
	}
	if existing, ok := c.dcs[dc.Name()]; ok {
		existing.Start()
		return true
	}
	c.dcs[dc.Name()] = dc
	dc.Start()
	return false
}

func (c *MinecraftConverter) GetDaemonConverter(name string) (define.DaemonConverter, bool) {
	dc, ok := c.dcs[name]
	return dc, ok
}

func (c *MinecraftConverter) StopDaemonConverter(name string) bool {
	dc, ok := c.dcs[name]
	if !ok {
		return false
	}
	dc.Stop()
	delete(c.dcs, name)
	return true
}

func (c *MinecraftConverter) drainEchoPackets(echo define.Conn) ([]packet.Packet, error) {
	echoConn, ok := echo.(*define.EchoConn)
	if !ok {
		return nil, fmt.Errorf("drainEchoPackets: unexpected echo connection type %T", echo)
	}

	var packets []packet.Packet
	for {
		select {
		case pk := <-echoConn.Packets():
			packets = append(packets, pk)
		default:
			return packets, nil
		}
	}
}

func (c *MinecraftConverter) HandlePacket(pk packet.Packet, sender define.Conn) (err error) {
	if pk == nil {
		return fmt.Errorf("HandlePacket: packet is nil")
	}
	if sender == nil {
		return fmt.Errorf("HandlePacket: sender is nil")
	}

	clientSend := sender == c.clientConn
	serverSend := sender == c.serverConn
	if !clientSend && !serverSend {
		return fmt.Errorf("HandlePacket: unknown sender")
	}

	if clientSend {
		switch pkt := pk.(type) {
		case *packet.InventoryTransaction:
			err = c.HandleInventoryTransaction(pkt)
		case *packet.MobEquipment:
			err = c.HandleMobEquipment(pkt)
		default:
			err = c.serverConnEcho.WritePacket(pk)
		}
	} else if serverSend {
		switch pkt := pk.(type) {
		case *packet.AddItemActor:
			err = c.HandleAddItemActor(pkt)
		case *packet.AddPlayer:
			err = c.HandleAddPlayer(pkt)
		case *packet.BiomeDefinitionList:
			err = c.HandleBiomeDefinitionList(pkt)
		case *packet.ClientCacheMissResponse:
			err = c.HandleClientCacheMissResponse(pkt)
		case *packet.CraftingData:
			return c.HandleCraftingData(pkt)
		case *packet.CreativeContent:
			err = c.HandleCreativeContent(pkt)
		case *packet.Event:
			err = c.HandleEvent(pkt)
		case *packet.InventoryContent:
			err = c.HandleInventoryContent(pkt)
		case *packet.InventorySlot:
			err = c.HandleInventorySlot(pkt)
		case *packet.ItemRegistry:
			err = c.HandleItemRegistry(pkt)
		case *packet.LevelChunk:
			err = c.HandleLevelChunk(pkt)
		case *packet.MobArmourEquipment:
			err = c.HandleMobArmourEquipment(pkt)
		case *packet.StartGame:
			err = c.HandleStartGame(pkt)
		case *packet.SubChunk:
			err = c.HandleSubChunk(pkt)
		case *packet.UpdateBlock:
			err = c.HandleUpdateBlock(pkt)
		case *packet.UpdateBlockSynced:
			err = c.HandleUpdateBlockSynced(pkt)
		case *packet.UpdateSubChunkBlocks:
			err = c.HandleUpdateSubChunkBlocks(pkt)
		default:
			err = c.clientConnEcho.WritePacket(pk)
		}
	}
	if err != nil {
		return fmt.Errorf("HandlePacket: failed to convert packet: %w", err)
	}

	for _, vc := range c.vcs {
		serverPks, err := c.drainEchoPackets(c.serverConnEcho)
		if err != nil {
			return fmt.Errorf("HandlePacket: failed to drain server echo: %w", err)
		}
		for _, pk := range serverPks {
			if err := vc.HandlePacket(pk, c.clientConn); err != nil {
				return fmt.Errorf("HandlePacket: failed to convert server packet: %w", err)
			}
		}

		clientPks, err := c.drainEchoPackets(c.clientConnEcho)
		if err != nil {
			return fmt.Errorf("HandlePacket: failed to drain client echo: %w", err)
		}
		for _, pk := range clientPks {
			if err := vc.HandlePacket(pk, c.serverConn); err != nil {
				return fmt.Errorf("HandlePacket: failed to convert client packet: %w", err)
			}
		}
	}

	clientPks, err := c.drainEchoPackets(c.clientConnEcho)
	if err != nil {
		return fmt.Errorf("HandlePacket: failed to drain client echo: %w", err)
	}
	for _, pk := range clientPks {
		if err := c.clientConn.WritePacket(pk); err != nil {
			return fmt.Errorf("HandlePacket: failed to write client packet: %w", err)
		}
	}

	serverPks, err := c.drainEchoPackets(c.serverConnEcho)
	if err != nil {
		return fmt.Errorf("HandlePacket: failed to drain server echo: %w", err)
	}
	for _, pk := range serverPks {
		if err := c.serverConn.WritePacket(pk); err != nil {
			return fmt.Errorf("HandlePacket: failed to write server packet: %w", err)
		}
	}

	return nil
}
