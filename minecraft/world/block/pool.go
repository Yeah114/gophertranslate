package block

import (
	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v16v100"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v16v210"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v17v0"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v17v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v17v30"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v17v40"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v18v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v18v30"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v0"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v20"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v50"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v60"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v70"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v19v80"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v0"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v30"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v40"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v50"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v60"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v70"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v20v80"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v0"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v100"
//	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v110"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v120"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v124"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v130"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v20"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v30"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v40"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v50"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v60"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v70"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v80"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v21v90"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v0"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v20v26"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v20"
)

// Pool holds functions that create BlockRuntimeIDTables for different Minecraft versions.
// The key is the protocol version.
var Pool = map[int32]func(bool) *block.BlockRuntimeIDTable{
	protocol.Protocol1v16v100:  v1v16v100.NewBlockRuntimeIDTable,
	protocol.Protocol1v16v210:  v1v16v210.NewBlockRuntimeIDTable,
	protocol.Protocol1v17v0:    v1v17v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v17v10:   v1v17v10.NewBlockRuntimeIDTable,
	protocol.Protocol1v17v30:   v1v17v30.NewBlockRuntimeIDTable,
	protocol.Protocol1v17v40:   v1v17v40.NewBlockRuntimeIDTable,
	protocol.Protocol1v18v10:   v1v18v10.NewBlockRuntimeIDTable,
	protocol.Protocol1v18v30:   v1v18v30.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v0:    v1v19v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v20:   v1v19v20.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v50:   v1v19v50.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v60:   v1v19v60.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v70:   v1v19v70.NewBlockRuntimeIDTable,
	protocol.Protocol1v19v80:   v1v19v80.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v0:    v1v20v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v10:   v1v20v10.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v30:   v1v20v30.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v40:   v1v20v40.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v50:   v1v20v50.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v60:   v1v20v60.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v70:   v1v20v70.NewBlockRuntimeIDTable,
	protocol.Protocol1v20v80:   v1v20v80.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v0:    v1v21v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v20:   v1v21v20.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v30:   v1v21v30.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v40:   v1v21v40.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v50:   v1v21v50.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v60:   v1v21v60.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v70:   v1v21v70.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v80:   v1v21v80.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v90:   v1v21v90.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v100:  v1v21v100.NewBlockRuntimeIDTable,
//	protocol.Protocol1v21v110:  v1v21v110.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v110:  v1v21v130.NewBlockRuntimeIDTable, // From PowerNukkitX
	protocol.Protocol1v21v120:  v1v21v120.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v124:  v1v21v124.NewBlockRuntimeIDTable,
	protocol.Protocol1v21v130:  v1v21v130.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v0:    v1v26v0.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v10:   v1v26v10.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v20v26: v1v26v20v26.NewBlockRuntimeIDTable,
	protocol.Protocol1v26v20:   v1v26v20.NewBlockRuntimeIDTable,
}
