package minecraft

import (
	"fmt"

	"github.com/Yeah114/gopherconvert/minecraft/utils"
	"github.com/Yeah114/gophertunnel/minecraft/protocol/packet"
)

// HandleCraftingData converts recipes, potion recipes, potion container change recipes and material
// reducers inside a CraftingData packet from server protocol to client protocol.
func (c *MinecraftConverter) HandleCraftingData(pk *packet.CraftingData) error {
	recipes, err := utils.ConvertSliceWithError(pk.Recipes, c.ic.ConvertServerRecipe)
	if err != nil {
		return fmt.Errorf("HandleCraftingData: failed to convert recipes: %w", err)
	}
	potionRecipes, err := utils.ConvertSliceWithError(pk.PotionRecipes, c.ic.ConvertServerPotionRecipe)
	if err != nil {
		return fmt.Errorf("HandleCraftingData: failed to convert potion recipes: %w", err)
	}
	potionContainerChangeRecipes, err := utils.ConvertSliceWithError(pk.PotionContainerChangeRecipes, c.ic.ConvertServerPotionContainerChangeRecipe)
	if err != nil {
		return fmt.Errorf("HandleCraftingData: failed to convert potion container change recipes: %w", err)
	}
	materialReducers, err := utils.ConvertSliceWithError(pk.MaterialReducers, c.ic.ConvertServerMaterialReducer)
	if err != nil {
		return fmt.Errorf("HandleCraftingData: failed to convert material reducers: %w", err)
	}
	dst := *pk
	dst.Recipes = recipes
	dst.PotionRecipes = potionRecipes
	dst.PotionContainerChangeRecipes = potionContainerChangeRecipes
	dst.MaterialReducers = materialReducers
	return c.clientConnEcho.WritePacket(&dst)
}
