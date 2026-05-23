package minecraft

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// HandleClientCacheMissResponse converts and writes a ClientCacheMissResponse packet.
func (c *MinecraftConverter) HandleClientCacheMissResponse(pk *packet.ClientCacheMissResponse) error {
	dst, err := c.cc.ConvertClientCacheMissResponse(pk)
	if err != nil {
		return err
	}
	return c.clientConnEcho.WritePacket(dst)
}
