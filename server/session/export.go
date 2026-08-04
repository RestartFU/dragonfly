package session

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SkinToProtocol converts a skin.Skin to its protocol representation.
func SkinToProtocol(s skin.Skin) protocol.Skin {
	return skinToProtocol(s)
}

// StackFromItem converts an item.Stack to its protocol representation, resolving any block
// it carries through br rather than the global registry.
func StackFromItem(br world.BlockRegistry, it item.Stack) protocol.ItemStack {
	return stackFromItem(br, it)
}
