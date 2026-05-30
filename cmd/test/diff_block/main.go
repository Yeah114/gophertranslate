package main

import (
	"fmt"
	"sort"
	"strings"

	bwo_block "github.com/Yeah114/bedrock-world-operator/block"
	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v10"
	"github.com/Yeah114/gopherconvert/minecraft/world/block/v1v26v20"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
)

const (
	leftVersion  = "1.26.20"
	rightVersion = "1.26.10"
)

type blockState struct {
	Name       string
	Properties map[string]any
}

func main() {
	world_block.Init(nil)

	v10 := v1v26v10.NewBlockRuntimeIDTable(false)
	v20 := v1v26v20.NewBlockRuntimeIDTable(false)
	leftProtocol, _ := protocol.GetProtocol(leftVersion)
	rightProtocol, _ := protocol.GetProtocol(rightVersion)
	leftProfile, _ := protocol.GetProfile(leftProtocol)
	_, _ = protocol.GetProfile(rightProtocol)
	version := leftProfile.BlockStateVersion()

	onlyIn20 := diffStates(scanStates(v20), scanStates(v10))

	fmt.Println("[]define.BlockState{")
	for _, state := range onlyIn20 {
		fmt.Printf("\t{\n")
		fmt.Printf("\t\tName:       %q,\n", state.Name)
		fmt.Printf("\t\tProperties: %s,\n", formatPropertiesLiteral(state.Properties))
		fmt.Printf("\t\tVersion:    int32(%d),\n", version)
		fmt.Printf("\t},\n")
	}
	fmt.Println("}")
}

func scanStates(table *bwo_block.BlockRuntimeIDTable) []blockState {
	var states []blockState
	for runtimeID := uint32(0); ; runtimeID++ {
		name, properties, found := table.RuntimeIDToState(runtimeID)
		if !found {
			break
		}
		states = append(states, blockState{
			Name:       name,
			Properties: cloneProperties(properties),
		})
	}
	return states
}

func diffStates(left, right []blockState) []blockState {
	rightSet := make(map[string]struct{}, len(right))
	for _, state := range right {
		rightSet[stateKey(state.Name, state.Properties)] = struct{}{}
	}
	var diff []blockState
	for _, state := range left {
		if _, ok := rightSet[stateKey(state.Name, state.Properties)]; !ok {
			diff = append(diff, state)
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
