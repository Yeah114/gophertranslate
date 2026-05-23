package minecraft

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// HandleLevelChunk converts and writes a LevelChunk packet.
func (c *MinecraftConverter) HandleLevelChunk(pk *packet.LevelChunk) error {
	dst, err := c.cc.ConvertLevelChunk(pk)
	if err != nil {
		return err
	}
	return c.clientConnEcho.WritePacket(dst)
}
