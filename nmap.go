package trie

type Nmap struct {
	buckets [10]*TrieNode
}

func NewNmap() Container {
	return &Nmap{}
}

func (m *Nmap) Set(k byte, v *TrieNode) {
	k = k - '0'
	for i := k - 1; i < 10; i-- {
		if m.buckets[i] != nil {
			m.buckets[i].next = v
			v.prev = m.buckets[i]
			break
		}
	}
	for i := k + 1; i < 10; i++ {
		if m.buckets[i] != nil {
			m.buckets[i].prev = v
			v.next = m.buckets[i]
			break
		}
	}
	m.buckets[k] = v
}
func (m *Nmap) Get(k byte) (*TrieNode, bool) {
	v := m.buckets[k-'0']
	return v, v != nil
}

func (m *Nmap) Del(k byte) bool {
	k = k - '0'
	v := m.buckets[k]
	defer func() {
		v.Free()
		v = nil
	}()
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

func (m *Nmap) Prev(k byte) *TrieNode {
	k = k - '0'
	for i := k - 1; i < k; i-- {
		if m.buckets[i] != nil {
			return m.buckets[i]
		}
	}
	return nil
}
func (m *Nmap) Next(k byte) *TrieNode {
	k = k - '0'
	for i := k + 1; i < 10; i++ {
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

func (m *Nmap) Pad() byte {
	return '0'
}
