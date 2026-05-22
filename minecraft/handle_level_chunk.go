package v1v26v10

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// ConvertLevelChunk converts a LevelChunk packet.
func (c *MinecraftConverter) ConvertLevelChunk(pk *packet.LevelChunk) (*packet.LevelChunk, error) {
	return c.cc.ConvertLevelChunk(pk)
}
