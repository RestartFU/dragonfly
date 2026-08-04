package session

import (
	"testing"

	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/recipe"
)

// A caller decoding recipes off the network can present anything, so the exported form
// reports an unknown pair rather than panicking the way the internal one does.
func TestMatchingStacks_ReportsUnknownItemsInsteadOfPanicking(t *testing.T) {
	type foreign struct{ recipe.Item }

	if _, ok := MatchingStacks(foreign{}, foreign{}); ok {
		t.Fatal("accepted a pair of unknown recipe items")
	}
	if _, ok := MatchingStacks(item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1), foreign{}); ok {
		t.Fatal("accepted an unknown expected item")
	}
	match, ok := MatchingStacks(
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1),
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1),
	)
	if !ok || !match {
		t.Fatalf("identical stacks did not match: match=%v ok=%v", match, ok)
	}
}
