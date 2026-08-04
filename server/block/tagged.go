package block

// Tagged represents a block that carries vanilla block tags. Tags decide which tools the
// client considers able to mine the block and how vanilla behaviour that keys on tags,
// such as tool efficiency, treats it.
//
// Tags are only sent for custom blocks: vanilla blocks carry the tags the client already
// knows for them.
type Tagged interface {
	// Tags returns the vanilla block tags of the block, such as
	// minecraft:is_pickaxe_item_destructible.
	Tags() []string
}
