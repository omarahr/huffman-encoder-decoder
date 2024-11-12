package huffman

import (
	"testing"
)

func TestBuildHuffmanCodes(t *testing.T) {
	input := map[rune]int64{
		'C': 32,
		'D': 42,
		'E': 120,
		'K': 7,
		'L': 42,
		'M': 24,
		'U': 37,
		'Z': 2,
	}
	BuildHuffmanCodes(input)
}
