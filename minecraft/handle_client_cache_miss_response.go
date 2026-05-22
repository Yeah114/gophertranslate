package v1v26v10

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// ConvertClientCacheMissResponse converts a ClientCacheMissResponse packet.
func (c *MinecraftConverter) ConvertClientCacheMissResponse(pk *packet.ClientCacheMissResponse) (*packet.ClientCacheMissResponse, error) {
	return c.cc.ConvertClientCacheMissResponse(pk)
}
