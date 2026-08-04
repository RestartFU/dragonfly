package item

import (
	"testing"

	"github.com/df-mc/dragonfly/server/world"
)

// testEnchantment is registered locally because server/item/enchantment imports this
// package, so a test inside it cannot pull the real types in.
type testEnchantment struct{}

func (testEnchantment) Name() string                                   { return "Test" }
func (testEnchantment) MaxLevel() int                                  { return 5 }
func (testEnchantment) Cost(int) (int, int)                            { return 1, 10 }
func (testEnchantment) Rarity() EnchantmentRarity                      { return EnchantmentRarityCommon }
func (testEnchantment) CompatibleWithEnchantment(EnchantmentType) bool { return true }
func (testEnchantment) CompatibleWithItem(world.Item) bool             { return true }

// Level 0 reaches this from malformed or hostile NBT. Enchantments are 1-indexed, so it is
// clamped rather than stored as written.
func TestReadEnchantments_ClampsLevelBelowOne(t *testing.T) {
	const id = 9999
	RegisterEnchantment(id, testEnchantment{})

	stack := NewStack(Sword{Tier: ToolTierDiamond}, 1)
	readEnchantments(map[string]any{
		"ench": []any{
			map[string]any{"id": int16(id), "lvl": int16(0)},
		},
	}, &stack)

	enchants := stack.Enchantments()
	if len(enchants) != 1 {
		t.Fatalf("expected 1 enchantment, got %d", len(enchants))
	}
	if enchants[0].Level() != 1 {
		t.Fatalf("expected clamped enchantment level 1, got %d", enchants[0].Level())
	}
}
