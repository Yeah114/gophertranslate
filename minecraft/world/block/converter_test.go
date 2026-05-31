package block

import (
	"testing"

	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func TestBlockConverterConvertClientBlockStateSameVersion(t *testing.T) {
	converter := newTestBlockConverter(t, protocol.Protocol1v26v20, protocol.Protocol1v26v20)
	properties := map[string]any{"foo": "bar"}
	name, convertedProperties, ok := converter.ConvertClientBlockState("minecraft:test", properties)
	if !ok {
		t.Fatal("expected conversion to succeed")
	}
	if name != "minecraft:test" {
		t.Fatalf("expected name to stay unchanged, got %q", name)
	}
	if convertedProperties["foo"] != "bar" {
		t.Fatalf("expected properties to stay unchanged, got %#v", convertedProperties)
	}
}

func TestBlockConverterConvertClientBlockRuntimeID(t *testing.T) {
	converter := newTestBlockConverter(t, protocol.Protocol1v26v10, protocol.Protocol1v26v20)

	clientRuntimeID, found := converter.clientTable.StateToRuntimeID("minecraft:air", nil)
	if !found {
		t.Fatal("failed to find source air runtime ID")
	}
	serverRuntimeID, ok := converter.ConvertClientBlockRuntimeID(clientRuntimeID)
	if !ok {
		t.Fatal("expected runtime ID conversion to succeed")
	}
	if serverRuntimeID != converter.serverTable.AirRuntimeID() {
		t.Fatalf("expected destination air runtime ID %d, got %d", converter.serverTable.AirRuntimeID(), serverRuntimeID)
	}
	if serverRuntimeID == clientRuntimeID {
		t.Fatalf("expected air runtime ID to change between versions, both were %d", serverRuntimeID)
	}
}

func TestBlockConverterConvertClientBlockRuntimeIDMissing(t *testing.T) {
	converter := newTestBlockConverter(t, protocol.Protocol1v26v10, protocol.Protocol1v26v20)
	if runtimeID, ok := converter.ConvertClientBlockRuntimeID(^uint32(0)); ok {
		t.Fatalf("expected conversion to fail for unknown runtime ID, got %d", runtimeID)
	}
}

func TestBlockConverterDoesNotMutateRuntimeTableProperties(t *testing.T) {
	converter := newTestBlockConverter(t, protocol.Protocol1v26v0, protocol.Protocol1v21v90)
	runtimeID, found := converter.clientTable.StateToRuntimeID("minecraft:lightning_rod", map[string]any{
		"facing_direction": int32(0),
		"powered_bit":      byte(0),
	})
	if !found {
		t.Fatal("failed to find lightning rod runtime ID")
	}

	for range 2 {
		if _, ok := converter.ConvertClientBlockRuntimeID(runtimeID); !ok {
			t.Fatal("expected runtime ID conversion to succeed")
		}
		if _, found := converter.clientTable.StateToRuntimeID("minecraft:lightning_rod", map[string]any{
			"facing_direction": int32(0),
			"powered_bit":      byte(0),
		}); !found {
			t.Fatal("runtime table properties were mutated by conversion")
		}
	}
}

func newTestBlockConverter(t *testing.T, clientProtocol, serverProtocol int32) *BlockConverter {
	t.Helper()
	clientConstructor, ok := Pool[clientProtocol]
	if !ok {
		t.Fatalf("missing source table constructor for protocol %d", clientProtocol)
	}
	serverConstructor, ok := Pool[serverProtocol]
	if !ok {
		t.Fatalf("missing destination table constructor for protocol %d", serverProtocol)
	}
	return NewBlockConverter(clientProtocol, clientConstructor(false), serverProtocol, serverConstructor(false))
}
