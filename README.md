
# Huffman Encoder/Decoder

A Go-based implementation of the Huffman encoding and decoding algorithm.


---

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Algorithm Overview](#algorithm-overview)
4. [License](#license)

---

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/omarahr/huffman-encoder-decoder.git
   cd huffman-encoder-decoder
   ```

2. **Install dependencies** (requires Go 1.20+):
   ```bash
   go mod tidy
   ```

3. **Build the project**:
   ```bash
   go build -o huff
   ```

---

## Usage

Run the compiled binary from the command line:

### Encoding
To encode a file:
```bash
./huff input.txt -o output.txt
```

- `-o`: Output file for the encoded data.

### Decoding
To decode a file:
```bash
./huff -d compressed.txt -o decompressed.txt
```

- `-d`: decompress mode.
- `-o`: Output file for the decoded data.

---

## Algorithm Overview

Huffman coding is a lossless data compression algorithm that works by assigning variable-length codes to characters based on their frequencies. Frequently occurring characters are assigned shorter codes, while rarer characters receive longer codes. 

### Steps:
1. **Frequency Analysis**: Calculate the frequency of each character in the input data.
2. **Huffman Tree Construction**:
   - Use a priority queue to build a binary tree.
   - Merge nodes with the smallest frequencies iteratively.
3. **Code Assignment**: Traverse the tree to assign binary codes to each character.
4. **Encoding**: Replace characters in the input with their corresponding binary codes.
5. **Decoding**: Reconstruct the original data using the binary codes and the Huffman tree.



## License

This project is licensed under the [MIT License](LICENSE).
Feel free to use, modify, and distribute it as per the license terms.

---

## Acknowledgments

- Inspired the solution from [Coding Challenges by John Crickett](https://codingchallenges.fyi/challenges/challenge-huffman)
- Explanation from [OpenDSA](https://opendsa-server.cs.vt.edu/ODSA/Books/CS3/html/Huffman.html)
- The Huffman algorithm is a foundational concept in computer science and data compression.
- Inspiration for this project came from [David A. Huffman’s original work](https://en.wikipedia.org/wiki/Huffman_coding).
