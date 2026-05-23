package minecraft

import "github.com/Yeah114/gophertunnel/minecraft/protocol/packet"

// HandleStartGame adjusts StartGame compatibility fields and writes the packet.
func (c *MinecraftConverter) HandleStartGame(pk *packet.StartGame) error {
	dst := *pk
	dst.UseBlockNetworkIDHashes = c.serverTable.UseNetworkIDHashes()
	return c.clientConnEcho.WritePacket(&dst)
}
