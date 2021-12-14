package trie

type Container interface {
	Set(k byte, v *TrieNode)
	Get(k byte) (*TrieNode, bool)
	Del(k byte) bool
	Prev(k uint8) *TrieNode
	Next(k byte) *TrieNode
	Head() *TrieNode
	Tail() *TrieNode
	Keys() []byte
	Pad() byte
}

type FsContainer interface {
	Set(k byte, v *FsTrieNode)
	Get(k byte) (*FsTrieNode, bool)
	Del(k byte) bool
	Prev(k uint8) *FsTrieNode
	Next(k byte) *FsTrieNode
	Head() *FsTrieNode
	Tail() *FsTrieNode
	Keys() []byte
	Pad() byte
}
