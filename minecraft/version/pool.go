package version

import (
	"sort"

	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v16v100"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v16v210"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v17v0"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v17v10"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v17v30"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v17v40"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v18v10"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v18v30"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v0"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v20"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v50"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v60"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v70"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v19v80"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v0"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v10"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v30"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v40"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v50"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v60"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v70"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v20v80"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v0"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v20"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v30"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v40"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v50"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v60"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v70"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v80"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v90"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v100"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v110"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v21v130"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v26v0"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v26v20"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

var Pool = map[int32]func(define.MinecraftConverter) define.VersionConverter{
	protocol.Protocol1v16v100: v1v16v100.NewVersionConverter,
	protocol.Protocol1v16v210: v1v16v210.NewVersionConverter,
	protocol.Protocol1v17v0: v1v17v0.NewVersionConverter,
	protocol.Protocol1v17v10: v1v17v10.NewVersionConverter,
	protocol.Protocol1v17v30: v1v17v30.NewVersionConverter,
	protocol.Protocol1v17v40: v1v17v40.NewVersionConverter,
	protocol.Protocol1v18v10: v1v18v10.NewVersionConverter,
	protocol.Protocol1v18v30: v1v18v30.NewVersionConverter,
	protocol.Protocol1v19v0: v1v19v0.NewVersionConverter,
	protocol.Protocol1v19v20: v1v19v20.NewVersionConverter,
	protocol.Protocol1v19v50: v1v19v50.NewVersionConverter,
	protocol.Protocol1v19v60: v1v19v60.NewVersionConverter,
	protocol.Protocol1v19v70: v1v19v70.NewVersionConverter,
	protocol.Protocol1v19v80: v1v19v80.NewVersionConverter,
	protocol.Protocol1v20v0: v1v20v0.NewVersionConverter,
	protocol.Protocol1v20v10: v1v20v10.NewVersionConverter,
	protocol.Protocol1v20v30: v1v20v30.NewVersionConverter,
	protocol.Protocol1v20v40: v1v20v40.NewVersionConverter,
	protocol.Protocol1v20v50: v1v20v50.NewVersionConverter,
	protocol.Protocol1v20v60: v1v20v60.NewVersionConverter,
	protocol.Protocol1v20v70: v1v20v70.NewVersionConverter,
	protocol.Protocol1v20v80: v1v20v80.NewVersionConverter,
	protocol.Protocol1v21v0: v1v21v0.NewVersionConverter,
	protocol.Protocol1v21v20: v1v21v20.NewVersionConverter,
	protocol.Protocol1v21v30: v1v21v30.NewVersionConverter,
	protocol.Protocol1v21v40: v1v21v40.NewVersionConverter,
	protocol.Protocol1v21v50: v1v21v50.NewVersionConverter,
	protocol.Protocol1v21v60: v1v21v60.NewVersionConverter,
	protocol.Protocol1v21v70: v1v21v70.NewVersionConverter,
	protocol.Protocol1v21v80: v1v21v80.NewVersionConverter,
	protocol.Protocol1v21v90: v1v21v90.NewVersionConverter,
	protocol.Protocol1v21v100: v1v21v100.NewVersionConverter,
	protocol.Protocol1v21v110: v1v21v110.NewVersionConverter,
	protocol.Protocol1v21v130: v1v21v130.NewVersionConverter,
	protocol.Protocol1v26v0: v1v26v0.NewVersionConverter,
	protocol.Protocol1v26v10: v1v26v10.NewVersionConverter,
	protocol.Protocol1v26v20v26: v1v26v20.NewVersionConverter,
	protocol.Protocol1v26v20: v1v26v20.NewVersionConverter,
}

func GetVersionConverters(sourceProto, targetProto int32) []func(define.MinecraftConverter) define.VersionConverter {
	type pair struct {
		proto int32
		ctor  func(define.MinecraftConverter) define.VersionConverter
	}
	var list []pair
	minProto, maxProto := min(sourceProto, targetProto), max(sourceProto, targetProto)

	for proto, ctor := range Pool {
		if proto > minProto && proto <= maxProto {
			list = append(list, pair{proto: proto, ctor: ctor})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].proto > list[j].proto
	})

	var res []func(define.MinecraftConverter) define.VersionConverter
	for _, item := range list {
		res = append(res, item.ctor)
	}
	return res
}
