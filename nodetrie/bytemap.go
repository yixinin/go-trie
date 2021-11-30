package nodetrie

type ByteMap struct {
	buckets [256]*TrieNode
}

func NewByteMap() *ByteMap {
	return &ByteMap{}
}

func (m *ByteMap) Set(k byte, v *TrieNode) {
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].Next = v
			v.Prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i > k; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].Prev = v
			v.Next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *ByteMap) Get(k byte) (*TrieNode, bool) {
	v := m.buckets[k]
	return v, v != nil
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
func (m *ByteMap) Head() (byte, *TrieNode) {
	for i := 0; i < 256; i++ {
		if v := m.buckets[byte(i)]; v != nil {
			return byte(i), v
		}
	}
	return 0, nil
}
func (m *ByteMap) Tail() (byte, *TrieNode) {
	for i := 255; i >= 0; i-- {
		if v := m.buckets[byte(i)]; v != nil {
			return byte(i), v
		}
	}
	return 0, nil
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
