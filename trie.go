package trie

import (
	"log"
	"strconv"
)

type Trie struct {
	keySize   int
	root      *TrieNode
	head      *TrieNode
	tail      *TrieNode
	container func() Container
	size      int
}

func NewTrie(keySize int, container func() Container) *Trie {
	return &Trie{
		root:      newTrieNode(0, container),
		keySize:   keySize,
		container: container,
	}
}

type TrieNode struct {
	prev     *TrieNode
	next     *TrieNode
	nodeKey  byte
	key      []byte
	val      interface{}
	children Container
}

func newTrieNodeLeaf(k byte, key []byte, val interface{}) *TrieNode {
	return &TrieNode{
		nodeKey: k,
		key:     key,
		val:     val,
	}
}

func newTrieNode(k byte, nodeContainer func() Container) *TrieNode {
	if nodeContainer == nil {
		panic("container is nil")
	}
	return &TrieNode{
		children: nodeContainer(),
		nodeKey:  k,
	}
}

func (node *TrieNode) Free() {
	if node == nil {
		return
	}
	node.children = nil
	node.key = nil
	node.next = nil
	node.prev = nil
	node.val = nil
	node = nil
}

func (t *Trie) Set(key []byte, val interface{}) {
	if len(key) != t.keySize {
		panic("key size should be " + strconv.Itoa(t.keySize))
	}

	cur := t.root
	var ok bool
	for level, nodeKey := range key {
		var node *TrieNode
		if cur.children == nil {
			log.Println(level)
		}
		if node, ok = cur.children.Get(nodeKey); !ok {
			if level == t.keySize-1 {
				node = newTrieNodeLeaf(nodeKey, key, val)
				t.size++
			} else {
				node = newTrieNode(nodeKey, t.container)
			}

			cur.children.Set(nodeKey, node)
			if tail := cur.children.Tail(); tail == node {
				if next := cur.next; next != nil {
					if nextHead := next.children.Head(); nextHead != nil {
						nextHead.prev = node
						node.next = nextHead
					}
				}
			}
			if head := cur.children.Head(); head == node {
				if prev := cur.prev; prev != nil {
					if prevTail := prev.children.Tail(); prevTail != nil {
						prevTail.next = node
						node.prev = prevTail
					}
				}
			}
		} else {
			node.val = val
		}
		cur = node
	}

	if t.head == nil {
		t.head = cur
		t.tail = cur
		return
	}

	if t.head.prev == cur {
		t.head = cur
	}
	if t.tail.next == cur {
		t.tail = cur
	}
}
func (t *Trie) Get(key []byte) (interface{}, bool) {
	if len(key) != t.keySize {
		return nil, false
	}
	cur := t.root
	var ok bool
	for _, k := range key {
		cur, ok = cur.children.Get(k)
		if !ok {
			return nil, false
		}
		if len(cur.key) > 0 {
			return cur.val, true
		}
	}
	return nil, false
}

func (t *Trie) Del(key []byte) bool {
	if len(key) != t.keySize {
		return false
	}
	cur := t.root
	var stack = NewStack()
	stack.Push(cur)
	var ok bool
	for _, k := range key {
		cur, ok = cur.children.Get(k)
		if !ok {
			return false
		}

		if len(cur.key) > 0 {
			// delete leaf
			stack.Top().children.Del(k)
			break
		}
		stack.Push(cur)
	}
	// delete parent
	for {
		cur := stack.Pop()
		if cur.children.Head() == nil {
			stack.Top().children.Del(cur.nodeKey)
		} else {
			break
		}
	}
	return false
}

func (t *Trie) Gt(key []byte) interface{} {
	if node := t.gt(key, false); node != nil {
		return node.val
	}
	return nil
}

func (t *Trie) Gte(key []byte) interface{} {
	if node := t.gt(key, true); node != nil {
		return node.val
	}
	return nil
}
func (t *Trie) gt(key []byte, e bool) *TrieNode {
	key = t.PadRight(key)
	cur := t.root
	toNext := false
	toHead := false
	var level = 0
	var stack = NewStack()
	stack.Push(cur)
	for level >= 0 && level < t.keySize && cur != nil {
		nodeKey := key[level]
		if toNext {
			c := cur.children.Next(nodeKey)
			if c == nil {
				level--
				toHead = true
				cur = stack.Pop()
				continue
			}
			cur = c
		} else if toHead {
			cur = cur.children.Head()
		} else {
			c, ok := cur.children.Get(nodeKey)
			if !ok {
				toNext = true
				stack.Pop()
				continue
			} else {
				cur = c
			}
		}

		if len(cur.key) > 0 {
			if e || toHead || toNext {
				return cur
			}
			toNext = true
			stack.Pop()
			continue
		}

		level++
		stack.Push(cur)
		toNext = false
	}

	return nil
}

func (t *Trie) Lt(key []byte) interface{} {
	if node := t.lt(key, false); node != nil {
		return node.val
	}
	return nil
}
func (t *Trie) Lte(key []byte) interface{} {
	if node := t.lt(key, true); node != nil {
		return node.val
	}
	return nil
}
func (t *Trie) lt(key []byte, e bool) *TrieNode {
	key = t.PadRight(key)
	cur := t.root
	toPrev := false
	toTail := false
	var level = 0
	var stack = NewStack()
	stack.Push(cur)
	for level >= 0 && level < t.keySize && cur != nil {
		nodeKey := key[level]
		if toPrev {
			c := cur.children.Prev(nodeKey)
			if c == nil {
				level--
				toTail = true
				cur = stack.Pop()
				continue
			}
			cur = c
		} else if toTail {
			cur = cur.children.Tail()
		} else {
			c, ok := cur.children.Get(nodeKey)
			if !ok {
				toPrev = true
				stack.Pop()
				continue
			} else {
				cur = c
			}
		}

		if len(cur.key) > 0 {
			if e || toTail || toPrev {
				return cur
			}
			toPrev = true
			stack.Pop()
			continue
		}
		level++
		stack.Push(cur)
		toPrev = false
	}

	return nil
}

func (t *Trie) Foreach(f func(key []byte, val interface{})) {
	for cur := t.head; cur != nil; cur = cur.next {
		f(cur.key, cur.val)
	}
}

func (t *Trie) PadRight(key []byte) []byte {
	if size := len(key); size < t.keySize {
		var fullKey = make([]byte, t.keySize)
		pad := t.root.children.Pad()
		for i := size - 1; i < t.keySize; i++ {
			fullKey[i] = pad
		}
		copy(fullKey[:size], key)
		return fullKey
	}
	return key[:t.keySize]
}

func (t *Trie) Len() int {
	return t.size
}
