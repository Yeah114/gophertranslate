package chunk

import (
	"bytes"
	"fmt"

	bwo_chunk "github.com/Yeah114/bedrock-world-operator/chunk"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gopherconvert/minecraft/world/chunk"
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
func (c *ChunkConverter) ConvertSubChunkEntryRawPayload(clientSubChunkEntryRawPayload []byte, r define.Range) (serverSubChunkEntryRawPayload []byte, err error) {
	bc := c.cc.BlockConverter()
	clientBuf := bytes.NewBuffer(clientSubChunkEntryRawPayload)
	clientSubChunk, index, err := bwo_chunk.DecodeSubChunk(clientBuf, r, bwo_chunk.NetworkEncoding, bc.ServerTable())
	if err != nil {
		return nil, fmt.Errorf("ConvertSubChunkEntryRawPayload: failed to decode source sub chunk: %w", err)
	}

	serverSubChunk, ok := c.cc.ConvertSubChunk(clientSubChunk)
	if !ok {
		return nil, fmt.Errorf("ConvertSubChunkEntryRawPayload: failed to convert sub chunk")
	}
	serverSubChunkPayload := bwo_chunk.EncodeSubChunk(serverSubChunk, r, index, bwo_chunk.NetworkEncoding, bc.ClientTable())

	return append(serverSubChunkPayload, clientBuf.Bytes()...), nil
}

// ConvertSubChunkBlobPayload converts a cache blob payload holding a sub chunk from one protocol version to another.
func (c *ChunkConverter) ConvertSubChunkBlobPayload(clientSubChunkBlobPayload []byte, r define.Range) (serverSubChunkBlobPayload []byte, err error) {
	serverPayload, err := c.ConvertSubChunkEntryRawPayload(clientSubChunkBlobPayload, r)
	if err != nil {
		return nil, err
	}
	if len(serverPayload) == 0 {
		return nil, fmt.Errorf("ConvertSubChunkBlobPayload: empty converted payload")
	}
	return serverPayload, nil
}

// ConvertSubChunkEntry converts a sub chunk entry from one protocol version to another.
func (c *ChunkConverter) ConvertSubChunkEntry(clientSubChunkEntry protocol.SubChunkEntry, r define.Range, cacheEnabled bool) (serverSubChunkEntry protocol.SubChunkEntry, err error) {
	var serverRawPayload []byte
	if len(clientSubChunkEntry.RawPayload) != 0 {
		if cacheEnabled {
			serverRawPayload = append([]byte{}, clientSubChunkEntry.RawPayload...)
		} else {
			serverRawPayload, err = c.ConvertSubChunkEntryRawPayload(clientSubChunkEntry.RawPayload, r)
			if err != nil {
				return protocol.SubChunkEntry{}, fmt.Errorf("ConvertSubChunkEntry: failed to convert sub chunk entry raw payload: %w", err)
			}
		}
	}
	serverSubChunkEntry = protocol.SubChunkEntry{
		Offset:              clientSubChunkEntry.Offset,
		Result:              clientSubChunkEntry.Result,
		RawPayload:          serverRawPayload,
		HeightMapType:       clientSubChunkEntry.HeightMapType,
		HeightMapData:       append([]int8{}, clientSubChunkEntry.HeightMapData...),
		RenderHeightMapType: clientSubChunkEntry.RenderHeightMapType,
		RenderHeightMapData: append([]int8{}, clientSubChunkEntry.RenderHeightMapData...),
		BlobHash:            clientSubChunkEntry.BlobHash,
	}
	return serverSubChunkEntry, nil
}

// ConvertSubChunk converts a sub chunk from one protocol version to another.
// It returns the converted sub chunk and a error if the conversion was unsuccessful.
func (c *ChunkConverter) ConvertSubChunk(clientSubChunk *packet.SubChunk) (serverSubChunk *packet.SubChunk, err error) {
	r, found := c.Ranges[clientSubChunk.Dimension]
	if !found {
		return nil, fmt.Errorf("ConvertSubChunk: unsupported dimension: %d", clientSubChunk.Dimension)
	}
	subChunkEntries, err := utils.ConvertSliceWithError(clientSubChunk.SubChunkEntries, func(clientSubChunkEntry protocol.SubChunkEntry) (protocol.SubChunkEntry, error) {
		return c.ConvertSubChunkEntry(clientSubChunkEntry, r, clientSubChunk.CacheEnabled)
	})
	if err != nil {
		return nil, fmt.Errorf("ConvertSubChunk: failed to convert sub chunk entries: %w", err)
	}
	serverSubChunk = &packet.SubChunk{
		CacheEnabled:    clientSubChunk.CacheEnabled,
		Dimension:       clientSubChunk.Dimension,
		Position:        clientSubChunk.Position,
		SubChunkEntries: subChunkEntries,
	}
	return serverSubChunk, nil
}

// ConvertLevelChunkRawPayload converts a LevelChunk raw payload from one protocol version to another.
func (c *ChunkConverter) ConvertLevelChunkRawPayload(clientLevelChunkRawPayload []byte, subChunkCount uint32, r define.Range) (serverLevelChunkRawPayload []byte, err error) {
	bc := c.cc.BlockConverter()
	clientBuf := bytes.NewBuffer(clientLevelChunkRawPayload)
	serverChunk := bwo_chunk.NewChunk(bc.ClientTable().AirRuntimeID(), r)

	for i := uint32(0); i < subChunkCount; i++ {
		clientSubChunk, index, err := bwo_chunk.DecodeSubChunk(clientBuf, r, bwo_chunk.NetworkEncoding, bc.ServerTable())
		if err != nil {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to decode source sub chunk %d: %w", i, err)
		}
		serverSubChunk, ok := c.cc.ConvertSubChunk(clientSubChunk)
		if !ok {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to convert sub chunk %d", i)
		}
		serverChunk.SetSubChunk(serverSubChunk, int16(index))
	}
	if err := bwo_chunk.DecodeBiomes(clientBuf, serverChunk, bwo_chunk.NetworkEncoding, bc.ServerTable()); err != nil {
		return nil, fmt.Errorf("ConvertLevelChunkRawPayload: failed to decode biomes: %w", err)
	}

	serverData := bwo_chunk.Encode(serverChunk, bwo_chunk.NetworkEncoding, bc.ClientTable())
	serverBuf := bytes.NewBuffer(make([]byte, 0, len(clientLevelChunkRawPayload)))
	for i := uint32(0); i < subChunkCount; i++ {
		if int(i) >= len(serverData.SubChunks) {
			return nil, fmt.Errorf("ConvertLevelChunkRawPayload: sub chunk count %d exceeds encoded sub chunk count %d", subChunkCount, len(serverData.SubChunks))
		}
		_, _ = serverBuf.Write(serverData.SubChunks[i])
	}
	_, _ = serverBuf.Write(serverData.Biomes)
	_, _ = serverBuf.Write(clientBuf.Bytes())
	return serverBuf.Bytes(), nil
}

// ConvertLevelChunk converts a LevelChunk packet from one protocol version to another.
func (c *ChunkConverter) ConvertLevelChunk(clientLevelChunk *packet.LevelChunk) (serverLevelChunk *packet.LevelChunk, err error) {
	serverLevelChunk = &packet.LevelChunk{
		Position:        clientLevelChunk.Position,
		Dimension:       clientLevelChunk.Dimension,
		HighestSubChunk: clientLevelChunk.HighestSubChunk,
		SubChunkCount:   clientLevelChunk.SubChunkCount,
		CacheEnabled:    clientLevelChunk.CacheEnabled,
		BlobHashes:      append([]uint64{}, clientLevelChunk.BlobHashes...),
		RawPayload:      append([]byte{}, clientLevelChunk.RawPayload...),
	}
	if clientLevelChunk.CacheEnabled || clientLevelChunk.SubChunkCount >= protocol.SubChunkRequestModeLimited {
		return serverLevelChunk, nil
	}

	r, found := c.Ranges[clientLevelChunk.Dimension]
	if !found {
		return nil, fmt.Errorf("ConvertLevelChunk: unsupported dimension: %d", clientLevelChunk.Dimension)
	}
	serverRawPayload, err := c.ConvertLevelChunkRawPayload(clientLevelChunk.RawPayload, clientLevelChunk.SubChunkCount, r)
	if err != nil {
		return nil, fmt.Errorf("ConvertLevelChunk: failed to convert raw payload: %w", err)
	}
	serverLevelChunk.RawPayload = serverRawPayload
	return serverLevelChunk, nil
}

// ConvertCacheBlob converts a client cache blob from one protocol version to another.
func (c *ChunkConverter) ConvertCacheBlob(clientBlob protocol.CacheBlob) (serverBlob protocol.CacheBlob, err error) {
	r, found := c.Ranges[c.CurrentDimension]
	if !found {
		return serverBlob, fmt.Errorf("ConvertCacheBlob: unsupported dimension: %d", c.CurrentDimension)
	}

	serverPayload, err := c.ConvertSubChunkBlobPayload(clientBlob.Payload, r)
	if err != nil {
		return protocol.CacheBlob{
			Hash:    clientBlob.Hash,
			Payload: append([]byte{}, clientBlob.Payload...),
		}, nil
	}

	return protocol.CacheBlob{
		Hash:    clientBlob.Hash,
		Payload: serverPayload,
	}, nil
}

// ConvertClientCacheMissResponse converts a ClientCacheMissResponse packet from one protocol version to another.
func (c *ChunkConverter) ConvertClientCacheMissResponse(clientClientCacheMissResponse *packet.ClientCacheMissResponse) (serverClientCacheMissResponse *packet.ClientCacheMissResponse, err error) {
	blobs, err := utils.ConvertSliceWithError(clientClientCacheMissResponse.Blobs, c.ConvertCacheBlob)
	if err != nil {
		return nil, fmt.Errorf("ConvertClientCacheMissResponse: failed to convert blobs: %w", err)
	}
	return &packet.ClientCacheMissResponse{Blobs: blobs}, nil
}
