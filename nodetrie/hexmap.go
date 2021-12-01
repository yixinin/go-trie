package nodetrie

type HexMap struct {
	buckets [16]*TrieNode
}

func ToHex(s string) []uint8 {
	var keys = make([]uint8, len(s))
	for i, v := range []uint8(s) {
		if v >= 'a' {
			keys[i] = v - 'a'
		} else {
			keys[i] = v - '0'
		}
	}
	return keys
}

func NewHexMap() *HexMap {
	return &HexMap{}
}

func (m *HexMap) Set(k uint8, v *TrieNode) {
	for i := k - 1; i < 17; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].Next = v
			v.Prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i < 17; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].Prev = v
			v.Next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *HexMap) Get(k uint8) (*TrieNode, bool) {
	v := m.buckets[k]
	return v, v != nil
}

func (m *HexMap) Prev(k uint8) *TrieNode {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *HexMap) Next(k uint8) *TrieNode {
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *HexMap) Head() (uint8, *TrieNode) {
	for i := 0; i < 10; i++ {
		if v := m.buckets[uint8(i)]; v != nil {
			return uint8(i), v
		}
	}
	return 0, nil
}
func (m *HexMap) Tail() (uint8, *TrieNode) {
	for i := 9; i < 10; i-- {
		if v := m.buckets[uint8(i)]; v != nil {
			return uint8(i), v
		}
	}
	return 0, nil
}
func (m *HexMap) Keys() []uint8 {
	var keys = make([]uint8, 0, 10)
	for i, v := range m.buckets {
		if v != nil {
			keys = append(keys, uint8(i))
		}
	}
	return keys
}
