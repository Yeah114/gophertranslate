package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/Yeah114/gophertunnel/minecraft/nbt"
)

type protocolInfo struct {
	Proto    int32
	DirName  string
	Protocol string
}

var protocols = []protocolInfo{
	{419, "v1v16v100", "Protocol1v16v100"},
	{428, "v1v16v210", "Protocol1v16v210"},
	{440, "v1v17v0", "Protocol1v17v0"},
	{448, "v1v17v10", "Protocol1v17v10"},
	{465, "v1v17v30", "Protocol1v17v30"},
	{471, "v1v17v40", "Protocol1v17v40"},
	{486, "v1v18v10", "Protocol1v18v10"},
	{503, "v1v18v30", "Protocol1v18v30"},
	{527, "v1v19v0", "Protocol1v19v0"},
	{544, "v1v19v20", "Protocol1v19v20"},
	{560, "v1v19v50", "Protocol1v19v50"},
	{567, "v1v19v60", "Protocol1v19v60"},
	{575, "v1v19v70", "Protocol1v19v70"},
	{582, "v1v19v80", "Protocol1v19v80"},
	{589, "v1v20v0", "Protocol1v20v0"},
	{594, "v1v20v10", "Protocol1v20v10"},
	{618, "v1v20v30", "Protocol1v20v30"},
	{622, "v1v20v40", "Protocol1v20v40"},
	{630, "v1v20v50", "Protocol1v20v50"},
	{649, "v1v20v60", "Protocol1v20v60"},
	{662, "v1v20v70", "Protocol1v20v70"},
	{671, "v1v20v80", "Protocol1v20v80"},
	{685, "v1v21v0", "Protocol1v21v0"},
	{712, "v1v21v20", "Protocol1v21v20"},
	{729, "v1v21v30", "Protocol1v21v30"},
	{748, "v1v21v40", "Protocol1v21v40"},
	{766, "v1v21v50", "Protocol1v21v50"},
	{776, "v1v21v60", "Protocol1v21v60"},
	{786, "v1v21v70", "Protocol1v21v70"},
	{800, "v1v21v80", "Protocol1v21v80"},
	{818, "v1v21v90", "Protocol1v21v90"},
	{827, "v1v21v100", "Protocol1v21v100"},
	{844, "v1v21v110", "Protocol1v21v110"},
}

var alreadyExisting = map[int32]bool{
	898: true, 924: true, 944: true, 974: true, 975: true,
}

func main() {
	datDir := `D:\Nukkit-MOT\src\main\resources`
	blockDir := `D:\gopherconvert\minecraft\world\block`

	var poolEntries []string
	newVersions := 0
	skipped := 0

	sort.Slice(protocols, func(i, j int) bool { return protocols[i].Proto < protocols[j].Proto })

	for _, info := range protocols {
		if alreadyExisting[info.Proto] {
			skipped++
			continue
		}

		datPath := filepath.Join(datDir, fmt.Sprintf("runtime_block_states_%d.dat", info.Proto))
		compressed, err := os.ReadFile(datPath)
		if err != nil {
			fmt.Printf("SKIP proto=%d: file not found: %v\n", info.Proto, err)
			skipped++
			continue
		}

		gz, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			fmt.Printf("SKIP proto=%d: gzip: %v\n", info.Proto, err)
			skipped++
			continue
		}
		decompressed, err := io.ReadAll(gz)
		gz.Close()
		if err != nil {
			fmt.Printf("SKIP proto=%d: read: %v\n", info.Proto, err)
			skipped++
			continue
		}

		var rawList []map[string]any
		if err := nbt.UnmarshalEncoding(decompressed, &rawList, nbt.BigEndian); err != nil {
			fmt.Printf("SKIP proto=%d: NBT decode: %v\n", info.Proto, err)
			skipped++
			continue
		}

		fmt.Printf("proto=%d: %d raw states -> ", info.Proto, len(rawList))

		var outBuf bytes.Buffer
		count := 0
		for _, raw := range rawList {
			name, _ := raw["name"].(string)
			if name == "" {
				continue
			}

			props := make(map[string]any)
			if statesRaw, ok := raw["states"]; ok {
				if sm, ok := statesRaw.(map[string]any); ok {
					props = sm
				}
			}

			version := int32(0)
			if v, ok := raw["version"]; ok {
				switch vt := v.(type) {
				case int32:
					version = vt
				case int64:
					version = int32(vt)
				case float64:
					version = int32(vt)
				}
			}

			state := map[string]any{
				"name":    name,
				"states":  props,
				"version": version,
			}

			var tagBuf bytes.Buffer
			if err := nbt.NewEncoderWithEncoding(&tagBuf, nbt.NetworkLittleEndian).Encode(state); err != nil {
				fmt.Printf("WARN: encode %s: %v\n", name, err)
				continue
			}
			outBuf.Write(tagBuf.Bytes())
			count++
		}

		if count == 0 {
			fmt.Printf("SKIP (no valid states)\n")
			skipped++
			continue
		}

		verDir := filepath.Join(blockDir, info.DirName)
		os.MkdirAll(verDir, 0755)

		nbtPath := filepath.Join(verDir, "block_states.nbt")
		if err := os.WriteFile(nbtPath, outBuf.Bytes(), 0644); err != nil {
			fmt.Printf("ERROR write: %v\n", err)
			continue
		}

		tableContent := fmt.Sprintf(`package %s

import (
	_ "embed"

	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/bedrock-world-operator/define"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
)

var (
	//go:embed block_states.nbt
	blockStatesBytes []byte
	blockStates      []define.BlockState
)

func init() {
	blockStates = utils.DecodeBlockStates(blockStatesBytes)
}

func NewBlockRuntimeIDTable(useNetworkIDHashes bool) *block.BlockRuntimeIDTable {
	return block.NewBlockRuntimeIDTableFromStates(blockStates, useNetworkIDHashes)
}
`, info.DirName)

		tablePath := filepath.Join(verDir, "table.go")
		os.WriteFile(tablePath, []byte(tableContent), 0644)

		fmt.Printf("OK (%d states, %d bytes)\n", count, len(outBuf.Bytes()))
		poolEntries = append(poolEntries, fmt.Sprintf("\tprotocol.%s: %s.NewBlockRuntimeIDTable,", info.Protocol, info.DirName))
		newVersions++
	}

	fmt.Printf("\n=== Summary ===\nNew: %d, Skipped: %d\n", newVersions, skipped)

	if len(poolEntries) > 0 {
		fmt.Printf("\n=== Pool entries to add ===\n")
		for _, e := range poolEntries {
			fmt.Println(e)
		}
	}
}
