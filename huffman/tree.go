package huffman

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"
)

type TreeNode struct {
	IsLeaf bool
	Freq   int64
	Value  rune
	Left   *TreeNode
	Right  *TreeNode
	Code   []bool
}

type RuneInfo struct {
	Freq  int64
	Value rune
	Code  []bool
}

func (t TreeNode) getCodes(dict map[rune]*RuneInfo) map[rune]*RuneInfo {
	if t.IsLeaf {
		dict[t.Value] = &RuneInfo{
			Freq:  t.Freq,
			Value: t.Value,
			Code:  t.Code,
		}
		return dict
	}

	if t.Left != nil {
		dict = t.Left.getCodes(dict)
	}

	if t.Right != nil {
		dict = t.Right.getCodes(dict)
	}

	return dict
}

func (t TreeNode) GetCodes() map[rune]*RuneInfo {
	return t.getCodes(map[rune]*RuneInfo{})
}

func (t TreeNode) Print() {
	if t.IsLeaf {
		fmt.Printf("[leaf](%c: %d), %+v\n", t.Value, t.Freq, t.Code)
	} else {
		fmt.Printf("[not-leaf](%d), %+v\n", t.Freq, t.Code)
	}

	if t.Left != nil {
		t.Left.Print()
	}

	if t.Right != nil {
		t.Right.Print()
	}
}

func buildCodes(t *TreeNode) {
	if t.Right != nil {
		t.Right.Code = append(append([]bool{}, t.Code[:]...), true)
		buildCodes(t.Right)
	}

	if t.Left != nil {
		t.Left.Code = append(append([]bool{}, t.Code[:]...), false)
		buildCodes(t.Left)
	}
}

func buildHuffmanTree(t []*TreeNode) *TreeNode {
	if len(t) < 1 {
		return nil
	}

	for len(t) > 1 {

		sort.Slice(t, func(i, j int) bool {
			return t[i].Freq < t[j].Freq
		})

		leftNode := t[0]
		rightNode := t[1]

		t = t[2:]

		t = append(t, &TreeNode{
			IsLeaf: false,
			Freq:   leftNode.Freq + rightNode.Freq,
			Left:   leftNode,
			Right:  rightNode,
		})
	}

	buildCodes(t[0])

	return t[0]
}

func SerializeTree(node *TreeNode) string {
	if node == nil {
		return "#"
	}

	if node.IsLeaf {
		return fmt.Sprintf("1%c", node.Value)
	}

	return fmt.Sprintf("0%s%s", SerializeTree(node.Left), SerializeTree(node.Right))
}

func CompressedSerialization(node *TreeNode, buf *bytes.Buffer) {
	if node == nil {
		return
	}
	if node.IsLeaf {
		buf.WriteByte(0x80)
		buf.Write([]byte(string(node.Value)))
	} else {
		buf.WriteByte(0x00)
	}

	CompressedSerialization(node.Left, buf)
	CompressedSerialization(node.Right, buf)
}

func DeserializeCompressed(data []byte, pos *int) *TreeNode {
	if *pos >= len(data) {
		return nil
	}

	bit := data[*pos] & 0x80
	*pos++

	if bit == 0x80 {
		// leaf node
		r, size := utf8.DecodeRune(data[*pos:])
		*pos += size
		return &TreeNode{
			IsLeaf: true,
			Value:  r,
		}
	}

	left := DeserializeCompressed(data, pos)
	right := DeserializeCompressed(data, pos)
	return &TreeNode{
		IsLeaf: false,
		Left:   left,
		Right:  right,
	}
}

func DeserializeTree(serialized string) *TreeNode {
	queue := strings.Split(serialized, "")
	var helper func() *TreeNode
	helper = func() *TreeNode {
		if len(queue) == 0 {
			return nil
		}

		val := queue[0]
		queue = queue[1:]

		if val == "1" {
			char := queue[0]
			queue = queue[1:]

			runeValue, _ := utf8.DecodeRuneInString(char)
			return &TreeNode{
				IsLeaf: true,
				Value:  runeValue,
			}
		} else if val == "#" {
			return nil
		} else {
			left := helper()
			right := helper()
			return &TreeNode{
				IsLeaf: false,
				Left:   left,
				Right:  right,
			}
		}
	}

	root := helper()

	buildCodes(root)

	return root
}

func convertFreqToTreeNodes(freq map[rune]int64) []*TreeNode {
	var currentTree []*TreeNode
	for k, v := range freq {
		currentTree = append(currentTree, &TreeNode{
			IsLeaf: true,
			Freq:   v,
			Value:  k,
		})
	}
	return currentTree
}

func BuildHuffmanCodes(freq map[rune]int64) *TreeNode {
	return buildHuffmanTree(convertFreqToTreeNodes(freq))
}
