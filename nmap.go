package trie

type Nmap struct {
	buckets [10]*TrieNode
}

func NewNmap() NodeContainer {
	return &Nmap{}
}

func (m *Nmap) Set(k byte, v *TrieNode) {
	for i := k - 1 - '0'; i < 10; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].Next = v
			v.Prev = m.buckets[i]
			break
		}
	}
	for i := k + 1 - '0'; i < 10; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].Prev = v
			v.Next = m.buckets[i]
			break
		}
	}
	m.buckets[k-'0'] = v
}
func (m *Nmap) Get(k byte) (*TrieNode, bool) {
	v := m.buckets[k-'0']
	return v, v != nil
}

func (m *Nmap) Prev(k byte) *TrieNode {
	for i := k - '0' - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *Nmap) Next(k byte) *TrieNode {
	for i := k - '0' + 1; i > k; i++ {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *Nmap) Head() *TrieNode {
	for i := 0; i < 10; i++ {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *Nmap) Tail() *TrieNode {
	for i := 9; i >= 0; i-- {
		if v := m.buckets[byte(i)]; v != nil {
			return v
		}
	}
	return nil
}
func (m *Nmap) Keys() []byte {
	var keys = make([]byte, 0, 10)
	for i, v := range m.buckets {
		if v != nil {
			keys = append(keys, byte(i)+'0')
		}
	}
	return keys
}
