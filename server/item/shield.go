package item

// Shield is a defensive tool used to block attacks while sneaking.
type Shield struct{}

// MaxCount always returns 1.
func (Shield) MaxCount() int {
	return 1
}

// OffHand returns true if a shield may be held in the off hand.
func (Shield) OffHand() bool {
	return true
}

// DurabilityInfo ...
func (Shield) DurabilityInfo() DurabilityInfo {
	return DurabilityInfo{
		MaxDurability: 336,
		BrokenItem:    simpleItem(Stack{}),
	}
}

// RepairableBy ...
func (Shield) RepairableBy(i Stack) bool {
	if planks, ok := i.Item().(interface{ RepairsWoodTools() bool }); ok {
		return planks.RepairsWoodTools()
	}
	return false
}

// EncodeItem ...
func (Shield) EncodeItem() (name string, meta int16) {
	return "minecraft:shield", 0
}
