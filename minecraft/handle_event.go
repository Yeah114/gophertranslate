package minecraft

import (
	"fmt"

	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleEvent converts and writes block runtime IDs inside an Event packet.
func (c *MinecraftConverter) HandleEvent(pk *packet.Event) error {
	eventData, err := c.bc.ConvertServerEvent(pk.Event)
	if err != nil {
		return fmt.Errorf("HandleEvent: failed to convert event data: %w", err)
	}
	dst := *pk
	dst.Event = eventData
	return c.clientConnEcho.WritePacket(&dst)
}
