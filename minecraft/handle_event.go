package v1v26v10

import (
	"fmt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ConvertEvent converts block runtime IDs inside an Event packet.
func (c *MinecraftConverter) ConvertEvent(pk *packet.Event) (*packet.Event, error) {
	eventData, err := c.ConvertEventData(pk.Event)
	if err != nil {
		return nil, fmt.Errorf("ConvertEvent: failed to convert event data: %w", err)
	}
	dst := *pk
	dst.Event = eventData
	return &dst, nil
}

// ConvertEventData converts block runtime IDs inside event data.
func (c *MinecraftConverter) ConvertEventData(event protocol.Event) (protocol.Event, error) {
	switch typedEvent := event.(type) {
	case nil:
		return nil, nil
	case *protocol.WaxedOrUnwaxedCopperEvent:
		return c.bc.ConvertWaxedOrUnwaxedCopperEvent(typedEvent)
	default:
		return event, nil
	}
}
