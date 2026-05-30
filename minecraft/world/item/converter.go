package item

import "github.com/Yeah114/gophertunnel/minecraft/protocol"

// ItemConverter converts item runtime IDs between two protocol versions.
type ItemConverter struct {
	clientInfo  protocol.Profile
	serverInfo  protocol.Profile
	clientTable *ItemRuntimeIDTable
	serverTable *ItemRuntimeIDTable
}

// NewItemConverter creates a new ItemConverter.
func NewItemConverter(clientProtocol int32, clientTable *ItemRuntimeIDTable, serverProtocol int32, serverTable *ItemRuntimeIDTable) *ItemConverter {
	clientInfo, _ := protocol.GetProfile(clientProtocol)
	serverInfo, _ := protocol.GetProfile(serverProtocol)
	return &ItemConverter{
		clientInfo:  clientInfo,
		serverInfo:  serverInfo,
		clientTable: clientTable,
		serverTable: serverTable,
	}
}

// ClientInfo returns the protocol info of the source protocol version.
func (c *ItemConverter) ClientInfo() protocol.Profile {
	return c.clientInfo
}

// ServerInfo returns the protocol info of the destination protocol version.
func (c *ItemConverter) ServerInfo() protocol.Profile {
	return c.serverInfo
}

// ClientTable returns the item runtime ID table of the source protocol version.
func (c *ItemConverter) ClientTable() *ItemRuntimeIDTable {
	return c.clientTable
}

// ServerTable returns the item runtime ID table of the destination protocol version.
func (c *ItemConverter) ServerTable() *ItemRuntimeIDTable {
	return c.serverTable
}

// ConvertItemRuntimeID converts an item runtime ID from one protocol version to another.
func (c *ItemConverter) ConvertItemRuntimeID(runtimeID int32) (int32, bool) {
	return c.ConvertClientItemRuntimeID(runtimeID)
}

// ConvertClientItemRuntimeID converts a client item runtime ID to the server protocol version.
func (c *ItemConverter) ConvertClientItemRuntimeID(runtimeID int32) (int32, bool) {
	return c.convertItemRuntimeID(c.clientTable, c.serverTable, runtimeID)
}

// ConvertServerItemRuntimeID converts a server item runtime ID to the client protocol version.
func (c *ItemConverter) ConvertServerItemRuntimeID(runtimeID int32) (int32, bool) {
	return c.convertItemRuntimeID(c.serverTable, c.clientTable, runtimeID)
}

func (c *ItemConverter) convertItemRuntimeID(fromTable, toTable *ItemRuntimeIDTable, runtimeID int32) (int32, bool) {
	if runtimeID == 0 {
		return runtimeID, true
	}
	name, found := fromTable.RuntimeIDToName(runtimeID)
	if !found {
		return 0, false
	}
	return toTable.NameToRuntimeID(name)
}
