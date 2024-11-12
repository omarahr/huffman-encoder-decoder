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

func Decompress(inputFilePath string, outputFilePath string) {
	bReader, inputFile, closer, err := reader(inputFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() { _ = closer() }()

	// 1.reading the header length
	var headerLength uint16
	if err = binary.Read(bReader, binary.LittleEndian, &headerLength); err != nil {
		log.Fatal(err)
		return
	}

	// 2.reading the header
	headerBytes := make([]byte, headerLength)
	if _, err = bReader.Read(headerBytes); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("header length: %d\n", headerLength)

	// 3.build the huffman tree
	root := huffman.DeserializeTree(string(headerBytes))

	writer, wCloser, err := getOutputWriter(outputFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() { _ = wCloser() }()

	// figure out how many bytes remaining to deserialize
	state, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
		return
	}

	remainingBytes := state.Size() - 2 - int64(headerLength)
	fmt.Printf("remaining bytes: %d\n", remainingBytes)

	// 4. decode the compressed file and write the original file
	decodeAndWrite(bReader, writer, root)
}

func decodeAndWrite(reader *bufio.Reader, outputFile *os.File, root *huffman.TreeNode) {
	buffer := make([]byte, 10*1024)
	var possibleFinals []byte

	currentNode := root
	for {
		n, err := reader.Read(buffer)

		if err == io.EOF {
			currentNode, _ = flushLastBytes(possibleFinals, currentNode, root, outputFile)
			break
		}

		if err != nil {
			log.Fatal(err)
			return
		}

		var readerBuffer []byte
		if n == 2 {
			// todo handle this case
			// case if n is two
			// last byte is buffer [0]
			// padding byte is buffer[1]
			// process possible finals
			readerBuffer = append([]byte{}, possibleFinals...)
			possibleFinals = append([]byte{}, buffer[:n]...)
		} else if n == 1 {
			// case if n is one
			// last byte is possible finals[1]
			// padding byte is buffer [0]
			// process possible buffer [1]
			readerBuffer = append([]byte{}, possibleFinals[0])
			possibleFinals = append([]byte{possibleFinals[1]}, buffer[0])
		} else {
			readerBuffer = append(append([]byte{}, possibleFinals...), buffer[:n-2]...)
			possibleFinals = append([]byte{}, buffer[n-2:n]...)
		}

		currentNode, err = writeBytes(readerBuffer, currentNode, root, outputFile)
		if err != nil {
			log.Fatal(err)
			return
		}

	}

	if err := outputFile.Sync(); err != nil {
		log.Fatal(err)
		return
	}
}

func writeBytes(readerBuffer []byte, currentNode, root *huffman.TreeNode, outputFile *os.File) (*huffman.TreeNode, error) {
	for j := 0; j < len(readerBuffer); j++ {
		currentByte := readerBuffer[j]
		for i := 7; i >= 0; i-- {
			currentBit := (currentByte >> i) & 1
			if currentBit == 0 {
				currentNode = currentNode.Left
			} else {
				currentNode = currentNode.Right
			}

			if currentNode.IsLeaf {
				if _, err := outputFile.WriteString(string(currentNode.Value)); err != nil {
					log.Fatal(err)
					return nil, err
				}
				currentNode = root
			}
		}
	}

	return currentNode, nil
}

func flushLastBytes(lastBytes []byte, currentNode, root *huffman.TreeNode, outputFile *os.File) (*huffman.TreeNode, error) {
	lastByte := lastBytes[0]
	paddingByte := lastBytes[1]

	fmt.Printf("lastByte: %d, paddingByte: %d\n", lastByte, paddingByte)

	for i := 7; i >= (8 - int(paddingByte)); i-- {
		currentBit := (lastByte >> i) & 1
		if currentBit == 0 {
			currentNode = currentNode.Left
		} else {
			currentNode = currentNode.Right
		}

		if currentNode.IsLeaf {
			if _, err := outputFile.WriteString(string(currentNode.Value)); err != nil {
				log.Fatal(err)
				return nil, err
			}
			currentNode = root
		}
	}

	return currentNode, nil
}
