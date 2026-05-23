package block

import (
	"github.com/Yeah114/bedrock-world-operator/block"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/worlddowngrader/blockdowngrader"
	"github.com/Yeah114/worldupgrader/blockupgrader"
)

// BlockConverter is a struct that can convert block states and runtime IDs between two protocol versions.
type BlockConverter struct {
	clientInfo  protocol.Info
	serverInfo  protocol.Info
	clientTable *block.BlockRuntimeIDTable
	serverTable *block.BlockRuntimeIDTable
}

// NewBlockConverter creates a new BlockConverter that can convert block states and runtime IDs between two protocol versions.
func NewBlockConverter(clientProtocol int32, clientTable *block.BlockRuntimeIDTable, serverProtocol int32, serverTable *block.BlockRuntimeIDTable) *BlockConverter {
	clientInfo := protocol.NewInfoByProtocol(clientProtocol)
	serverInfo := protocol.NewInfoByProtocol(serverProtocol)

	return &BlockConverter{
		clientInfo:  clientInfo,
		serverInfo:  serverInfo,
		clientTable: clientTable,
		serverTable: serverTable,
	}
}

// ServerInfo returns the protocol info of the destination protocol version.
func (c *BlockConverter) ClientInfo() protocol.Info {
	return c.clientInfo
}

// ServerInfo returns the protocol info of the destination protocol version.
func (c *BlockConverter) ServerInfo() protocol.Info {
	return c.serverInfo
}

// ClientTable returns the block runtime ID table of the source protocol version.
func (c *BlockConverter) ClientTable() *block.BlockRuntimeIDTable {
	return c.clientTable
}

// ServerTable returns the block runtime ID table of the destination protocol version.
func (c *BlockConverter) ServerTable() *block.BlockRuntimeIDTable {
	return c.serverTable
}

// ConvertClientBlockState converts a client block state to the server protocol version.
// It returns the converted block state and a boolean indicating whether the conversion was successful.
func (c *BlockConverter) ConvertClientBlockState(name string, properties map[string]interface{}) (string, map[string]interface{}, bool) {
	if c.clientInfo.Version() < c.serverInfo.Version() {
		blockState := blockupgrader.BlockState{
			Name:       name,
			Properties: properties,
			Version:    c.clientInfo.Version(),
		}
		serverBlockState := blockupgrader.UpgradeToVersion(blockState, c.serverInfo.Ver())
		return serverBlockState.Name, serverBlockState.Properties, true
	} else if c.clientInfo.Version() > c.serverInfo.Version() {
		blockState := blockdowngrader.BlockState{
			Name:       name,
			Properties: properties,
			Version:    c.clientInfo.Version(),
		}
		serverBlockState := blockdowngrader.DowngradeToVersion(blockState, c.serverInfo.Ver())
		return serverBlockState.Name, serverBlockState.Properties, true
	}
	return name, properties, true
}

// ConvertServerBlockState converts a server block state to the client protocol version.
// It returns the converted block state and a boolean indicating whether the conversion was successful.
func (c *BlockConverter) ConvertServerBlockState(name string, properties map[string]interface{}) (string, map[string]interface{}, bool) {
	if c.clientInfo.Version() < c.serverInfo.Version() {
		blockState := blockdowngrader.BlockState{
			Name:       name,
			Properties: properties,
			Version:    c.serverInfo.Version(),
		}
		clientBlockState := blockdowngrader.DowngradeToVersion(blockState, c.clientInfo.Ver())
		return clientBlockState.Name, clientBlockState.Properties, true
	} else if c.clientInfo.Version() > c.serverInfo.Version() {
		blockState := blockupgrader.BlockState{
			Name:       name,
			Properties: properties,
			Version:    c.serverInfo.Version(),
		}
		clientBlockState := blockupgrader.UpgradeToVersion(blockState, c.clientInfo.Ver())
		return clientBlockState.Name, clientBlockState.Properties, true
	}
	return name, properties, true
}

// ConvertClientBlockRuntimeID converts a client block runtime ID to the server protocol version.
// It returns the converted block runtime ID and a boolean indicating whether the conversion was successful.
func (c *BlockConverter) ConvertClientBlockRuntimeID(runtimeID uint32) (uint32, bool) {
	name, properties, found := c.clientTable.RuntimeIDToState(runtimeID)
	if !found {
		return 0, false
	}

	serverName, serverProperties, found := c.ConvertClientBlockState(name, properties)
	if !found {
		return 0, false
	}

	return c.serverTable.StateToRuntimeID(serverName, serverProperties)
}

// ConvertServerBlockRuntimeID converts a server block runtime ID to the client protocol version.
// It returns the converted block runtime ID and a boolean indicating whether the conversion was successful.
func (c *BlockConverter) ConvertServerBlockRuntimeID(runtimeID uint32) (uint32, bool) {
	name, properties, found := c.serverTable.RuntimeIDToState(runtimeID)
	if !found {
		return 0, false
	}

	clientName, clientProperties, found := c.ConvertServerBlockState(name, properties)
	if !found {
		return 0, false
	}

	return c.clientTable.StateToRuntimeID(clientName, clientProperties)
}
