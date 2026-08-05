package item

import "testing"

func TestRepeatStacks(t *testing.T) {
	// A sword stacks to 1, so each repetition needs its own stack.
	swords := RepeatStacks([]Stack{NewStack(Sword{Tier: ToolTierDiamond}, 1)}, 3)
	if len(swords) != 3 {
		t.Fatalf("3 swords produced %d stacks, want 3", len(swords))
	}

	// 20 of a 64-stacking item repeated 4 times is 80: one full stack and a remainder.
	dirt := RepeatStacks([]Stack{NewStack(testItem{}, 20)}, 4)
	if len(dirt) != 2 {
		t.Fatalf("80 items produced %d stacks, want 2", len(dirt))
	}
	if got := dirt[0].Count() + dirt[1].Count(); got != 80 {
		t.Fatalf("stacks hold %d items in total, want 80", got)
	}
	if dirt[0].Count() != 64 {
		t.Fatalf("first stack holds %d, want it filled to 64", dirt[0].Count())
	}
}

type testItem struct{}

func (testItem) EncodeItem() (string, int16) { return "minecraft:test", 0 }
