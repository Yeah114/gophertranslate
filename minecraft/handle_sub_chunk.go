package v1v26v10

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// ConvertSubChunk converts a SubChunk packet.
func (c *MinecraftConverter) ConvertSubChunk(pk *packet.SubChunk) (*packet.SubChunk, error) {
	return c.cc.ConvertSubChunk(pk)
}

// HandleSubChunk converts and writes a SubChunk packet.
func (c *MinecraftConverter) HandleSubChunk(pk *packet.SubChunk) error {
	dst, err := c.ConvertSubChunk(pk)
	if err != nil {
		return err
	}
	return c.dstConn.WritePacket(dst)
}
