package trie

type LinkMap struct {
	buckets map[byte]*TrieNode
	head    *TrieNode
	tail    *TrieNode
}

func NewLinkmap() NodeContainer {
	return &LinkMap{
		buckets: make(map[byte]*TrieNode, 1),
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
	if k < m.head.k {
		head := m.head
		v.Next = head
		head.Prev = v
		m.head = v
		return
	}
	if k > m.tail.k {
		tail := m.head
		v.Prev = tail
		tail.Next = v
		m.tail = v
		return
	}
	left := m.head
	right := m.tail
	for left.k < right.k {
		if left.k < k && left.Next.k > k {
			next := left.Next

			left.Next = v
			v.Prev = left

			v.Next = next
			next.Prev = v
			return
		}
		if right.k > k && right.Prev.k < k {
			prev := right.Prev

			right.Prev = v
			v.Next = right

			v.Prev = prev
			prev.Next = v
			return
		}
		left = left.Next
		right = right.Prev
	}
	panic(k)
}

func (m *LinkMap) Get(k byte) (*TrieNode, bool) {
	v, ok := m.buckets[k]
	return v, ok
}

func (m *LinkMap) Prev(k uint8) *TrieNode {
	v, ok := m.buckets[k]
	if ok {
		return v.Prev
	}
	return nil
}
func (m *LinkMap) Next(k byte) *TrieNode {
	v, ok := m.buckets[k]
	if ok {
		return v.Next
	}
	return nil
}
func (m *LinkMap) Head() *TrieNode {
	if m.tail != nil {
		return m.head
	}
	return m.head
}
func (m *LinkMap) Tail() *TrieNode {
	if m.tail != nil {
		return m.tail
	}
	return m.tail
}
func (m *LinkMap) Keys() []byte {
	var keys = make([]byte, 0, len(m.buckets))
	for cur := m.head; cur != nil; cur = cur.Next {
		keys = append(keys, cur.k)
	}
	return keys
}
