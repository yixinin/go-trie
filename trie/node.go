package trie

type TrieNode struct {
	children []*TrieNode
	b        byte
	vals     []string
}

func NewTrieNode(val ...string) *TrieNode {
	return &TrieNode{
		children: make([]*TrieNode, 256),
		vals:     val,
	}
}

func (n *TrieNode) IsEnd() bool {
	return len(n.vals) > 0
}
