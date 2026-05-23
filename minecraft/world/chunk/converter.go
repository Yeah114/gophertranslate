package chunk

import (
	_ "runtime"
	"unsafe"

	"github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/gopherconvert/minecraft/world/block"
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
func (c *ChunkConverter) ConvertSubChunk(clientSubChunk *chunk.SubChunk) (serverSubChunk *chunk.SubChunk, ok bool) {
	serverSubChunk = &chunk.SubChunk{}
	memmove(unsafe.Pointer(serverSubChunk), unsafe.Pointer(clientSubChunk), unsafe.Sizeof(*clientSubChunk))

	serverAir := serverSubChunk.Air()
	*serverAir, ok = c.bc.ConvertServerBlockRuntimeID(*serverAir)
	if !ok {
		return nil, false
	}

	ok = true
	for _, storage := range clientSubChunk.Layers() {
		storage.Palette().Replace(func(clientBlockRuntimeID uint32) uint32 {
			serverBlockRuntimeID, found := c.bc.ConvertServerBlockRuntimeID(clientBlockRuntimeID)
			if !found {
				ok = false
				return clientBlockRuntimeID
			}
			return serverBlockRuntimeID
		})
		if !ok {
			return nil, false
		}
	}
	return serverSubChunk, true
}

// ConvertChunk converts a chunk from one protocol version to another.
// It returns the converted chunk and a boolean indicating whether the conversion was successful.
func (c *ChunkConverter) ConvertChunk(clientChunk *chunk.Chunk) (serverChunk *chunk.Chunk, ok bool) {
	serverChunk = &chunk.Chunk{}
	memmove(unsafe.Pointer(serverChunk), unsafe.Pointer(clientChunk), unsafe.Sizeof(*clientChunk))

	serverAir := serverChunk.Air()
	*serverAir, ok = c.bc.ConvertServerBlockRuntimeID(*serverAir)
	if !ok {
		return nil, false
	}

	ok = true
	for _, sub := range clientChunk.Sub() {
		serverSub, subOk := c.ConvertSubChunk(sub)
		if !subOk {
			ok = false
			break
		}
		sub = serverSub
	}
	if !ok {
		return nil, false
	}
	return serverChunk, true
}
