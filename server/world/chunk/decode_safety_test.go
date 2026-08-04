package chunk

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// A sub-chunk's index width is fixed by the storage, not by how many palette entries were
// sent, so a payload can pack an index the palette does not cover. That has to resolve to
// something rather than panicking deep inside a later read.
func TestPalette_ValueToleratesAnUncoveredIndex(t *testing.T) {
	p := newPalette(4, []uint32{7, 8})
	if got := p.Value(0); got != 7 {
		t.Fatalf("Value(0) = %d, want 7", got)
	}
	if got := p.Value(15); got != 7 {
		t.Fatalf("Value(15) = %d, want the first entry (7) as the fallback", got)
	}
	if got := newPalette(4, nil).Value(3); got != 0 {
		t.Fatalf("Value on an empty palette = %d, want 0", got)
	}
}

// The count is read off the wire and used to size an allocation, so it must be bounded
// before that allocation happens rather than failing later on a short entry read.
func TestDecodePalette_RejectsHostileCounts(t *testing.T) {
	encodings := map[string]interface {
		decodePalette(buf *bytes.Buffer, blockSize paletteSize, pe paletteEncoding) (*Palette, error)
	}{
		"network":           networkEncoding{},
		"networkPersistent": networkPersistentEncoding{},
	}
	for name, e := range encodings {
		t.Run(name+"/beyond the index width", func(t *testing.T) {
			buf := &bytes.Buffer{}
			// 4 bits per block addresses 16 entries; claim more.
			_ = protocol.WriteVarint32(buf, 17)
			_, err := e.decodePalette(buf, paletteSize(4), nil)
			if err == nil {
				t.Fatal("accepted a palette larger than the storage can address")
			}
			if !strings.Contains(err.Error(), "exceeds") {
				t.Fatalf("rejected only after allocating: %v", err)
			}
		})
		t.Run(name+"/truncated count", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("panicked instead of returning an error: %v", r)
				}
			}()
			if _, err := e.decodePalette(&bytes.Buffer{}, paletteSize(4), nil); err == nil {
				t.Fatal("accepted a truncated palette count")
			}
		})
		t.Run(name+"/zero count", func(t *testing.T) {
			buf := &bytes.Buffer{}
			_ = protocol.WriteVarint32(buf, 0)
			if _, err := e.decodePalette(buf, paletteSize(4), nil); err == nil {
				t.Fatal("accepted a zero palette count")
			}
		})
	}
}
