package chunk

import "bytes"

// DecodeSubChunk decodes a single sub-chunk from buf into c at the index pointed to by
// index, using the Encoding given. It is the decoding counterpart to EncodeSubChunk, and
// exists for callers that stream sub-chunks individually rather than through a full
// NetworkDecode: a proxy receiving SubChunk packets applies them one at a time.
func DecodeSubChunk(buf *bytes.Buffer, c *Chunk, index *byte, e Encoding) (*SubChunk, error) {
	return decodeSubChunk(buf, c, index, e)
}
