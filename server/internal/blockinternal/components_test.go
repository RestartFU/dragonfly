package blockinternal

import (
	"testing"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/customblock"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
)

// taggedBlock is a custom block carrying vanilla block tags.
type taggedBlock struct {
	tags []string
}

func (taggedBlock) EncodeBlock() (string, map[string]any) { return "test:tagged", nil }
func (taggedBlock) Model() world.BlockModel               { return nil }
func (taggedBlock) Hash() (uint64, uint64)                { return 0, 0 }
func (b taggedBlock) Tags() []string                      { return b.tags }

func (taggedBlock) Properties() customblock.Properties {
	return customblock.Properties{Cube: true}
}

// untaggedBlock is the same block without tags, so the component map must not gain the
// field at all rather than gaining an empty one.
type untaggedBlock struct{}

func (untaggedBlock) EncodeBlock() (string, map[string]any) { return "test:untagged", nil }
func (untaggedBlock) Model() world.BlockModel               { return nil }
func (untaggedBlock) Hash() (uint64, uint64)                { return 0, 0 }

func (untaggedBlock) Properties() customblock.Properties {
	return customblock.Properties{Cube: true}
}

// breakableBlock is a custom block with the hardness of stone, harvestable only with a
// tool, so its bare-handed time is the hardness times a hundred ticks rather than thirty.
type breakableBlock struct{}

func (breakableBlock) EncodeBlock() (string, map[string]any) { return "test:breakable", nil }
func (breakableBlock) Model() world.BlockModel               { return nil }
func (breakableBlock) Hash() (uint64, uint64)                { return 0, 0 }

func (breakableBlock) Properties() customblock.Properties {
	return customblock.Properties{Cube: true}
}

func (breakableBlock) BreakInfo() block.BreakInfo {
	return block.BreakInfo{
		Hardness:    1.5,
		Harvestable: func(item.Tool) bool { return false },
		Effective:   func(item.Tool) bool { return false },
		Drops:       func(item.Tool, []item.Enchantment) []item.Stack { return nil },
	}
}

// The client reads the component as seconds to destroy, so sending the hardness through
// makes every custom block break several times too fast.
func TestComponents_DestructibleByMiningIsSeconds(t *testing.T) {
	components := Components("test:breakable", breakableBlock{}, 10000)["components"].(map[string]any)

	raw, ok := components["minecraft:destructible_by_mining"]
	if !ok {
		t.Fatal("breakable block has no destructible_by_mining component")
	}
	value := raw.(map[string]any)["value"].(float32)
	// 1.5 hardness, unharvestable by hand: 1.5 * 100 ticks.
	if value != 7.5 {
		t.Fatalf("destructible_by_mining = %v, want 7.5 seconds (hardness is 1.5)", value)
	}
}

func TestComponents_BlockTags(t *testing.T) {
	tags := []string{"minecraft:is_pickaxe_item_destructible", "minecraft:stone"}
	components := Components("test:tagged", taggedBlock{tags: tags}, 10000)

	raw, ok := components["blockTags"]
	if !ok {
		t.Fatal("tagged block has no blockTags")
	}
	// The client decodes the field as a list, so it is encoded as []any rather than a
	// typed slice.
	got, ok := raw.([]any)
	if !ok {
		t.Fatalf("blockTags = %T, want []any", raw)
	}
	if len(got) != len(tags) {
		t.Fatalf("blockTags = %d entries, want %d", len(got), len(tags))
	}
	for i, tag := range tags {
		if got[i] != tag {
			t.Errorf("blockTags[%d] = %v, want %s", i, got[i], tag)
		}
	}
}

func TestComponents_NoBlockTagsWithoutTagged(t *testing.T) {
	components := Components("test:untagged", untaggedBlock{}, 10000)
	if _, ok := components["blockTags"]; ok {
		t.Fatal("a block without tags must not carry a blockTags field")
	}
}

// An empty tag slice must be treated as no tags: sending an empty list would tell the
// client the block has no valid tools rather than leaving it to vanilla defaults.
func TestComponents_EmptyTagsAreOmitted(t *testing.T) {
	components := Components("test:tagged", taggedBlock{}, 10000)
	if _, ok := components["blockTags"]; ok {
		t.Fatal("an empty tag slice must not produce a blockTags field")
	}
}
