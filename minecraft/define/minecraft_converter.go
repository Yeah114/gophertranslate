package define

import (
	"context"

	bwo_define "github.com/Yeah114/bedrock-world-operator/define"
	world_block "github.com/Yeah114/gopherconvert/minecraft/world/block"
	world_item "github.com/Yeah114/gopherconvert/minecraft/world/item"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// BlockConverter converts block runtime IDs and block-related structures between client and server protocols.
type BlockConverter interface {
	BlockConverter() *world_block.BlockConverter
	ConvertClientBlockRuntimeID(clientBlockRuntimeID uint32) (uint32, error)
	ConvertClientBlockRuntimeIDInt32(clientBlockRuntimeID int32) (int32, error)
	ConvertClientItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error)
	ConvertClientItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error)
	ConvertClientCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error)
	ConvertClientCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error)
	ConvertClientInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error)
	ConvertClientBlockChangeEntry(clientBlockChangeEntry protocol.BlockChangeEntry) (protocol.BlockChangeEntry, error)
	ConvertClientUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error)
	ConvertClientUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error)
	ConvertClientReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error)
	ConvertClientWaxedOrUnwaxedCopperEvent(clientEvent *protocol.WaxedOrUnwaxedCopperEvent) (*protocol.WaxedOrUnwaxedCopperEvent, error)
	ConvertClientEvent(clientEvent protocol.Event) (protocol.Event, error)
	ConvertClientBiomeDefinition(clientBiomeDefinition protocol.BiomeDefinition) (protocol.BiomeDefinition, error)
	ConvertClientBiomeChunkGeneration(clientChunkGeneration protocol.BiomeChunkGeneration) (protocol.BiomeChunkGeneration, error)
	ConvertClientBiomeMountainParameters(clientParameters protocol.BiomeMountainParameters) (protocol.BiomeMountainParameters, error)
	ConvertClientBiomeElementData(clientElementData protocol.BiomeElementData) (protocol.BiomeElementData, error)
	ConvertClientBiomeSurfaceMaterial(clientSurfaceMaterial protocol.BiomeSurfaceMaterial) (protocol.BiomeSurfaceMaterial, error)
	ConvertClientBiomeMesaSurface(clientMesaSurface protocol.BiomeMesaSurface) (protocol.BiomeMesaSurface, error)
	ConvertClientBiomeCappedSurface(clientCappedSurface protocol.BiomeCappedSurface) (protocol.BiomeCappedSurface, error)
	ConvertServerBlockRuntimeID(serverBlockRuntimeID uint32) (uint32, error)
	ConvertServerBlockRuntimeIDInt32(serverBlockRuntimeID int32) (int32, error)
	ConvertServerItemStack(serverItemStack protocol.ItemStack) (protocol.ItemStack, error)
	ConvertServerItemInstance(serverItemInstance protocol.ItemInstance) (protocol.ItemInstance, error)
	ConvertServerCreativeGroup(serverCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error)
	ConvertServerCreativeItem(serverCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error)
	ConvertServerInventoryAction(serverInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error)
	ConvertServerBlockChangeEntry(serverBlockChangeEntry protocol.BlockChangeEntry) (protocol.BlockChangeEntry, error)
	ConvertServerUseItemTransactionData(serverData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error)
	ConvertServerUseItemOnEntityTransactionData(serverData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error)
	ConvertServerReleaseItemTransactionData(serverData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error)
	ConvertServerWaxedOrUnwaxedCopperEvent(serverEvent *protocol.WaxedOrUnwaxedCopperEvent) (*protocol.WaxedOrUnwaxedCopperEvent, error)
	ConvertServerEvent(serverEvent protocol.Event) (protocol.Event, error)
	ConvertServerBiomeDefinition(serverBiomeDefinition protocol.BiomeDefinition) (protocol.BiomeDefinition, error)
	ConvertServerBiomeChunkGeneration(serverChunkGeneration protocol.BiomeChunkGeneration) (protocol.BiomeChunkGeneration, error)
	ConvertServerBiomeMountainParameters(serverParameters protocol.BiomeMountainParameters) (protocol.BiomeMountainParameters, error)
	ConvertServerBiomeElementData(serverElementData protocol.BiomeElementData) (protocol.BiomeElementData, error)
	ConvertServerBiomeSurfaceMaterial(serverSurfaceMaterial protocol.BiomeSurfaceMaterial) (protocol.BiomeSurfaceMaterial, error)
	ConvertServerBiomeMesaSurface(serverMesaSurface protocol.BiomeMesaSurface) (protocol.BiomeMesaSurface, error)
	ConvertServerBiomeCappedSurface(serverCappedSurface protocol.BiomeCappedSurface) (protocol.BiomeCappedSurface, error)
}

// ChunkConverter converts chunk payloads between client and server protocols.
type ChunkConverter interface {
	ConvertSubChunkEntryRawPayload(clientSubChunkEntryRawPayload []byte, r bwo_define.Range) ([]byte, error)
	ConvertSubChunkBlobPayload(clientSubChunkBlobPayload []byte, r bwo_define.Range) ([]byte, error)
	ConvertSubChunkEntry(clientSubChunkEntry protocol.SubChunkEntry, r bwo_define.Range, cacheEnabled bool) (protocol.SubChunkEntry, error)
	ConvertSubChunk(clientSubChunk *packet.SubChunk) (*packet.SubChunk, error)
	ConvertLevelChunkRawPayload(clientLevelChunkRawPayload []byte, subChunkCount uint32, r bwo_define.Range) ([]byte, error)
	ConvertLevelChunk(clientLevelChunk *packet.LevelChunk) (*packet.LevelChunk, error)
	ConvertCacheBlob(clientBlob protocol.CacheBlob) (protocol.CacheBlob, error)
	ConvertClientCacheMissResponse(clientClientCacheMissResponse *packet.ClientCacheMissResponse) (*packet.ClientCacheMissResponse, error)
}

// ItemConverter converts item runtime IDs and item-related structures between client and server protocols.
type ItemConverter interface {
	ItemConverter() *world_item.ItemConverter
	ConvertItemRuntimeID(clientItemRuntimeID int32) (int32, error)
	ConvertClientItemRuntimeID(clientItemRuntimeID int32) (int32, error)
	ConvertServerItemRuntimeID(serverItemRuntimeID int32) (int32, error)
	ConvertItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error)
	ConvertClientItemStack(clientItemStack protocol.ItemStack) (protocol.ItemStack, error)
	ConvertServerItemStack(serverItemStack protocol.ItemStack) (protocol.ItemStack, error)
	ConvertItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error)
	ConvertClientItemInstance(clientItemInstance protocol.ItemInstance) (protocol.ItemInstance, error)
	ConvertServerItemInstance(serverItemInstance protocol.ItemInstance) (protocol.ItemInstance, error)
	ConvertCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error)
	ConvertClientCreativeGroup(clientCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error)
	ConvertServerCreativeGroup(serverCreativeGroup protocol.CreativeGroup) (protocol.CreativeGroup, error)
	ConvertCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error)
	ConvertClientCreativeItem(clientCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error)
	ConvertServerCreativeItem(serverCreativeItem protocol.CreativeItem) (protocol.CreativeItem, error)
	ConvertInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error)
	ConvertClientInventoryAction(clientInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error)
	ConvertServerInventoryAction(serverInventoryAction protocol.InventoryAction) (protocol.InventoryAction, error)
	ConvertClientInventoryTransactionData(clientData protocol.InventoryTransactionData) (protocol.InventoryTransactionData, error)
	ConvertUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error)
	ConvertClientUseItemTransactionData(clientData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error)
	ConvertServerUseItemTransactionData(serverData *protocol.UseItemTransactionData) (*protocol.UseItemTransactionData, error)
	ConvertUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error)
	ConvertClientUseItemOnEntityTransactionData(clientData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error)
	ConvertServerUseItemOnEntityTransactionData(serverData *protocol.UseItemOnEntityTransactionData) (*protocol.UseItemOnEntityTransactionData, error)
	ConvertReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error)
	ConvertClientReleaseItemTransactionData(clientData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error)
	ConvertServerReleaseItemTransactionData(serverData *protocol.ReleaseItemTransactionData) (*protocol.ReleaseItemTransactionData, error)
	ConvertServerRecipe(serverRecipe protocol.Recipe) (protocol.Recipe, error)
	ConvertServerFurnaceRecipe(recipe *protocol.FurnaceRecipe) (*protocol.FurnaceRecipe, error)
	ConvertServerFurnaceDataRecipe(recipe *protocol.FurnaceDataRecipe) (*protocol.FurnaceDataRecipe, error)
	ConvertServerShapelessRecipe(recipe *protocol.ShapelessRecipe) (*protocol.ShapelessRecipe, error)
	ConvertServerShulkerBoxRecipe(recipe *protocol.ShulkerBoxRecipe) (*protocol.ShulkerBoxRecipe, error)
	ConvertServerShapelessChemistryRecipe(recipe *protocol.ShapelessChemistryRecipe) (*protocol.ShapelessChemistryRecipe, error)
	ConvertServerShapedRecipe(recipe *protocol.ShapedRecipe) (*protocol.ShapedRecipe, error)
	ConvertServerShapedChemistryRecipe(recipe *protocol.ShapedChemistryRecipe) (*protocol.ShapedChemistryRecipe, error)
	ConvertServerSmithingTransformRecipe(recipe *protocol.SmithingTransformRecipe) (*protocol.SmithingTransformRecipe, error)
	ConvertServerSmithingTrimRecipe(recipe *protocol.SmithingTrimRecipe) (*protocol.SmithingTrimRecipe, error)
	ConvertServerItemDescriptorCount(item protocol.ItemDescriptorCount) (protocol.ItemDescriptorCount, error)
	ConvertServerItemDescriptor(descriptor protocol.ItemDescriptor) (protocol.ItemDescriptor, error)
	ConvertServerPotionRecipe(serverRecipe protocol.PotionRecipe) (protocol.PotionRecipe, error)
	ConvertServerPotionContainerChangeRecipe(serverRecipe protocol.PotionContainerChangeRecipe) (protocol.PotionContainerChangeRecipe, error)
	ConvertServerMaterialReducer(serverMaterialReducer protocol.MaterialReducer) (protocol.MaterialReducer, error)
	ConvertServerMaterialReducerOutput(output protocol.MaterialReducerOutput) (protocol.MaterialReducerOutput, error)
}

// MinecraftConverter converts packets from a source connection to a destination connection.
type MinecraftConverter interface {
	// ClientConn returns the source connection packets are read from.
	ClientConn() Conn
	// ClientConnEcho returns the echo connection for packets from the source side.
	ClientConnEcho() Conn
	// ServerConn returns the destination connection packets are written to.
	ServerConn() Conn
	// ServerConnEcho returns the echo connection for packets from the destination side.
	ServerConnEcho() Conn
	// BlockConverter returns the protocol block converter.
	BlockConverter() BlockConverter
	// ChunkConverter returns the protocol chunk converter.
	ChunkConverter() ChunkConverter
	// ItemConverter returns the protocol item converter.
	ItemConverter() ItemConverter
	StartGameContext(ctx context.Context, data *minecraft.GameData) error
	// HandlePacket converts a packet from sender when needed and writes it to the opposite connection.
	HandlePacket(pk packet.Packet, sender Conn) error
	RunDaemonConverter(dc DaemonConverter) bool
	GetDaemonConverter(name string) (DaemonConverter, bool)
	StopDaemonConverter(name string) bool
}

type VersionConverter interface {
	StartGame(data *minecraft.GameData) error
	HandlePacket(pk packet.Packet, sender Conn) error
}

type DaemonConverter interface {
	Name() string
	ProtocolInfo() protocol.Info
	Start()
	Stop()
}
