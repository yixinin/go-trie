package trie

import (
	"strconv"
)

type Trie struct {
	keySize   int
	root      *TrieNode
	head      *TrieNode
	tail      *TrieNode
	container func() Container
}

func NewTrie(keySize int, container func() Container) *Trie {
	return &Trie{
		root:      newTrieNode(0, nil, nil, container),
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

func newTrieNode(k byte, key []byte, val interface{}, nodeContainer func() Container) *TrieNode {
	if nodeContainer == nil {
		panic("container is nil")
	}
	return &TrieNode{
		children: nodeContainer(),
		nodeKey:  k,
		key:      key,
		val:      val,
	}
}

func (t *Trie) Set(key []byte, v interface{}) {
	if len(key) != t.keySize {
		panic("key size should be " + strconv.Itoa(t.keySize))
	}
	cur := t.root
	for level, nodeKey := range key {
		if _, ok := cur.children.Get(nodeKey); !ok {
			var node *TrieNode
			if level == t.keySize-1 {
				node = newTrieNode(nodeKey, key, v, t.container)
			} else {
				node = newTrieNode(nodeKey, nil, nil, t.container)
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
		}
		cur, _ = cur.children.Get(nodeKey)
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

func (t *Trie) Gt(key []byte) interface{} {
	return t.gt(key, false)
}

func (t *Trie) Gte(key []byte) interface{} {
	return t.gt(key, true)
}
func (t *Trie) gt(key []byte, e bool) interface{} {
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
				return cur.val
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
	return t.lt(key, false)
}
func (t *Trie) Lte(key []byte) interface{} {
	return t.lt(key, true)
}
func (t *Trie) lt(key []byte, e bool) interface{} {
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
				return cur.val
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
