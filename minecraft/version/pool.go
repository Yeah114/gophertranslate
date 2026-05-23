package version

import (
	"sort"

	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/version/v1v26v20"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

var Pool = map[int32]func(define.MinecraftConverter) define.VersionConverter{
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
