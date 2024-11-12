package encoder

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/omarahr/huffman-encoder-decoder/huffman"
)

const (
	MaxBitsSize = 8 * 16 * 1024
)

func reader(filePath string) (reader *bufio.Reader, file *os.File, closer func() error, err error) {
	file, err = os.Open(filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	// creating reader
	reader = bufio.NewReader(file)

	return reader, file, file.Close, nil
}

func Compress(inputFilePath string, outputFilePath string) {
	bReader, _, closer, err := reader(inputFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	var freqMap map[rune]int64
	if freqMap, err = getFreqMap(bReader); err != nil {
		log.Fatal(err)
		return
	}

	_ = closer()

	root := huffman.BuildHuffmanCodes(freqMap)
	serializedTree := huffman.SerializeTree(root)

	// encoding file
	file, wCloser, err := getOutputWriter(outputFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() { _ = wCloser() }()

	// encoding header
	// header contains
	// first two bytes are the header size
	// after that the huffman tree encoded in a pre-order fashion
	if err = writeHeader(file, serializedTree); err != nil {
		log.Fatal(err)
		return
	}

	// encoding the file content,
	// Add the padding byte at the end of the file as well
	if bReader, _, closer, err = reader(inputFilePath); err != nil {
		log.Fatal(err)
		return
	}
	defer func() { _ = closer() }()

	if err = writeFileCompressedData(file, bReader, root); err != nil {
		log.Fatal(err)
		return
	}
}

// from frequency map to building the encoding
func getFreqMap(reader *bufio.Reader) (map[rune]int64, error) {
	freq := make(map[rune]int64)
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		freq[r]++
	}
	return freq, nil
}

func getOutputWriter(outputFilePath string) (*os.File, func() error, error) {
	file, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file, file.Close, nil
}

func writeHeader(file *os.File, serializedTree string) error {
	headerBytes := []byte(serializedTree)
	headerLength := uint16(len(headerBytes))

	fmt.Printf("header length: %d\n", headerLength)
	// writing header length
	if err := binary.Write(file, binary.LittleEndian, headerLength); err != nil {
		return err
	}

	// writing header
	_, err := file.Write(headerBytes)

	return err
}

func writeFileCompressedData(outputFile *os.File, inputReader *bufio.Reader, root *huffman.TreeNode) error {
	dict := root.GetCodes()      // fetching the code map
	bs := huffman.NewBitString() // creating bit string

	currentSize := 0
	originalSize := 0
	var info *huffman.RuneInfo
	for {
		r, rs, err := inputReader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			return err
		}

		originalSize += rs

		if info = dict[r]; info == nil {
			log.Fatal(fmt.Sprintf("unknown rune: %c", r))
			return err
		}

		bs.AddBits(info.Code)

		if bs.Size() >= MaxBitsSize {
			readyBytes := bs.GetReadyBytes()
			currentSize += len(readyBytes)
			if _, err = outputFile.Write(readyBytes); err != nil {
				log.Fatal(err)
				return err
			}
		}
	}

	// write remaining bytes
	compressedData := bs.GetBytes()
	currentSize += len(compressedData)

	reduction := (1.0 - (float64(currentSize) / float64(originalSize))) * 100.0
	fmt.Printf("original length: %d, compressed data length: %d, reduction: %.2f%%\n",
		originalSize,
		currentSize,
		reduction,
	)

	if _, err := outputFile.Write(compressedData); err != nil {
		log.Fatal(err)
		return err
	}

	// add padding
	paddingByte := bs.GetTrailingSize()
	fmt.Printf("padding byte: %d\n", paddingByte)
	if _, err := outputFile.Write([]byte{byte(paddingByte)}); err != nil {
		log.Fatal(err)
		return err
	}

	// flush
	if err := outputFile.Sync(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
