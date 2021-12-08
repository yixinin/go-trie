package trie

type HexMap struct {
	buckets [16]*TrieNode
}

func NewHexMap(init bool) Container {
	if !init {
		return nil
	}
	return &HexMap{}
}

func toIndex(k uint8) uint8 {
	if k >= 'a' {
		return k - 'a' + 10
	}
	return k - '0'
}
func toChar(k int) uint8 {
	if k > 9 {
		return uint8(k) + 'a' - 10
	}
	return uint8(k) + '0'
}
func (m *HexMap) Set(k uint8, v *TrieNode) {
	k = toIndex(k)
	for i := k - 1; i < 16; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].next = v
			v.prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i < 16; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].prev = v
			v.next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *HexMap) Get(k uint8) (*TrieNode, bool) {
	if m == nil {
		return nil, false
	}
	k = toIndex(k)
	v := m.buckets[k]
	return v, v != nil
}

func (m *HexMap) Del(k byte) bool {
	if m == nil {
		return false
	}
	k = toIndex(k)
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

func (m *HexMap) Prev(k uint8) *TrieNode {
	if m == nil {
		return nil
	}
	k = toIndex(k)
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *HexMap) Next(k uint8) *TrieNode {
	if m == nil {
		return nil
	}
	k = toIndex(k)
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *HexMap) Head() *TrieNode {
	if m == nil {
		return nil
	}
	for i := 0; i < 16; i++ {
		if v := m.buckets[uint8(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *HexMap) Tail() *TrieNode {
	if m == nil {
		return nil
	}
	for i := 15; i >= 0; i-- {
		if v := m.buckets[uint8(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *HexMap) Keys() []uint8 {
	if m == nil {
		return nil
	}
	var keys = make([]uint8, 0, 16)
	for i, v := range m.buckets {
		if v != nil {
			keys = append(keys, toChar(i))
		}
	}
	return keys
}

func (m *HexMap) Pad() byte {
	return '0'
}
