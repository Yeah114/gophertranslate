package item

import "github.com/Yeah114/gophertunnel/minecraft/protocol"

// ItemRuntimeIDTable maps item network/runtime IDs to names and back.
// Bedrock uses the item runtime ID from the item registry as the network ID
// found in item stacks, so both names are exposed for call sites.
type ItemRuntimeIDTable struct {
	runtimeIDToName map[int32]string
	nameToRuntimeID map[string]int32
}

// NewItemRuntimeIDTable creates an item runtime ID table.
func NewItemRuntimeIDTable() *ItemRuntimeIDTable {
	table := &ItemRuntimeIDTable{
		runtimeIDToName: make(map[int32]string),
		nameToRuntimeID: make(map[string]int32),
	}
	return table
}

// ReplaceItems replaces all item mappings with the provided entries.
func (t *ItemRuntimeIDTable) ReplaceItems(entries []protocol.ItemEntry) {
	t.runtimeIDToName = make(map[int32]string, len(entries))
	t.nameToRuntimeID = make(map[string]int32, len(entries))
	t.RegisterItems(entries)
}

// RegisterItems registers item entries into the table.
func (t *ItemRuntimeIDTable) RegisterItems(entries []protocol.ItemEntry) {
	for _, entry := range entries {
		t.RegisterItem(entry)
	}
}

// RegisterItem registers an item entry into the table.
func (t *ItemRuntimeIDTable) RegisterItem(entry protocol.ItemEntry) {
	if entry.Name == "" {
		return
	}
	runtimeID := int32(entry.RuntimeID)
	t.runtimeIDToName[runtimeID] = entry.Name
	t.nameToRuntimeID[entry.Name] = runtimeID
}

// RuntimeIDToName returns the item name for a runtime ID.
func (t *ItemRuntimeIDTable) RuntimeIDToName(runtimeID int32) (string, bool) {
	name, ok := t.runtimeIDToName[runtimeID]
	return name, ok
}

// NetworkIDToName returns the item name for a network ID.
func (t *ItemRuntimeIDTable) NetworkIDToName(networkID int32) (string, bool) {
	return t.RuntimeIDToName(networkID)
}

// NameToRuntimeID returns the runtime ID for an item name.
func (t *ItemRuntimeIDTable) NameToRuntimeID(name string) (int32, bool) {
	runtimeID, ok := t.nameToRuntimeID[name]
	return runtimeID, ok
}

// NameToNetworkID returns the network ID for an item name.
func (t *ItemRuntimeIDTable) NameToNetworkID(name string) (int32, bool) {
	return t.NameToRuntimeID(name)
}

// RuntimeIDToNetworkID converts an item runtime ID to its network ID.
func (t *ItemRuntimeIDTable) RuntimeIDToNetworkID(runtimeID int32) (int32, bool) {
	_, ok := t.runtimeIDToName[runtimeID]
	return runtimeID, ok
}

// NetworkIDToRuntimeID converts an item network ID to its runtime ID.
func (t *ItemRuntimeIDTable) NetworkIDToRuntimeID(networkID int32) (int32, bool) {
	return t.RuntimeIDToNetworkID(networkID)
}
