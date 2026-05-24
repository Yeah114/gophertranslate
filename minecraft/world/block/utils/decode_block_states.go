package utils

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gophertunnel/minecraft/nbt"
)

// DecodeBlockStates decodes block states from the given byte slice.
func DecodeBlockStates(blockStatesBytes []byte) (blockStates []define.BlockState) {
	blockStatesBytes = decompressBlockStates(blockStatesBytes)
	existHash := make(map[uint32]struct{})

	for offset := 0; offset < len(blockStatesBytes); {
		buf := bytes.NewBuffer(blockStatesBytes[offset:])
		dec := nbt.NewDecoder(buf)

		var s define.BlockState
		if err := dec.Decode(&s); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if _, ok := err.(nbt.BufferOverrunError); ok {
				break
			}
			panic(fmt.Errorf("DecodeBlockStates: Failed to decode block state from NBT: %v", err))
		}

		consumed := len(blockStatesBytes[offset:]) - buf.Len()
		if consumed <= 0 {
			panic("DecodeBlockStates: decoder consumed no bytes")
		}
		offset += consumed

		hash := block.ComputeBlockHash(s.Name, s.Properties)
		if _, has := existHash[hash]; has {
			continue
		}
		existHash[hash] = struct{}{}
		blockStates = append(blockStates, s)
	}
	return blockStates
}

func decompressBlockStates(blockStatesBytes []byte) []byte {
	if len(blockStatesBytes) < 2 || blockStatesBytes[0] != 0x1f || blockStatesBytes[1] != 0x8b {
		return blockStatesBytes
	}

	gzipReader, err := gzip.NewReader(bytes.NewReader(blockStatesBytes))
	if err != nil {
		panic(fmt.Errorf("DecodeBlockStates: Failed to open gzip-compressed NBT: %v", err))
	}
	decompressed, err := io.ReadAll(gzipReader)
	if closeErr := gzipReader.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		panic(fmt.Errorf("DecodeBlockStates: Failed to read gzip-compressed NBT: %v", err))
	}
	return decompressed
}
