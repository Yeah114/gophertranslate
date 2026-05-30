package main

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	bwo_define "github.com/Yeah114/bedrock-world-operator/define"
	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gophertunnel/minecraft/nbt"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/worlddowngrader/blockdowngrader"
)

func main() {
	out := flag.String("out", "minecraft/world/block/full/full_block_states.nbt", "output full block states NBT path")
	flag.Parse()
	world_block.Init(nil)

	states, err := buildFullBlockStates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen_full_block_states: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	enc := nbt.NewEncoder(&buf)
	for _, state := range states {
		if err := enc.Encode(state); err != nil {
			fmt.Fprintf(os.Stderr, "gen_full_block_states: encode %s: %v\n", state.Name, err)
			os.Exit(1)
		}
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "gen_full_block_states: mkdir: %v\n", err)
		os.Exit(1)
	}
	compressed, err := gzipBytes(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "gen_full_block_states: gzip: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(*out, compressed, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "gen_full_block_states: write: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("generated %s (%d block states)\n", *out, len(states))
}

func gzipBytes(data []byte) ([]byte, error) {
	var compressed bytes.Buffer
	writer, err := gzip.NewWriterLevel(&compressed, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := writer.Write(data); err != nil {
		_ = writer.Close()
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return compressed.Bytes(), nil
}

func buildFullBlockStates() ([]bwo_define.BlockState, error) {
	tableFunc, ok := world_block.Pool[protocol.CurrentProtocol]
	if !ok {
		return nil, fmt.Errorf("no latest block runtime ID table found for protocol %d", protocol.CurrentProtocol)
	}

	table := tableFunc(false)
	table.FinaliseRegister()

	profiles := protocol.Profiles()
	versions := make([]protocol.Profile, 0, len(profiles))
	for _, profile := range profiles {
		if profile.NetEase {
			continue
		}
		versions = append(versions, profile)
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Protocol > versions[j].Protocol
	})

	latestStates := make([]bwo_define.BlockState, 0)
	for runtimeID := uint32(0); ; runtimeID++ {
		name, properties, found := table.RuntimeIDToState(runtimeID)
		if !found {
			break
		}
		latestStates = append(latestStates, bwo_define.BlockState{
			Name:       name,
			Properties: properties,
			Version:    protocol.CurrentProfile.BlockStateVersion(),
		})
	}

	versionStates := make([][]bwo_define.BlockState, len(versions))
	var wg sync.WaitGroup
	for versionIndex, profile := range versions {
		versionIndex, profile := versionIndex, profile
		wg.Add(1)
		go func() {
			defer wg.Done()

			states := make([]bwo_define.BlockState, 0, len(latestStates))
			for _, state := range latestStates {
				downgraded := blockdowngrader.DowngradeToVersion(blockdowngrader.BlockState{
					Name:       state.Name,
					Properties: state.Properties,
					Version:    state.Version,
				}, profile.Version)
				states = append(states, bwo_define.BlockState{
					Name:       downgraded.Name,
					Properties: downgraded.Properties,
					Version:    downgraded.Version,
				})
			}
			versionStates[versionIndex] = states
		}()
	}
	wg.Wait()

	seen := make(map[uint32]struct{})
	fullStates := make([]bwo_define.BlockState, 0, len(latestStates))
	for _, state := range latestStates {
		addUniqueBlockState(&fullStates, seen, state)
	}
	for _, states := range versionStates {
		for _, state := range states {
			addUniqueBlockState(&fullStates, seen, state)
		}
	}
	return fullStates, nil
}

func addUniqueBlockState(states *[]bwo_define.BlockState, seen map[uint32]struct{}, state bwo_define.BlockState) {
	hash := bwo_block.ComputeBlockHash(state.Name, state.Properties)
	if _, ok := seen[hash]; ok {
		return
	}
	seen[hash] = struct{}{}
	*states = append(*states, state)
}
