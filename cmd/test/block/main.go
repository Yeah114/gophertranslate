package main

import (
	"fmt"

	bwo_chunk "github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v20"
	"github.com/Yeah114/gopherconvert/minecraft/world/chunk"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

func main() {
	block.Init(nil)

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

	clientAir := srcTable.AirRuntimeID()
	clientSubChunk := bwo_chunk.NewSubChunk(clientAir)
	clientDirt, found := srcTable.StateToRuntimeID("minecraft:cinnabar", nil)
	if !found {
		panic("failed to find source dirt runtime ID")
	}
	clientSubChunk.SetBlock(1, 2, 3, 0, clientDirt)

	serverSubChunk, ok := cc.ConvertSubChunk(clientSubChunk)
	if !ok {
		panic("failed to convert subchunk")
	}
	fmt.Printf("Converted block at (1, 2, 3) from runtime ID %d to %d\n", clientSubChunk.Block(1, 2, 3, 0), serverSubChunk.Block(1, 2, 3, 0))

	name, properties, found := srcTable.RuntimeIDToState(1366)
	fmt.Println(name, properties, found)
}
