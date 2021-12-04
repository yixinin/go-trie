package trie

func SliceEq(s1, s2 []byte) (bool, int) {
	if len(s1) != len(s2) {
		return false, 0
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false, i
		}
	}
	return true, 0
}
