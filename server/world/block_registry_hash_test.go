package world

import "testing"

// The hashes are a wire contract with the client, so they are pinned rather than derived.
func TestNetworkBlockHash_IsStable(t *testing.T) {
	for name, want := range map[string]uint64{
		"minecraft:air":   0xbd584baf00003448,
		"minecraft:stone": 0x2727739c84aa7913,
	} {
		if got := NetworkBlockHash(name); got != want {
			t.Errorf("NetworkBlockHash(%q) = %#016x, want %#016x", name, got, want)
		}
	}
}
