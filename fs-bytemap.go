package trie

type FsByteMap struct {
	buckets [256]*FsTrieNode
}

func NewFsByteMap() FsContainer {
	return &FsByteMap{}
}

func (m *FsByteMap) Set(k byte, v *FsTrieNode) {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].MemNode.next = v
			v.MemNode.prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].MemNode.prev = v
			v.MemNode.next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *FsByteMap) Get(k byte) (*FsTrieNode, bool) {
	v := m.buckets[k]
	return v, v != nil
}
func (m *FsByteMap) Del(k byte) bool {
	v := m.buckets[k]
	if v != nil {
		prev := v.MemNode.prev
		next := v.MemNode.next
		if prev != nil {
			prev.MemNode.next = next
		}
		if next != nil {
			next.MemNode.prev = prev
		}
		m.buckets[k] = nil
		return true
	}
	return false
}

func (m *FsByteMap) Prev(k byte) *FsTrieNode {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *FsByteMap) Next(k byte) *FsTrieNode {
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *FsByteMap) Head() *FsTrieNode {
	for i := 0; i < 256; i++ {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *FsByteMap) Tail() *FsTrieNode {
	for i := 255; i >= 0; i-- {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *FsByteMap) Keys() []byte {
	var keys = make([]byte, 0, 256)
	for i, v := range m.buckets {
		if v != nil {
			keys = append(keys, byte(i))
		}
	}
	return keys
}

func (m *FsByteMap) Pad() byte {
	return 0
}
