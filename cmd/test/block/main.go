package main

import (
	"fmt"

	bwo_chunk "github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/block"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/block/v1v26v20"
	"github.com/Yeah114/gopherconvert/minecraft/chunk"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func main() {
	srcTable := v1v26v20.NewBlockRuntimeIDTable(false)
	dstTable := v1v26v10.NewBlockRuntimeIDTable(false)
	err := dstTable.RegisterCustomBlock(define.BlockState{
		Name: "minecraft:cinnabar",
	})
	if err != nil {
		panic("failed to register custom block")
	}
	dstTable.FinaliseRegister()

	bc := block.NewBlockConverter(protocol.Protocol1v26v20, srcTable, protocol.Protocol1v26v10, dstTable)
	cc := chunk.NewChunkConverter(bc)

	srcAir := srcTable.AirRuntimeID()
	srcSubChunk := bwo_chunk.NewSubChunk(srcAir)
	srcDirt, found := srcTable.StateToRuntimeID("minecraft:cinnabar", nil)
	if !found {
		panic("failed to find source dirt runtime ID")
	}
	srcSubChunk.SetBlock(1, 2, 3, 0, srcDirt)

	dstSubChunk, ok := cc.ConvertSubChunk(srcSubChunk)
	if !ok {
		panic("failed to convert subchunk")
	}
	fmt.Printf("Converted block at (1, 2, 3) from runtime ID %d to %d\n", srcSubChunk.Block(1, 2, 3, 0), dstSubChunk.Block(1, 2, 3, 0))
}
