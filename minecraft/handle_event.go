package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertEvent converts block runtime IDs inside an Event packet.
func (c *MinecraftConverter) ConvertEvent(pk *packet.Event) (*packet.Event, error) {
	eventData, err := c.bc.ConvertEvent(pk.Event)
	if err != nil {
		return nil, fmt.Errorf("ConvertEvent: failed to convert event data: %w", err)
	}
	dst := *pk
	dst.Event = eventData
	return &dst, nil
}
