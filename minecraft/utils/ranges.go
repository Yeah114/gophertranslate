package utils

import (
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gophertunnel/minecraft"
)

// RangesFromGameData extracts dimension ranges from the provided GameData
// and returns a mapping of dimension types to their respective ranges.
func RangesFromGameData(data minecraft.GameData) map[int32]define.Range {
	ranges := map[int32]define.Range{
		0: define.Dimension(0).Range(),
		1: define.Dimension(1).Range(),
		2: define.Dimension(2).Range(),
	}

	for _, dimension := range data.Dimensions {
		r := dimension.Range
		ranges[dimension.DimensionType] = define.Range{int(r[0]), int(r[1])}
	}

	return ranges
}
