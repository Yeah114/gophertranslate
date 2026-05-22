package block

import (
	"testing"

	runtimeblock "github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v21v130"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v0"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v20"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v20v26"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func TestBlockRuntimeIDTable(t *testing.T) {
	tests := []struct {
		name         string
		protocol     int32
		airRuntimeID uint32
		constructor  func(bool) *runtimeblock.BlockRuntimeIDTable
	}{
		{name: "1.21.130", protocol: protocol.Protocol1v21v130, airRuntimeID: 12530, constructor: v1v21v130.NewBlockRuntimeIDTable},
		{name: "1.26.0", protocol: protocol.Protocol1v26v0, airRuntimeID: 12530, constructor: v1v26v0.NewBlockRuntimeIDTable},
		{name: "1.26.10", protocol: protocol.Protocol1v26v10, airRuntimeID: 12531, constructor: v1v26v10.NewBlockRuntimeIDTable},
		{name: "1.26.20.26", protocol: protocol.Protocol1v26v20v26, airRuntimeID: 13080, constructor: v1v26v20v26.NewBlockRuntimeIDTable},
		{name: "1.26.20", protocol: protocol.Protocol1v26v20, airRuntimeID: 13080, constructor: v1v26v20.NewBlockRuntimeIDTable},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			table := test.constructor(false)
			airRuntimeID, found := table.StateToRuntimeID("minecraft:air", nil)
			if !found {
				t.Fatalf("protocol %d: failed to find runtime ID for air block", test.protocol)
			}
			if airRuntimeID != table.AirRuntimeID() {
				t.Fatalf("protocol %d: expected air block runtime ID to be %d, got %d", test.protocol, table.AirRuntimeID(), airRuntimeID)
			}
			if airRuntimeID != test.airRuntimeID {
				t.Fatalf("protocol %d: expected vanilla air runtime ID to be %d, got %d", test.protocol, test.airRuntimeID, airRuntimeID)
			}
		})
	}
}
