package block

import (
	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v21v130"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v0"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v20"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v20v26"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

// Pool holds functions that create BlockRuntimeIDTables for different Minecraft versions.
// The key is the protocol version.
var Pool = map[int32]func(bool) *block.BlockRuntimeIDTable{
	protocol.Protocol1v26v20:    v1v26v20.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v20v26: v1v26v20v26.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v10:    v1v26v10.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v0:     v1v26v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v130:   v1v21v130.NewBlockRuntimeIDTable,
}
