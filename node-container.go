package trie

type NodeContainer interface {
	Set(k byte, v *TrieNode)
	Get(k byte) (*TrieNode, bool)
	Prev(k uint8) *TrieNode
	Next(k byte) *TrieNode
	Head() *TrieNode
	Tail() *TrieNode
	Keys() []byte
}
