package v1v19v60

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// VersionConverter applies the protocol-specific adjustments required for
// Minecraft 1.19.60.
type VersionConverter struct {
	c define.MinecraftConverter
}

// NewVersionConverter creates a version converter for Minecraft 1.19.60.
func NewVersionConverter(c define.MinecraftConverter) define.VersionConverter {
	return &VersionConverter{c: c}
}

// StartGame registers the custom block states required by this protocol
// version before the game starts.
func (c *VersionConverter) StartGame(data *minecraft.GameData) (err error) {
	if c.c.ClientConnEcho().Proto().ID() < c.c.ServerConnEcho().Proto().ID() {
		data.CustomBlocks = append(data.CustomBlocks, utils.BlockStatesToBlockEntries(BlockStates)...)
		table := c.c.BlockConverter().BlockConverter().ClientTable()
		for _, state := range BlockStates {
			err := table.RegisterCustomBlock(state)
			if err != nil {
				return fmt.Errorf("v1v19v60.VersionConverter.StartGame: failed to register custom block: %w", err)
			}
		}
		table.FinaliseRegister()
	}
	return nil
}

// HandlePacket processes echo packets from the main converter, applying version-specific
// transformations where needed.
func (c *VersionConverter) HandlePacket(pk packet.Packet, sender define.Conn) (err error) {
	if sender == c.c.ServerConn() {
		return c.c.ClientConnEcho().WritePacket(pk)
	}
	if sender == c.c.ClientConn() {
		return c.c.ServerConnEcho().WritePacket(pk)
	}
	return fmt.Errorf("v1v19v60.VersionConverter.HandlePacket: unknown sender")
}
