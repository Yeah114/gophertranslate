package v1v26v20

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
)

type VersionConverter struct {
	c define.MinecraftConverter
}

func NewVersionConverter(c define.MinecraftConverter) define.VersionConverter {
	return &VersionConverter{c: c}
}

func (c *VersionConverter) StartGame(data *minecraft.GameData) (err error) {
	data.CustomBlocks = append(data.CustomBlocks, utils.BlockStatesToBlockEntries(BlockStates)...)
	table := c.c.BlockConverter().BlockConverter().ClientTable()
	for _, state := range BlockStates {
		err := table.RegisterCustomBlock(state)
		if err != nil {
			return fmt.Errorf("v1.26v20.VersionConverter.StartGame: failed to register custom block: %w", err)
		}
	}
	table.FinaliseRegister()

	return nil
}

func (c *VersionConverter) HandlePacket(pk packet.Packet, sender define.Conn) (err error) {
	return nil
}
