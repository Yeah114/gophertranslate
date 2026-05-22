package v1v26v10

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// HandleSubChunk converts and writes a SubChunk packet.
func (c *MinecraftConverter) HandleSubChunk(pk *packet.SubChunk) error {
	dst, err := c.cc.ConvertSubChunk(pk)
	if err != nil {
		return err
	}
	return c.dstConn.WritePacket(dst)
}
