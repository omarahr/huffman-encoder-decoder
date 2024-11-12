package huffman

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitString_GetReadyBytes(t *testing.T) {
	// Arrange
	prefillBoolean := func(s []bool, b bool) []bool {
		for i := range s {
			s[i] = b
		}
		return s
	}

	type args struct {
		bytesToAdd []bool
	}
	tests := []struct {
		name           string
		args           args
		expectedResult []byte
		expectedSize   int
		expectedBytes  []byte
	}{
		{
			name: "empty bit-string, return empty byte array",
			args: args{
				bytesToAdd: nil,
			},
			expectedResult: []byte{},
			expectedSize:   0,
			expectedBytes:  []byte{},
		},
		{
			name: "one bit, no bytes ready, return nil",
			args: args{
				bytesToAdd: []bool{true},
			},
			expectedResult: []byte{},
			expectedSize:   1,
			expectedBytes:  []byte{1 << 7},
		},
		{
			name: "7 bit, no bytes ready, return nil",
			args: args{
				bytesToAdd: prefillBoolean(make([]bool, 7), true),
			},
			expectedResult: []byte{},
			expectedSize:   7,
			expectedBytes:  []byte{0xFE}, // 11111110
		},
		{
			name: "8 bit, 1 byte ready",
			args: args{
				bytesToAdd: prefillBoolean(make([]bool, 8), true),
			},
			expectedResult: []byte{0xFF},
			expectedSize:   0,
			expectedBytes:  []byte{},
		},
		{
			name: "9 bit, 1 byte ready",
			args: args{
				bytesToAdd: prefillBoolean(make([]bool, 9), true),
			},
			expectedResult: []byte{0xFF},
			expectedSize:   1,
			expectedBytes:  []byte{1 << 7},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitString()
			bs.AddBits(tt.args.bytesToAdd)

			// Act
			result := bs.GetReadyBytes()
			size := bs.Size()
			remainingBytes := bs.GetBytes()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedSize, size)
			assert.Equal(t, tt.expectedBytes, remainingBytes)
		})
	}

}

func TestBitString_GetTrailingSize(t *testing.T) {
	// Arrange
	type args struct {
		bytesToAdd []bool
	}

	tests := []struct {
		name           string
		args           args
		expectedResult int
	}{
		{
			name: "empty bit-string, return 0",
			args: args{
				bytesToAdd: nil,
			},
			expectedResult: 0,
		},
		{
			name: "added One Bit, return 1",
			args: args{
				bytesToAdd: make([]bool, 1),
			},
			expectedResult: 1,
		},
		{
			name: "added Two Bits, return 2",
			args: args{
				bytesToAdd: make([]bool, 2),
			},
			expectedResult: 2,
		},
		{
			name: "added Three Bits, return 3",
			args: args{
				bytesToAdd: make([]bool, 3),
			},
			expectedResult: 3,
		},
		{
			name: "added 8 bits, return 8",
			args: args{
				bytesToAdd: make([]bool, 8),
			},
			expectedResult: 8,
		},
		{
			name: "added 9 bits, return 1",
			args: args{
				bytesToAdd: make([]bool, 9),
			},
			expectedResult: 1,
		},
		{
			name: "added 16 bits, return 8",
			args: args{
				bytesToAdd: make([]bool, 16),
			},
			expectedResult: 8,
		},
		{
			name: "added 17 bits, return 1",
			args: args{
				bytesToAdd: make([]bool, 17),
			},
			expectedResult: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitString()
			bs.AddBits(tt.args.bytesToAdd)

			// Act
			result := bs.GetTrailingSize()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestBitString_GetBytes(t *testing.T) {
	// Arrange
	type args struct {
		bytesToAdd []bool
	}

	tests := []struct {
		name           string
		args           args
		expectedResult []byte
		expectedSize   int
	}{
		{
			name: "empty bit-string, return empty byte array",
			args: args{
				bytesToAdd: nil,
			},
			expectedResult: []byte{},
			expectedSize:   0,
		},
		{
			name: "added One (true) Bit, return one byte",
			args: args{
				bytesToAdd: []bool{true},
			},
			expectedResult: []byte{1 << 7}, // 10000000
			expectedSize:   1,
		},
		{
			name: "added One (false) Bit, return one byte",
			args: args{
				bytesToAdd: []bool{false},
			},
			expectedResult: []byte{0}, // 00000000
			expectedSize:   1,
		},
		{
			name: "added Two (true) Bits, return one byte",
			args: args{
				bytesToAdd: []bool{true, true},
			},
			expectedResult: []byte{(1 << 7) | (1 << 6)}, // 11000000
			expectedSize:   2,
		},
		{
			name: "added Two (true, false) Bits, return one byte",
			args: args{
				bytesToAdd: []bool{true, false},
			},
			expectedResult: []byte{1 << 7}, // 10000000
			expectedSize:   2,
		},
		{
			name: "added Two (false, true) Bits, return one byte",
			args: args{
				bytesToAdd: []bool{false, true},
			},
			expectedResult: []byte{1 << 6}, // 01000000
			expectedSize:   2,
		},
		{
			name: "added 8 true Bits, return one byte",
			args: args{
				bytesToAdd: []bool{true, true, true, true, true, true, true, true},
			},
			expectedResult: []byte{0xFF}, // 11111111
			expectedSize:   8,
		},
		{
			name: "added 9 true Bits, return two bytes",
			args: args{
				bytesToAdd: []bool{true, true, true, true, true, true, true, true, true},
			},
			expectedResult: []byte{0xFF, 1 << 7}, // 11111111, 10000000
			expectedSize:   9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := NewBitString()

			bs.AddBits(tt.args.bytesToAdd)

			// Act
			result := bs.GetBytes()
			size := bs.Size()

			// Assert
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedSize, size)
		})
	}
}
