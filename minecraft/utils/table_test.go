package utils

import (
	"testing"

	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func TestFullBlockRuntimeIDTableFromGameDataAfterConversion(t *testing.T) {
	converter := newTestBlockConverter(t, protocol.Protocol1v26v0, protocol.Protocol1v21v90)
	runtimeID, found := converter.ClientTable().StateToRuntimeID("minecraft:lightning_rod", map[string]any{
		"facing_direction": int32(0),
		"powered_bit":      byte(0),
	})
	if !found {
		t.Fatal("failed to find lightning rod runtime ID")
	}
	if _, ok := converter.ConvertClientBlockRuntimeID(runtimeID); !ok {
		t.Fatal("expected runtime ID conversion to succeed")
	}

	for range 2 {
		if _, err := FullBlockRuntimeIDTableFromGameData(minecraft.GameData{UseBlockNetworkIDHashes: true}); err != nil {
			t.Fatalf("failed to create full block runtime ID table: %v", err)
		}
	}
}

func newTestBlockConverter(t *testing.T, clientProtocol, serverProtocol int32) *world_block.BlockConverter {
	t.Helper()
	clientConstructor, ok := world_block.Pool[clientProtocol]
	if !ok {
		t.Fatalf("missing source table constructor for protocol %d", clientProtocol)
	}
	serverConstructor, ok := world_block.Pool[serverProtocol]
	if !ok {
		t.Fatalf("missing destination table constructor for protocol %d", serverProtocol)
	}
	return world_block.NewBlockConverter(clientProtocol, clientConstructor(false), serverProtocol, serverConstructor(false))
}
