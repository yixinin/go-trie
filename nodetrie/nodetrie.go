package nodetrie

type NodeTrie struct {
	root *TrieNode
	head *TrieNode
	tail *TrieNode
}

func NewNodeTrie() *NodeTrie {
	return &NodeTrie{
		root: NewTrieNode(0, nil),
	}
}

type TrieNode struct {
	Prev     *TrieNode
	Next     *TrieNode
	k        byte
	keys     []byte
	children *ByteMap
}

func NewTrieNode(k byte, key []byte) *TrieNode {
	return &TrieNode{
		children: NewByteMap(),
		k:        k,
		keys:     key,
	}
}

func (t *NodeTrie) Insert(ks []byte) {
	var size = len(ks)
	cur := t.root
	for i, k := range ks {
		if _, ok := cur.children.Get(k); !ok {
			var node *TrieNode
			if i == size-1 {
				node = NewTrieNode(k, ks)
			} else {
				node = NewTrieNode(k, nil)
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

func (t *NodeTrie) Foreach(f func(k []byte)) {
	for cur := t.head; cur != nil; cur = cur.Next {
		f(cur.keys)
	}
}
