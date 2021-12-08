package trie

type LinkMap struct {
	buckets map[byte]*TrieNode
	head    *TrieNode
	tail    *TrieNode
}

func NewLinkmap(init bool) Container {
	if !init {
		return nil
	}
	return &LinkMap{
		buckets: make(map[byte]*TrieNode, 2),
	}
}

func (m *LinkMap) Set(k byte, v *TrieNode) {
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
		v.next = head
		head.prev = v
		m.head = v
		return
	}
	if k > m.tail.nodeKey {
		tail := m.tail
		v.prev = tail
		tail.next = v
		m.tail = v
		return
	}
	left := m.head
	right := m.tail
	for left.nodeKey < right.nodeKey {
		if left.nodeKey < k && left.next.nodeKey > k {
			next := left.next

			left.next = v
			v.prev = left

			v.next = next
			next.prev = v
			return
		}
		if right.nodeKey > k && right.prev.nodeKey < k {
			prev := right.prev

			right.prev = v
			v.next = right

			v.prev = prev
			prev.next = v
			return
		}
		left = left.next
		right = right.prev
	}
	panic(k)
}

func (m *LinkMap) Get(k byte) (*TrieNode, bool) {
	if m == nil {
		return nil, false
	}
	v, ok := m.buckets[k]
	return v, ok
}

func (m *LinkMap) Del(k byte) bool {
	if m == nil {
		return false
	}
	v, ok := m.buckets[k]
	if ok {
		prev := v.prev
		next := v.next
		if prev != nil {
			prev.next = next
		}
		if next != nil {
			next.prev = prev
		}
		delete(m.buckets, k)
		return true
	}
	return false
}

func (m *LinkMap) Prev(k uint8) *TrieNode {
	if m == nil {
		return nil
	}
	v, ok := m.buckets[k]
	if ok {
		return v.prev
	}
	return nil
}
func (m *LinkMap) Next(k byte) *TrieNode {
	if m == nil {
		return nil
	}
	v, ok := m.buckets[k]
	if ok {
		return v.next
	}
	return nil
}
func (m *LinkMap) Head() *TrieNode {
	if m == nil {
		return nil
	}
	if m.tail != nil {
		return m.head
	}
	return m.head
}
func (m *LinkMap) Tail() *TrieNode {
	if m == nil {
		return nil
	}
	if m.tail != nil {
		return m.tail
	}
	return m.tail
}
func (m *LinkMap) Keys() []byte {
	if m == nil {
		return nil
	}
	var keys = make([]byte, 0, len(m.buckets))
	for cur := m.head; cur != nil; cur = cur.next {
		keys = append(keys, cur.nodeKey)
	}
	return keys
}

func (m *LinkMap) Pad() byte {
	return 0
}
