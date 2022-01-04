package trie

type FsLinkMap struct {
	buckets map[byte]*FsTrieNode
	head    *FsTrieNode
	tail    *FsTrieNode
}

func NewFsLinkMap() FsContainer {
	return &FsLinkMap{
		buckets: make(map[byte]*FsTrieNode, 2),
	}
}

func (m *FsLinkMap) Set(k byte, v *FsTrieNode) {
	_, ok := m.buckets[k]
	if ok {
		return
	}
	m.buckets[k] = v
	if m.head == nil {
		m.head = v
		m.tail = v
		return
	}
	if k < m.head.nodeKey {
		head := m.head
		v.MemNode.next = head
		head.MemNode.prev = v
		m.head = v
		return
	}
	if k > m.tail.nodeKey {
		tail := m.tail
		v.MemNode.prev = tail
		tail.MemNode.next = v
		m.tail = v
		return
	}
	left := m.head
	right := m.tail
	for left.nodeKey < right.nodeKey {
		if left.nodeKey < k && left.MemNode.next.nodeKey > k {
			next := left.MemNode.next

			left.MemNode.next = v
			v.MemNode.prev = left

			v.MemNode.next = next
			next.MemNode.prev = v
			return
		}
		if right.nodeKey > k && right.MemNode.prev.nodeKey < k {
			prev := right.MemNode.prev

			right.MemNode.prev = v
			v.MemNode.next = right

			v.MemNode.prev = prev
			prev.MemNode.next = v
			return
		}
		left = left.MemNode.next
		right = right.MemNode.prev
	}
	panic(k)
}

func (m *FsLinkMap) Get(k byte) (*FsTrieNode, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.buckets[k]
	return v, ok
}

func (m *FsLinkMap) Del(k byte) bool {
	if m == nil {
		return false
	}
	v, ok := m.buckets[k]
	if ok {
		prev := v.MemNode.prev
		next := v.MemNode.next
		if prev != nil {
			prev.MemNode.next = next
		}
		if next != nil {
			next.MemNode.prev = prev
		}
		delete(m.buckets, k)
		return true
	}
	return false
}

func (m *FsLinkMap) Prev(k uint8) *FsTrieNode {
	if m == nil {
		return nil
	}
	v, ok := m.buckets[k]
	if ok {
		return v.MemNode.prev
	}
	return nil
}
func (m *FsLinkMap) Next(k byte) *FsTrieNode {
	if m == nil {
		return nil
	}
	v, ok := m.buckets[k]
	if ok {
		return v.MemNode.next
	}
	return nil
}
func (m *FsLinkMap) Head() *FsTrieNode {
	if m == nil {
		return nil
	}
	if m.tail != nil {
		return m.head
	}
	return m.head
}
func (m *FsLinkMap) Tail() *FsTrieNode {
	if m == nil {
		return nil
	}
	if m.tail != nil {
		return m.tail
	}
	return m.tail
}
func (m *FsLinkMap) Keys() []byte {
	if m == nil {
		return nil
	}
	var keys = make([]byte, 0, len(m.buckets))
	for cur := m.head; cur != nil; cur = cur.MemNode.next {
		keys = append(keys, cur.nodeKey)
	}
	return keys
}

func (m *FsLinkMap) Pad() byte {
	return 0
}
