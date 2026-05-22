package v1v26v10

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// ConvertStartGame adjusts StartGame block runtime ID compatibility fields.
func (c *MinecraftConverter) ConvertStartGame(pk *packet.StartGame) (*packet.StartGame, error) {
	dst := *pk
	dst.UseBlockNetworkIDHashes = c.dstTable.UseNetworkIDHashes()
	return &dst, nil
}
