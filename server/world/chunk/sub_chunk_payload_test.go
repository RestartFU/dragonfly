package chunk_test

import (
	"testing"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// The client lights a sub-chunk from this map, so a surface lying wholly outside the
// sub-chunk has to collapse to TooHigh/TooLow rather than ship 256 sentinel heights.
func TestSubChunkHeightMap(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	t.Parallel()

	c := chunk.New(world.DefaultBlockRegistry, world.Overworld.Range())
	stone, ok := world.DefaultBlockRegistry.StateToRuntimeID("minecraft:stone", map[string]any{"stone_type": "stone"})
	if !ok {
		stone = world.DefaultBlockRegistry.BlockRuntimeID(block.Stone{})
	}
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			c.SetBlock(x, 64, z, 0, stone)
		}
	}
	surface := c.SubIndex(64)

	t.Run("surface above the sub-chunk", func(t *testing.T) {
		mapType, heights := chunk.SubChunkHeightMap(c, surface-1)
		if mapType != protocol.HeightMapDataTooHigh {
			t.Fatalf("type = %d, want TooHigh", mapType)
		}
		if heights != nil {
			t.Errorf("carried %d heights, want none", len(heights))
		}
	})

	t.Run("surface below the sub-chunk", func(t *testing.T) {
		mapType, heights := chunk.SubChunkHeightMap(c, surface+1)
		if mapType != protocol.HeightMapDataTooLow {
			t.Fatalf("type = %d, want TooLow", mapType)
		}
		if heights != nil {
			t.Errorf("carried %d heights, want none", len(heights))
		}
	})

	t.Run("surface inside the sub-chunk", func(t *testing.T) {
		mapType, heights := chunk.SubChunkHeightMap(c, surface)
		if mapType != protocol.HeightMapDataHasData {
			t.Fatalf("type = %d, want HasData", mapType)
		}
		if len(heights) != 256 {
			t.Fatalf("carried %d heights, want 256", len(heights))
		}
		// The height map holds the first free Y above the surface, expressed relative to
		// the sub-chunk's own base rather than to world Y.
		want := int8(65 - c.SubY(surface))
		for i, h := range heights {
			if h != want {
				t.Fatalf("height[%d] = %d, want %d", i, h, want)
			}
		}
	})
}
