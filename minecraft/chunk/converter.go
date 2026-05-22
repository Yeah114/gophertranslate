package chunk

import (
	_ "runtime"
	"unsafe"

	"github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/gopherconvert/minecraft/block"
)

//go:noescape
//go:linkname memmove runtime.memmove
//goland:noinspection GoUnusedParameter
func memmove(to, from unsafe.Pointer, n uintptr)

// ChunkConverter is a struct that can convert sub chunks between two protocol versions using a BlockConverter.
type ChunkConverter struct {
	bc *block.BlockConverter
}

// NewChunkConverter creates a new ChunkConverter that can convert sub chunks between two protocol versions using a BlockConverter.
func NewChunkConverter(bc *block.BlockConverter) *ChunkConverter {
	return &ChunkConverter{bc: bc}
}

// BlockConverter returns the BlockConverter used by this ChunkConverter.
func (c *ChunkConverter) BlockConverter() *block.BlockConverter {
	return c.bc
}

// ConvertSubChunk converts a sub chunk from one protocol version to another.
// It returns the converted sub chunk and a boolean indicating whether the conversion was successful.
func (c *ChunkConverter) ConvertSubChunk(srcSubChunk *chunk.SubChunk) (dstSubChunk *chunk.SubChunk, ok bool) {
	dstSubChunk = &chunk.SubChunk{}
	memmove(unsafe.Pointer(dstSubChunk), unsafe.Pointer(srcSubChunk), unsafe.Sizeof(*srcSubChunk))

	dstAir := dstSubChunk.Air()
	*dstAir, ok = c.bc.ConvertBlockRuntimeID(*dstAir)
	if !ok {
		return nil, false
	}

	ok = true
	for _, storage := range srcSubChunk.Layers() {
		storage.Palette().Replace(func(srcBlockRuntimeID uint32) uint32 {
			dstBlockRuntimeID, found := c.bc.ConvertBlockRuntimeID(srcBlockRuntimeID)
			if !found {
				ok = false
				return srcBlockRuntimeID
			}
			return dstBlockRuntimeID
		})
		if !ok {
			return nil, false
		}
	}
	return dstSubChunk, true
}

// ConvertChunk converts a chunk from one protocol version to another.
// It returns the converted chunk and a boolean indicating whether the conversion was successful.
func (c *ChunkConverter) ConvertChunk(srcChunk *chunk.Chunk) (dstChunk *chunk.Chunk, ok bool) {
	dstChunk = &chunk.Chunk{}
	memmove(unsafe.Pointer(dstChunk), unsafe.Pointer(srcChunk), unsafe.Sizeof(*srcChunk))

	dstAir := dstChunk.Air()
	*dstAir, ok = c.bc.ConvertBlockRuntimeID(*dstAir)
	if !ok {
		return nil, false
	}

	ok = true
	for _, sub := range srcChunk.Sub() {
		dstSub, subOk := c.ConvertSubChunk(sub)
		if !subOk {
			ok = false
			break
		}
		sub = dstSub
	}
	if !ok {
		return nil, false
	}
	return dstChunk, true
}
