package session

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SkinToProtocol converts a skin.Skin to its protocol representation, including its
// animations, geometry and cape. Callers that build player list or skin packets outside a
// Session need the same conversion the Session performs internally.
func SkinToProtocol(s skin.Skin) protocol.Skin {
	return skinToProtocol(s)
}

// StackFromItem converts an item.Stack to its protocol representation, resolving any block
// the item carries through the world.BlockRegistry given rather than the global one. A
// caller holding a registry built from another server's palette must pass that registry, so
// the block runtime ID it writes means the same thing to the receiving client.
func StackFromItem(br world.BlockRegistry, it item.Stack) protocol.ItemStack {
	return stackFromItem(br, it)
}
