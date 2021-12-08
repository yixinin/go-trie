package trie

type ByteMap struct {
	buckets []*TrieNode
}

func NewByteMap(init bool) Container {
	if !init {
		return nil
	}
	return &ByteMap{
		buckets: make([]*TrieNode, 256),
	}
}

func (m *ByteMap) Set(k byte, v *TrieNode) {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].next = v
			v.prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].prev = v
			v.next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *ByteMap) Get(k byte) (*TrieNode, bool) {
	v := m.buckets[k]
	return v, v != nil
}
func (m *ByteMap) Del(k byte) bool {
	v := m.buckets[k]
	if v != nil {
		prev := v.prev
		next := v.next
		if prev != nil {
			prev.next = next
		}
		if next != nil {
			next.prev = prev
		}
		m.buckets[k] = nil
		return true
	}
	return false
}

func (m *ByteMap) Prev(k byte) *TrieNode {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *ByteMap) Next(k byte) *TrieNode {
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *ByteMap) Head() *TrieNode {
	for i := 0; i < 256; i++ {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *ByteMap) Tail() *TrieNode {
	for i := 255; i >= 0; i-- {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *ByteMap) Keys() []byte {
	var keys = make([]byte, 0, 256)
	for i, v := range m.buckets {
		if v != nil {
			keys = append(keys, byte(i))
		}
	}
	return keys
}

func (m *ByteMap) Pad() byte {
	return 0
}
