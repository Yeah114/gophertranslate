package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

type versionEntry struct {
	proto       int32
	versionName string // e.g. "v1v26v20"
	versionStr  string // e.g. "1.26.20"
	versionInt  int32  // e.g. 18486272 from Version()
}

type blockStateRecord struct {
	Name       string
	Properties map[string]any
}

func main() {
	block.Init(nil)

	var versions []versionEntry
	for proto := range block.Pool {
		info := protocol.Info(proto)
		verStr := info.Ver()
		verName := versionStrToDirName(verStr)
		verInt := info.Version()
		versions = append(versions, versionEntry{
			proto:       proto,
			versionName: verName,
			versionStr:  verStr,
			versionInt:  verInt,
		})
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].proto < versions[j].proto
	})

	fmt.Printf("Found %d versions in block pool\n", len(versions))
	for _, v := range versions {
		fmt.Printf("  %s (proto=%d, dir=%s)\n", v.versionStr, v.proto, v.versionName)
	}

	outDir := `D:\gopherconvert\minecraft\version`

	// First version gets empty BlockStates + converter
	{
		first := versions[0]
		genBlockStateFile(outDir, first.versionName, nil, first.versionInt)
		genConverterFile(outDir, first.versionName, first.versionStr)
		fmt.Printf("  First version %s: converter generated\n", first.versionStr)
	}

	// For each consecutive pair, compute diff and generate files
	for i := 1; i < len(versions); i++ {
		cur := versions[i]
		prev := versions[i-1]

		// v1v26v20v26 is an alias for v1v26v20 - same NBT, same converter folder
		if cur.versionName == "v1v26v20v26" {
			fmt.Printf("\n--- Skip %s (alias of v1v26v20) ---\n", cur.versionStr)
			continue
		}

		// For v1v26v20, previous is v1v26v10 (skip the v1v26v20v26 alias)
		var actualPrev versionEntry
		if cur.versionName == "v1v26v20" {
			actualPrev = prev
			// If prev was v1v26v20v26, skip back to v1v26v10
			if prev.versionName == "v1v26v20v26" {
				actualPrev = versions[i-2]
			}
		} else {
			actualPrev = prev
			// If prev is v1v26v20v26 alias, skip it
			if prev.versionName == "v1v26v20v26" {
				actualPrev = versions[i-2]
			}
		}

		fmt.Printf("\n--- Diff %s -> %s ---\n", actualPrev.versionStr, cur.versionStr)

		newStates := loadBlockStates(cur.versionName)
		oldStates := loadBlockStates(actualPrev.versionName)

		diffStates := computeDiffV2(newStates, oldStates)
		fmt.Printf("  New blocks in %s: %d\n", cur.versionStr, len(diffStates))

		// For block_state.go, always write (even if empty, block_state.go has BlockStates = []define.BlockState{})
		genBlockStateFile(outDir, cur.versionName, diffStates, cur.versionInt)

		// For converter.go, use cur.versionName as directory name
		// Special case: v1v26v20 already exists with crafting handling, skip
		if cur.versionName != "v1v26v20" {
			genConverterFile(outDir, cur.versionName, cur.versionStr)
		} else {
			fmt.Printf("  Skipped converter for %s (already exists)\n", cur.versionName)
		}
	}

	// Update pool.go
	genPoolFile(outDir, versions)

	fmt.Println("\nDone!")
}

func versionStrToDirName(s string) string {
	parts := strings.Split(s, ".")
	var out strings.Builder
	out.WriteString("v")
	for i, p := range parts {
		if i > 0 {
			out.WriteString("v")
		}
		out.WriteString(p)
	}
	return out.String()
}

// loadBlockStates reads block_states.nbt from a block version directory and returns
// deduplicated block states (by name+properties). It avoids creating BlockRuntimeIDTable
// which panics on duplicates.
func loadBlockStates(versionName string) []blockStateRecord {
	nbtPath := filepath.Join(`D:\gopherconvert\minecraft\world\block`, versionName, "block_states.nbt")
	data, err := os.ReadFile(nbtPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read %s: %v", nbtPath, err))
	}
	rawStates := utils.DecodeBlockStates(data)
	seen := make(map[string]struct{})
	var unique []blockStateRecord
	for _, s := range rawStates {
		key := stateKey(s.Name, s.Properties)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, blockStateRecord{
			Name:       s.Name,
			Properties: cloneProperties(s.Properties),
		})
	}
	return unique
}

func computeDiffV2(newStates, oldStates []blockStateRecord) []blockStateRecord {
	oldSet := make(map[string]struct{})
	for _, s := range oldStates {
		oldSet[stateKey(s.Name, s.Properties)] = struct{}{}
	}

	var diff []blockStateRecord
	for _, s := range newStates {
		if _, ok := oldSet[stateKey(s.Name, s.Properties)]; !ok {
			diff = append(diff, s)
		}
	}

	sort.Slice(diff, func(i, j int) bool {
		return stateKey(diff[i].Name, diff[i].Properties) < stateKey(diff[j].Name, diff[j].Properties)
	})
	return diff
}

func stateKey(name string, properties map[string]any) string {
	return name + "\x00" + formatPropertiesKey(properties)
}

func formatPropertiesKey(properties map[string]any) string {
	if len(properties) == 0 {
		return ""
	}
	keys := make([]string, 0, len(properties))
	for key := range properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, key := range keys {
		if i > 0 {
			b.WriteByte(';')
		}
		fmt.Fprintf(&b, "%s=%v", key, properties[key])
	}
	return b.String()
}

func cloneProperties(properties map[string]any) map[string]any {
	if len(properties) == 0 {
		return nil
	}
	cloned := make(map[string]any, len(properties))
	for key, value := range properties {
		cloned[key] = value
	}
	return cloned
}

func formatPropertiesLiteral(properties map[string]any) string {
	if len(properties) == 0 {
		return "nil"
	}
	keys := make([]string, 0, len(properties))
	for key := range properties {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var b strings.Builder
	b.WriteString("map[string]any{")
	for i, key := range keys {
		if i > 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%q: %s", key, formatValueLiteral(properties[key]))
	}
	b.WriteByte('}')
	return b.String()
}

func formatValueLiteral(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int8:
		return fmt.Sprintf("int8(%d)", v)
	case int16:
		return fmt.Sprintf("int16(%d)", v)
	case int32:
		return fmt.Sprintf("int32(%d)", v)
	case int64:
		return fmt.Sprintf("int64(%d)", v)
	case int:
		return fmt.Sprintf("int(%d)", v)
	case uint8:
		return fmt.Sprintf("uint8(%d)", v)
	case uint16:
		return fmt.Sprintf("uint16(%d)", v)
	case uint32:
		return fmt.Sprintf("uint32(%d)", v)
	case uint64:
		return fmt.Sprintf("uint64(%d)", v)
	case uint:
		return fmt.Sprintf("uint(%d)", v)
	case float32:
		return fmt.Sprintf("float32(%v)", v)
	case float64:
		return fmt.Sprintf("float64(%v)", v)
	default:
		return fmt.Sprintf("%#v", v)
	}
}

func genBlockStateFile(outDir, versionDir string, states []blockStateRecord, versionInt int32) {
	dirPath := filepath.Join(outDir, versionDir)
	os.MkdirAll(dirPath, 0755)

	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", versionDir))
	b.WriteString(`import "github.com/Yeah114/bedrock-world-operator/define"` + "\n\n")
	b.WriteString("var BlockStates = []define.BlockState{\n")
	for _, state := range states {
		b.WriteString("\t{\n")
		b.WriteString(fmt.Sprintf("\t\tName:       %q,\n", state.Name))
		b.WriteString(fmt.Sprintf("\t\tProperties: %s,\n", formatPropertiesLiteral(state.Properties)))
		b.WriteString(fmt.Sprintf("\t\tVersion:    int32(%d),\n", versionInt))
		b.WriteString("\t},\n")
	}
	b.WriteString("}\n")

	outPath := filepath.Join(dirPath, "block_state.go")
	err := os.WriteFile(outPath, []byte(b.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
	} else {
		fmt.Printf("  Generated %s (%d block states)\n", outPath, len(states))
	}
}

func genConverterFile(outDir, versionDir, versionStr string) {
	dirPath := filepath.Join(outDir, versionDir)
	os.MkdirAll(dirPath, 0755)

	outPath := filepath.Join(dirPath, "converter.go")
	if _, err := os.Stat(outPath); err == nil {
		fmt.Printf("  Skipped %s (already exists)\n", outPath)
		return
	}

	// Build the content by appending parts
	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\n", versionDir)
	// The %% in WriteString becomes literal % in output. %%w -> %w in generated Go code.
	b.WriteString(`import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// VersionConverter applies the protocol-specific adjustments required for
// Minecraft `)
	b.WriteString(versionStr)
	b.WriteString(`.
type VersionConverter struct {
	c define.MinecraftConverter
}

// NewVersionConverter creates a version converter for Minecraft `)
	b.WriteString(versionStr)
	b.WriteString(`.
func NewVersionConverter(c define.MinecraftConverter) define.VersionConverter {
	return &VersionConverter{c: c}
}

// StartGame registers the custom block states required by this protocol
// version before the game starts.
func (c *VersionConverter) StartGame(data *minecraft.GameData) (err error) {
	if c.c.ClientConnEcho().Proto().ID() < c.c.ServerConnEcho().Proto().ID() {
		data.CustomBlocks = append(data.CustomBlocks, utils.BlockStatesToBlockEntries(BlockStates)...)
		table := c.c.BlockConverter().BlockConverter().ClientTable()
		for _, state := range BlockStates {
			err := table.RegisterCustomBlock(state)
			if err != nil {
				return fmt.Errorf("`)
	fmt.Fprintf(&b, "%s", versionDir)
	b.WriteString(`.VersionConverter.StartGame: failed to register custom block: %w", err)
			}
		}
		table.FinaliseRegister()
	}
	return nil
}

// HandlePacket processes echo packets from the main converter, applying version-specific
// transformations where needed.
func (c *VersionConverter) HandlePacket(pk packet.Packet, sender define.Conn) (err error) {
	if sender == c.c.ServerConn() {
		return c.c.ClientConnEcho().WritePacket(pk)
	}
	if sender == c.c.ClientConn() {
		return c.c.ServerConnEcho().WritePacket(pk)
	}
	return fmt.Errorf("`)
	fmt.Fprintf(&b, "%s", versionDir)
	b.WriteString(`.VersionConverter.HandlePacket: unknown sender")
}
`)

	err := os.WriteFile(outPath, []byte(b.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
	} else {
		fmt.Printf("  Generated %s\n", outPath)
	}
}

func genPoolFile(outDir string, versions []versionEntry) {
	var b strings.Builder
	b.WriteString("package version\n\n")
	b.WriteString("import (\n")
	b.WriteString("\t\"sort\"\n\n")
	b.WriteString("\t\"github.com/Yeah114/gopherconvert/minecraft/define\"\n")

	seen := make(map[string]bool)
	for _, v := range versions {
		name := v.versionName
		if v.versionStr == "1.26.20.26" {
			name = "v1v26v20"
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		b.WriteString(fmt.Sprintf("\t\"github.com/Yeah114/gopherconvert/minecraft/version/%s\"\n", name))
	}

	b.WriteString("\t\"github.com/Yeah114/gophertunnel/minecraft/protocol\"\n")
	b.WriteString(")\n\n")

	b.WriteString("var Pool = map[int32]func(define.MinecraftConverter) define.VersionConverter{\n")
	for _, v := range versions {
		pkgName := v.versionName
		if v.versionStr == "1.26.20.26" {
			pkgName = "v1v26v20"
		}
		// versionName "v1v16v100" -> protocol const Protocol1v16v100, so strip leading "v"
		protoSuffix := v.versionName[1:]
		b.WriteString(fmt.Sprintf("\tprotocol.Protocol%s: %s.NewVersionConverter,\n", protoSuffix, pkgName))
	}
	b.WriteString("}\n\n")

	b.WriteString(`func GetVersionConverters(sourceProto, targetProto int32) []func(define.MinecraftConverter) define.VersionConverter {
	type pair struct {
		proto int32
		ctor  func(define.MinecraftConverter) define.VersionConverter
	}
	var list []pair
	minProto, maxProto := min(sourceProto, targetProto), max(sourceProto, targetProto)

	for proto, ctor := range Pool {
		if proto > minProto && proto <= maxProto {
			list = append(list, pair{proto: proto, ctor: ctor})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].proto > list[j].proto
	})

	var res []func(define.MinecraftConverter) define.VersionConverter
	for _, item := range list {
		res = append(res, item.ctor)
	}
	return res
}
`)

	outPath := filepath.Join(outDir, "pool.go")
	err := os.WriteFile(outPath, []byte(b.String()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", outPath, err)
	} else {
		fmt.Printf("  Generated %s (%d entries)\n", outPath, len(versions))
	}
}
