package chunk

import (
	"bytes"
	"fmt"

	bwo_chunk "github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/chunk"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// ChunkConverter is a struct that can convert sub chunks between two protocol versions using a BlockConverter.
type ChunkConverter struct {
	cc *chunk.ChunkConverter
	// Ranges is a map from dimension ID to the corresponding block state ID range.
	Ranges map[int32]define.Range
	// CurrentDimension is the current dimension ID being converted, used for logging purposes.
	CurrentDimension int32
}

// NewChunkConverter creates a new ChunkConverter that can convert sub chunks between two protocol versions using a BlockConverter.
func NewChunkConverter(
	cc *chunk.ChunkConverter,
	ranges map[int32]define.Range,
	currentDimension int32,
) *ChunkConverter {
	if ranges == nil {
		ranges = map[int32]define.Range{
			0: define.Dimension(0).Range(),
			1: define.Dimension(1).Range(),
			2: define.Dimension(2).Range(),
		}
	}
	return &ChunkConverter{
		cc:               cc,
		Ranges:           ranges,
		CurrentDimension: currentDimension,
	}
}

// ConvertSubChunkEntryRawPayload converts a sub chunk entry raw payload from one protocol version to another.
func (c *ChunkConverter) ConvertSubChunkEntryRawPayload(srcSubChunkEntryRawPayload []byte, r define.Range) (dstSubChunkEntryRawPayload []byte, err error) {
	bc := c.cc.BlockConverter()
	srcBuf := bytes.NewBuffer(srcSubChunkEntryRawPayload)
	srcSubChunk, index, err := bwo_chunk.DecodeSubChunk(srcBuf, r, bwo_chunk.NetworkEncoding, bc.SrcTable())
	if err != nil {
		return nil, fmt.Errorf("ConvertSubChunkEntryRawPayload: failed to decode source sub chunk: %w", err)
	}

	dstSubChunk, ok := c.cc.ConvertSubChunk(srcSubChunk)
	if !ok {
		return nil, fmt.Errorf("ConvertSubChunkEntryRawPayload: failed to convert sub chunk")
	}
	dstSubChunkPayload := bwo_chunk.EncodeSubChunk(dstSubChunk, r, index, bwo_chunk.NetworkEncoding, bc.DstTable())

	return append(dstSubChunkPayload, srcBuf.Bytes()...), nil
}

// ConvertSubChunkBlobPayload converts a cache blob payload holding a sub chunk from one protocol version to another.
func (c *ChunkConverter) ConvertSubChunkBlobPayload(srcSubChunkBlobPayload []byte, r define.Range) (dstSubChunkBlobPayload []byte, err error) {
	dstPayload, err := c.ConvertSubChunkEntryRawPayload(srcSubChunkBlobPayload, r)
	if err != nil {
		return nil, err
	}
	if len(dstPayload) == 0 {
		return nil, fmt.Errorf("ConvertSubChunkBlobPayload: empty converted payload")
	}
	return dstPayload, nil
}

// ConvertSubChunkEntry converts a sub chunk entry from one protocol version to another.
func (c *ChunkConverter) ConvertSubChunkEntry(srcSubChunkEntry protocol.SubChunkEntry, r define.Range, cacheEnabled bool) (dstSubChunkEntry protocol.SubChunkEntry, err error) {
	var dstRawPayload []byte
	if len(srcSubChunkEntry.RawPayload) != 0 {
		if cacheEnabled {
			dstRawPayload = append([]byte{}, srcSubChunkEntry.RawPayload...)
		} else {
			dstRawPayload, err = c.ConvertSubChunkEntryRawPayload(srcSubChunkEntry.RawPayload, r)
			if err != nil {
				return protocol.SubChunkEntry{}, fmt.Errorf("ConvertSubChunkEntry: failed to convert sub chunk entry raw payload: %w", err)
			}
		}
	}
	dstSubChunkEntry = protocol.SubChunkEntry{
		Offset:              srcSubChunkEntry.Offset,
		Result:              srcSubChunkEntry.Result,
		RawPayload:          dstRawPayload,
		HeightMapType:       srcSubChunkEntry.HeightMapType,
		HeightMapData:       append([]int8{}, srcSubChunkEntry.HeightMapData...),
		RenderHeightMapType: srcSubChunkEntry.RenderHeightMapType,
		RenderHeightMapData: append([]int8{}, srcSubChunkEntry.RenderHeightMapData...),
		BlobHash:            srcSubChunkEntry.BlobHash,
	}
	return dstSubChunkEntry, nil
}

// ConvertSubChunk converts a sub chunk from one protocol version to another.
// It returns the converted sub chunk and a error if the conversion was unsuccessful.
func (c *ChunkConverter) ConvertSubChunk(srcSubChunk *packet.SubChunk) (dstSubChunk *packet.SubChunk, err error) {
	r, found := c.Ranges[srcSubChunk.Dimension]
	if !found {
		return nil, fmt.Errorf("ConvertSubChunk: unsupported dimension: %d", srcSubChunk.Dimension)
	}
	subChunkEntries, err := utils.ConvertSliceWithError(srcSubChunk.SubChunkEntries, func(srcSubChunkEntry protocol.SubChunkEntry) (protocol.SubChunkEntry, error) {
		return c.ConvertSubChunkEntry(srcSubChunkEntry, r, srcSubChunk.CacheEnabled)
	})
	if err != nil {
		return nil, fmt.Errorf("ConvertSubChunk: failed to convert sub chunk entries: %w", err)
	}
	dstSubChunk = &packet.SubChunk{
		CacheEnabled:    srcSubChunk.CacheEnabled,
		Dimension:       srcSubChunk.Dimension,
		Position:        srcSubChunk.Position,
		SubChunkEntries: subChunkEntries,
	}
	return dstSubChunk, nil
}

// ConvertLevelChunkRawPayload converts a LevelChunk raw payload from one protocol version to another.
func (c *ChunkConverter) ConvertLevelChunkRawPayload(srcLevelChunkRawPayload []byte, subChunkCount uint32, r define.Range) (dstLevelChunkRawPayload []byte, err error) {
	bc := c.cc.BlockConverter()
	srcBuf := bytes.NewBuffer(srcLevelChunkRawPayload)
	dstChunk := bwo_chunk.NewChunk(bc.DstTable().AirRuntimeID(), r)

	for i := uint32(0); i < subChunkCount; i++ {
		srcSubChunk, index, err := bwo_chunk.DecodeSubChunk(srcBuf, r, bwo_chunk.NetworkEncoding, bc.SrcTable())
		if err != nil {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to decode source sub chunk %d: %w", i, err)
		}
		dstSubChunk, ok := c.cc.ConvertSubChunk(srcSubChunk)
		if !ok {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to convert sub chunk %d", i)
		}
		dstChunk.SetSubChunk(dstSubChunk, int16(index))
	}
	if err := bwo_chunk.DecodeBiomes(srcBuf, dstChunk, bwo_chunk.NetworkEncoding, bc.SrcTable()); err != nil {
		return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to decode biomes: %w", err)
	}

	dstData := bwo_chunk.Encode(dstChunk, bwo_chunk.NetworkEncoding, bc.DstTable())
	dstBuf := bytes.NewBuffer(make([]byte, 0, len(srcLevelChunkRawPayload)))
	for i := uint32(0); i < subChunkCount; i++ {
		if int(i) >= len(dstData.SubChunks) {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: sub chunk count %d exceeds encoded sub chunk count %d", subChunkCount, len(dstData.SubChunks))
		}
		_, _ = dstBuf.Write(dstData.SubChunks[i])
	}
	_, _ = dstBuf.Write(dstData.Biomes)
	_, _ = dstBuf.Write(srcBuf.Bytes())
	return dstBuf.Bytes(), nil
}

// ConvertLevelChunk converts a LevelChunk packet from one protocol version to another.
func (c *ChunkConverter) ConvertLevelChunk(srcLevelChunk *packet.LevelChunk) (dstLevelChunk *packet.LevelChunk, err error) {
	dstLevelChunk = &packet.LevelChunk{
		Position:        srcLevelChunk.Position,
		Dimension:       srcLevelChunk.Dimension,
		HighestSubChunk: srcLevelChunk.HighestSubChunk,
		SubChunkCount:   srcLevelChunk.SubChunkCount,
		CacheEnabled:    srcLevelChunk.CacheEnabled,
		BlobHashes:      append([]uint64{}, srcLevelChunk.BlobHashes...),
		RawPayload:      append([]byte{}, srcLevelChunk.RawPayload...),
	}
	if srcLevelChunk.CacheEnabled || srcLevelChunk.SubChunkCount >= protocol.SubChunkRequestModeLimited {
		return dstLevelChunk, nil
	}

	r, found := c.Ranges[srcLevelChunk.Dimension]
	if !found {
		return nil, fmt.Errorf("ConvertLevelChunk: unsupported dimension: %d", srcLevelChunk.Dimension)
	}
	dstRawPayload, err := c.ConvertLevelChunkRawPayload(srcLevelChunk.RawPayload, srcLevelChunk.SubChunkCount, r)
	if err != nil {
		return nil, fmt.Errorf("ConvertLevelChunk: failed to convert raw payload: %w", err)
	}
	dstLevelChunk.RawPayload = dstRawPayload
	return dstLevelChunk, nil
}

// ConvertCacheBlob converts a client cache blob from one protocol version to another.
func (c *ChunkConverter) ConvertCacheBlob(srcBlob protocol.CacheBlob) (dstBlob protocol.CacheBlob, err error) {
	r, found := c.Ranges[c.CurrentDimension]
	if !found {
		return dstBlob, fmt.Errorf("ConvertCacheBlob: unsupported dimension: %d", c.CurrentDimension)
	}

	dstPayload, err := c.ConvertSubChunkBlobPayload(srcBlob.Payload, r)
	if err != nil {
		return protocol.CacheBlob{
			Hash:    srcBlob.Hash,
			Payload: append([]byte{}, srcBlob.Payload...),
		}, nil
	}

	return protocol.CacheBlob{
		Hash:    srcBlob.Hash,
		Payload: dstPayload,
	}, nil
}

// ConvertClientCacheMissResponse converts a ClientCacheMissResponse packet from one protocol version to another.
func (c *ChunkConverter) ConvertClientCacheMissResponse(srcClientCacheMissResponse *packet.ClientCacheMissResponse) (dstClientCacheMissResponse *packet.ClientCacheMissResponse, err error) {
	blobs, err := utils.ConvertSliceWithError(srcClientCacheMissResponse.Blobs, c.ConvertCacheBlob)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientCacheMissResponse: failed to convert blobs: %w", err)
	}
	return &packet.ClientCacheMissResponse{Blobs: blobs}, nil
}
