package chunk

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SubChunkHeightMap describes where the surface sits within the sub-chunk of c at index,
// which is what the client lights that sub-chunk from. The returned type is one of the
// protocol.HeightMapData constants; the heights are nil unless it is HeightMapDataHasData,
// because a surface lying wholly above or below the sub-chunk needs no per-column data.
func SubChunkHeightMap(c *Chunk, index int16) (byte, []int8) {
	heights := c.HeightMap()
	heightMap := make([]int8, 256)
	above, below := true, true
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			y := heights.At(x, z)
			at, i := c.SubIndex(y), (uint16(z)<<4)|uint16(x)
			switch {
			case at > index:
				heightMap[i], below = 16, false
			case at < index:
				heightMap[i], above = -1, false
			default:
				heightMap[i], below, above = int8(y-c.SubY(at)), false, false
			}
		}
	}
	switch {
	case above:
		return protocol.HeightMapDataTooHigh, nil
	case below:
		return protocol.HeightMapDataTooLow, nil
	}
	return protocol.HeightMapDataHasData, heightMap
}
