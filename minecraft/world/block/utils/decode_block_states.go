package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gophertunnel/minecraft/nbt"
)

// DecodeBlockStates decodes block states from the given byte slice.
func DecodeBlockStates(blockStatesBytes []byte) (blockStates []define.BlockState) {
	dec := nbt.NewDecoder(bytes.NewBuffer(blockStatesBytes))
	existHash := make(map[uint32]struct{})

	for {
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
		hash := block.ComputeBlockHash(s.Name, s.Properties)
		if _, has := existHash[hash]; has {
			continue
		}
		existHash[hash] = struct{}{}
		blockStates = append(blockStates, s)
	}
	return blockStates
}
