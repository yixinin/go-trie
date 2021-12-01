package nodetrie

import "strconv"

type NodeTrie struct {
	keySize int
	root    *TrieNode
	head    *TrieNode
	tail    *TrieNode
}

func NewNodeTrie(keySize int) *NodeTrie {
	return &NodeTrie{
		root:    NewTrieNode(0, nil, nil),
		keySize: keySize,
	}
}

type TrieNode struct {
	Prev     *TrieNode
	Next     *TrieNode
	k        byte
	keys     []byte
	val      interface{}
	children *Nmap
}

func NewTrieNode(k byte, key []byte, val interface{}) *TrieNode {
	return &TrieNode{
		children: NewNmap(),
		k:        k,
		keys:     key,
		val:      val,
	}
}

func (t *NodeTrie) Set(key []byte, v interface{}) {
	if len(key) != t.keySize {
		panic("key size should be " + strconv.Itoa(t.keySize))
	}
	cur := t.root
	for i, k := range key {
		if _, ok := cur.children.Get(k); !ok {
			var node *TrieNode
			if i == t.keySize-1 {
				node = NewTrieNode(k, key, v)
			} else {
				node = NewTrieNode(k, nil, nil)
			}

			cur.children.Set(k, node)
			if _, tail := cur.children.Tail(); tail == node {
				if next := cur.Next; next != nil {
					if _, nextHead := next.children.Head(); nextHead != nil {
						nextHead.Prev = node
						node.Next = nextHead
					}
				}
			}
			if _, head := cur.children.Head(); head == node {
				if prev := cur.Prev; prev != nil {
					if _, prevTail := prev.children.Tail(); prevTail != nil {
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
func (t *NodeTrie) Get(key []byte) (interface{}, bool) {
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
func (t *NodeTrie) Gt(key []byte, e bool) interface{} {
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

func (t *NodeTrie) Lt(key []byte, e bool) {

}

func (t *NodeTrie) Foreach(f func(k []byte)) {
	for cur := t.head; cur != nil; cur = cur.Next {
		f(cur.keys)
	}
}
