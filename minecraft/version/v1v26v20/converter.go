package v1v26v20

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/define"
	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft"
	"github.com/Yeah114/gophertunnel/minecraft/protocol"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
	"github.com/google/uuid"
)

// VersionConverter applies the protocol-specific adjustments required for
// Minecraft 1.26.20.
type VersionConverter struct {
	c define.MinecraftConverter
}

// NewVersionConverter creates a version converter for Minecraft 1.26.20.
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
				return fmt.Errorf("v1v26v20.VersionConverter.StartGame: failed to register custom block: %w", err)
			}
		}
	}

	return nil
}

// HandlePacket processes echo packets from the main converter, applying version-specific
// transformations where needed.
func (c *VersionConverter) HandlePacket(pk packet.Packet, sender define.Conn) (err error) {
	if sender == c.c.ServerConn() {
		switch pkt := pk.(type) {
		case *packet.CraftingData:
			return c.HandleCraftingData(pkt)
		default:
			return c.c.ClientConnEcho().WritePacket(pk)
		}
	}
	if sender == c.c.ClientConn() {
		return c.c.ServerConnEcho().WritePacket(pk)
	}
	return fmt.Errorf("v1.26v20.VersionConverter.HandlePacket: unknown sender")
}

// HandleCraftingData downgrades furnace recipes to shapeless recipes for an older client
// that does not understand native FurnaceRecipe/FurnaceDataRecipe types. If the client is not older
// than the server, the packet is forwarded unchanged.
func (c *VersionConverter) HandleCraftingData(pk *packet.CraftingData) error {
	if c.c.ClientConnEcho().Proto().ID() >= c.c.ServerConnEcho().Proto().ID() {
		return c.c.ClientConnEcho().WritePacket(pk)
	}
	for i, recipe := range pk.Recipes {
		dst, err := c.ConvertServerFurnaceRecipe(recipe, uint32(i)+1)
		if err != nil {
			continue
		}
		pk.Recipes[i] = dst
	}
	return c.c.ClientConnEcho().WritePacket(pk)
}

// ConvertServerFurnaceRecipe converts a legacy furnace recipe (FurnaceRecipe or FurnaceDataRecipe)
// to a ShapelessRecipe compatible with older protocol versions. This is required because
// protocol versions prior to 1.27.x do not support the dedicated furnace recipe packet fields.
func (c *VersionConverter) ConvertServerFurnaceRecipe(recipe protocol.Recipe, recipeNetworkID uint32) (protocol.Recipe, error) {
	switch recipe := recipe.(type) {
	case *protocol.FurnaceRecipe:
		return c.ConvertServerFurnaceRecipeWithData(recipe, false, recipeNetworkID)
	case *protocol.FurnaceDataRecipe:
		return c.ConvertServerFurnaceRecipeWithData(&recipe.FurnaceRecipe, true, recipeNetworkID)
	default:
		return nil, fmt.Errorf("ConvertServerFurnaceRecipe: not a furnace recipe")
	}
}

// ConvertServerFurnaceRecipeWithData builds a ShapelessRecipe from a furnace recipe's fields.
// When hasData is true, the recipe's InputType metadata is used; otherwise a wildcard
// metadata (0x7fff) is used.
func (c *VersionConverter) ConvertServerFurnaceRecipeWithData(recipe *protocol.FurnaceRecipe, hasData bool, recipeNetworkID uint32) (protocol.Recipe, error) {
	metadata := int16(0x7fff)
	if hasData {
		metadata = int16(recipe.InputType.MetadataValue)
	}
	id := uuid.New()
	return &protocol.ShapelessRecipe{
		RecipeID: id.String(),
		Input: []protocol.ItemDescriptorCount{{
			Descriptor: &protocol.DefaultItemDescriptor{
				NetworkID:     int16(recipe.InputType.NetworkID),
				MetadataValue: metadata,
			},
			Count: 1,
		}},
		Output: []protocol.ItemStack{recipe.Output},
		UUID:   id,
		Block:  recipe.Block,
		UnlockRequirement: protocol.RecipeUnlockRequirement{
			Context: protocol.RecipeUnlockContextAlwaysUnlocked,
		},
		RecipeNetworkID: recipeNetworkID,
	}, nil
}

// ConvertShapelessToFurnaceRecipe converts a ShapelessRecipe to a FurnaceRecipe or
// FurnaceDataRecipe if its block is "furnace" or "blast_furnace". Recipes with a
// MetadataValue other than 0x7fff become FurnaceDataRecipe; otherwise FurnaceRecipe.
func (c *VersionConverter) ConvertShapelessToFurnaceRecipe(recipe protocol.Recipe) protocol.Recipe {
	shapeless, ok := recipe.(*protocol.ShapelessRecipe)
	if !ok || len(shapeless.Input) != 1 || len(shapeless.Output) != 1 {
		return recipe
	}
	desc, ok := shapeless.Input[0].Descriptor.(*protocol.DefaultItemDescriptor)
	if !ok {
		return recipe
	}
	switch shapeless.Block {
	case "furnace", "blast_furnace":
	default:
		return recipe
	}
	furnace := &protocol.FurnaceRecipe{
		InputType: protocol.ItemType{
			NetworkID:     int32(desc.NetworkID),
			MetadataValue: uint32(desc.MetadataValue),
		},
		Output: shapeless.Output[0],
		Block:  shapeless.Block,
	}
	if desc.MetadataValue == 0x7fff {
		return furnace
	}
	return &protocol.FurnaceDataRecipe{FurnaceRecipe: *furnace}
}
