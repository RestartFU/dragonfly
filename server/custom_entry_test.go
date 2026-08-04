package server

import (
	"image"
	"testing"

	"github.com/df-mc/dragonfly/server/block/customblock"
	"github.com/df-mc/dragonfly/server/item/category"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// customItem is a custom item that is not also a block.
type customItem struct{}

func (customItem) EncodeItem() (string, int16) { return "test:item", 0 }
func (customItem) Name() string                { return "Test Item" }
func (customItem) Texture() image.Image        { return image.NewNRGBA(image.Rect(0, 0, 1, 1)) }
func (customItem) Category() category.Category { return category.Construction() }

// customBlockItem is the item form of a custom block, which the client places rather than
// treating as a component-based item.
type customBlockItem struct{ customItem }

func (customBlockItem) EncodeBlock() (string, map[string]any) { return "test:item", nil }
func (customBlockItem) Model() world.BlockModel               { return nil }
func (customBlockItem) Hash() (uint64, uint64)                { return 0, 0 }

func (customBlockItem) Properties() customblock.Properties {
	return customblock.Properties{Cube: true}
}

func TestCustomItemEntry(t *testing.T) {
	for _, tc := range []struct {
		name           string
		item           world.CustomItem
		wantVersion    int32
		wantComponents bool
	}{
		{"plain item", customItem{}, protocol.ItemEntryVersionDataDriven, true},
		// A custom block's item is placed by the block, so it carries no components.
		{"custom block item", customBlockItem{}, protocol.ItemEntryVersionNone, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			entry := CustomItemEntry(tc.item, 1234)
			if entry.Name != "test:item" {
				t.Errorf("name = %q, want test:item", entry.Name)
			}
			if entry.RuntimeID != 1234 {
				t.Errorf("runtime ID = %d, want the one passed (1234)", entry.RuntimeID)
			}
			if entry.Version != tc.wantVersion {
				t.Errorf("version = %d, want %d", entry.Version, tc.wantVersion)
			}
			if entry.ComponentBased != tc.wantComponents {
				t.Errorf("component based = %v, want %v", entry.ComponentBased, tc.wantComponents)
			}
			components, ok := entry.Data["components"].(map[string]any)
			if !ok {
				t.Fatalf("entry data has no components: %v", entry.Data)
			}
			properties, ok := components["item_properties"].(map[string]any)
			if !ok {
				t.Fatalf("components have no item properties: %v", components)
			}
			if got := properties["creative_category"]; got != int32(category.Construction().Uint8()) {
				t.Errorf("creative category = %v, want the item's own", got)
			}
			icon, ok := properties["minecraft:icon"].(map[string]any)
			if !ok {
				t.Fatalf("item properties have no icon: %v", properties)
			}
			textures, _ := icon["textures"].(map[string]any)
			if got := textures["default"]; got != "test:item" {
				t.Errorf("icon texture = %v, want the item identifier", got)
			}
		})
	}
}
