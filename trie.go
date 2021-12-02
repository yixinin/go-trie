package trie

import (
	"strconv"
)

type Trie struct {
	keySize   int
	root      *TrieNode
	head      *TrieNode
	tail      *TrieNode
	container func() NodeContainer
}

func NewTrie(keySize int, container func() NodeContainer) *Trie {
	return &Trie{
		root:      NewTrieNode(0, nil, nil, container),
		keySize:   keySize,
		container: container,
	}
}

type TrieNode struct {
	Prev     *TrieNode
	Next     *TrieNode
	k        byte
	keys     []byte
	val      interface{}
	children NodeContainer
}

func NewTrieNode(k byte, key []byte, val interface{}, nodeContainer func() NodeContainer) *TrieNode {
	if nodeContainer == nil {
		panic("nodeContainer is nil")
	}
	return &TrieNode{
		children: nodeContainer(),
		k:        k,
		keys:     key,
		val:      val,
	}
}

func (t *Trie) Set(key []byte, v interface{}) {
	if len(key) != t.keySize {
		panic("key size should be " + strconv.Itoa(t.keySize))
	}
	cur := t.root
	for i, k := range key {
		if _, ok := cur.children.Get(k); !ok {
			var node *TrieNode
			if i == t.keySize-1 {
				node = NewTrieNode(k, key, v, t.container)
			} else {
				node = NewTrieNode(k, nil, nil, t.container)
			}

			cur.children.Set(k, node)
			if tail := cur.children.Tail(); tail == node {
				if next := cur.Next; next != nil {
					if nextHead := next.children.Head(); nextHead != nil {
						nextHead.Prev = node
						node.Next = nextHead
					}
				}
			}
			if head := cur.children.Head(); head == node {
				if prev := cur.Prev; prev != nil {
					if prevTail := prev.children.Tail(); prevTail != nil {
						prevTail.Next = node
						node.Prev = prevTail
					}
				}
			}
		}
		cur, _ = cur.children.Get(k)
	}

	if t.head == nil {
		t.head = cur
		t.tail = cur
		return
	}

	if t.head.Prev == cur {
		t.head = cur
	}
	if t.tail.Next == cur {
		t.tail = cur
	}
}
func (t *Trie) Get(key []byte) (interface{}, bool) {
	cur := t.root
	var ok bool
	for _, k := range key {
		cur, ok = cur.children.Get(k)
		if !ok {
			return nil, false
		}
		if len(cur.keys) > 0 {
			return cur.val, true
		}
	}
	return nil, false
}
func (t *Trie) Gt(key []byte, e bool) interface{} {
	cur := t.root
	var ok bool
	for _, k := range key {
		cur, ok = cur.children.Get(k)
		if !ok {
			return nil
		}
		if e && len(cur.keys) > 0 {
			return cur.val
		}
	}
	return nil
}

func (t *Trie) Lt(key []byte, e bool) {

}

func (t *Trie) Foreach(f func(k []byte)) {
	for cur := t.head; cur != nil; cur = cur.Next {
		f(cur.keys)
	}
}
