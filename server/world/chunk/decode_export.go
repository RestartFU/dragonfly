package chunk

import "bytes"

// DecodeSubChunk decodes a single sub-chunk from buf into c at *index. It is the
// counterpart to EncodeSubChunk, for callers streaming sub-chunks one at a time.
func DecodeSubChunk(buf *bytes.Buffer, c *Chunk, index *byte, e Encoding) (*SubChunk, error) {
	return decodeSubChunk(buf, c, index, e)
}
