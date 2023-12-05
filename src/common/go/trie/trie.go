package trie

type Node struct {
	Children map[string]*Node
	Val      int
	WordEnds bool
}

type ITrie struct {
	root *Node
}

func Trie() *ITrie {
	t := new(ITrie)
	t.root = new(Node)
	return t
}

func (t *ITrie) Insert(word string, val int) {
	current := t.root
	for _, cRune := range word {
		char := string(cRune)

		if current.Children == nil {
			current.Children = make(map[string]*Node)
		}

		if current.Children[char] == nil {
			current.Children[char] = new(Node)
		}
		current = current.Children[char]
	}
	current.WordEnds = true
	current.Val = val
}

func (t *ITrie) Path(partialWord string) *Node {
	current := t.root
	for i, cRune := range partialWord {
		childIdx := string(cRune)

		if i == 0 && current.Children[childIdx] == nil {
			return nil
		}

		if current.Children[childIdx] == nil {
			return current
		}

		current = current.Children[childIdx]
	}

	return current
}
