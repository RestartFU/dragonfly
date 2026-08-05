package chunk_test

import (
	"strings"
	"testing"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/chunk"
)

func TestEncodeLevelChunkPayloadWithBlockEntitiesRoundTrips(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	t.Parallel()

	ch := chunk.New(world.DefaultBlockRegistry, world.Overworld.Range())
	data := chunk.Encode(ch, chunk.NetworkEncoding)
	pos := cube.Pos{3, 64, 5}

	raw, err := chunk.EncodeLevelChunkPayload(data, []chunk.BlockEntity{{
		Pos:  pos,
		Data: map[string]any{"id": "Chest", "pairx": int32(4), "pairz": int32(5)},
	}})
	if err != nil {
		t.Fatalf("EncodeLevelChunkPayload() error = %v", err)
	}

	_, blockEntities, err := chunk.NetworkDecodeWithBlockEntities(world.DefaultBlockRegistry, raw, len(data.SubChunks), world.Overworld.Range())
	if err != nil {
		t.Fatalf("NetworkDecodeWithBlockEntities() error = %v", err)
	}
	if len(blockEntities) != 1 {
		t.Fatalf("decoded %d block entities, want 1", len(blockEntities))
	}
	got := blockEntities[0]
	if got.Pos != pos {
		t.Fatalf("decoded block entity pos = %v, want %v", got.Pos, pos)
	}
	if got.Data["id"] != "Chest" || got.Data["pairx"] != int32(4) || got.Data["pairz"] != int32(5) {
		t.Fatalf("decoded block entity data = %#v", got.Data)
	}
}

func TestEncodeLevelChunkPayloadWithBlockEntitiesReturnsEncodeError(t *testing.T) {
	world.DefaultBlockRegistry.Finalize()
	t.Parallel()

	ch := chunk.New(world.DefaultBlockRegistry, world.Overworld.Range())
	data := chunk.Encode(ch, chunk.NetworkEncoding)

	raw, err := chunk.EncodeLevelChunkPayload(data, []chunk.BlockEntity{{
		Pos:  cube.Pos{3, 64, 5},
		Data: map[string]any{"id": "Chest", "bad": func() {}},
	}})
	if err == nil {
		t.Fatal("expected encode error, got nil")
	}
	if raw != nil {
		t.Fatalf("expected nil payload on encode error, got %d bytes", len(raw))
	}
	if !strings.Contains(err.Error(), "encode block entity") {
		t.Fatalf("error = %q, want block entity context", err)
	}
}

// The client only draws a sub-chunk's block actors when their compounds trail the payload
// carrying their position.
func TestEncodeBlockEntitiesCarriesPositions(t *testing.T) {
	t.Parallel()

	raw, err := chunk.EncodeBlockEntities([]chunk.BlockEntity{
		{Pos: cube.Pos{3, 64, 5}, Data: map[string]any{"id": "Chest"}},
		{Pos: cube.Pos{1, 2, 3}, Data: nil},
	})
	if err != nil {
		t.Fatalf("EncodeBlockEntities() error = %v", err)
	}
	if len(raw) == 0 {
		t.Fatal("encoded nothing")
	}
	if !strings.Contains(string(raw), "Chest") {
		t.Error("encoded data does not carry the block entity id")
	}
	for _, axis := range []string{"x", "y", "z"} {
		if !strings.Contains(string(raw), axis) {
			t.Errorf("encoded data is missing the %q position key", axis)
		}
	}
}
