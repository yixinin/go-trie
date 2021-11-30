package trie

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{
		root: NewTrieNode(),
	}
}

func (t *Trie) Insert(key []byte, val ...string) {
	var cur = t.root
	for _, b := range key {
		if cur.children[b] == nil {
			node := NewTrieNode()
			node.b = b
			cur.children[b] = node
		}
		cur = cur.children[b]
	}
	cur.vals = append(cur.vals, val...)
}
func (t *Trie) Delete(key []byte) {

}

func (t *Trie) FindNext(key []byte) []byte {
	cur := t.root

	var stack = NewStack()
	for _, k := range key {
		if node := cur.children[k]; node != nil {
			stack.Push(node)
			cur = node
			if cur.IsEnd() {
				return stackToKey(stack)
			}
		} else {
			break
		}
	}
	if len(key) == stack.Len() {
		return key
	}

	for cur != nil {
		var b = key[stack.Len()]
		for kb := b + 1; kb > b; kb++ {
			if node := cur.children[kb]; node != nil {
				stack.Push(node)
				cur = node
				if cur.IsEnd() {
					return stackToKey(stack)
				}

				for i := 0; i <= 255; i++ {
					if node := cur.children[byte(i)]; node != nil {
						stack.Push(node)
						cur = node
						if cur.IsEnd() {
							return stackToKey(stack)
						}
						i = 0
					}
				}
			}
		}
		e := stack.Top()
		if e != nil {
			stack.Pop()
			if stack.Top() != nil {
				cur = stack.Pop()
			} else {
				cur = t.root
				// return stackToKey(stack)
			}
		} else {
			return stackToKey(stack)
		}
	}
	return stackToKey(stack)
}
func (t *Trie) FindPrev(key []byte) []byte {
	cur := t.root

	var stack = NewStack()
	for _, k := range key {
		if node := cur.children[k]; node != nil {
			stack.Push(node)
			cur = node
			if cur.IsEnd() {
				return stackToKey(stack)
			}
		} else {
			break
		}
	}
	if len(key) == stack.Len() {
		return key
	}

	for cur != nil {
		var b = key[stack.Len()]
		for kidx := b - 1; kidx < b; kidx-- {
			if node := cur.children[kidx]; node != nil {
				stack.Push(node)
				cur = node
				if cur.IsEnd() {
					return stackToKey(stack)
				}

				for i := 255; i >= 0; i-- {
					if node := cur.children[byte(i)]; node != nil {
						stack.Push(node)
						cur = node
						if cur.IsEnd() {
							return stackToKey(stack)
						}
						i = 255
					}
				}
			}
		}
		e := stack.Top()
		if e != nil {
			stack.Pop()
			if stack.Top() != nil {
				cur = stack.Top()
			} else {
				return stackToKey(stack)
			}
		} else {
			return stackToKey(stack)
		}
	}
	return stackToKey(stack)
}

func (t *Trie) Foreach(f func(key []byte), revert ...bool) {
	if len(revert) > 0 && revert[0] {
		return
	}
	t.foreach(f)
}

func (t *Trie) foreach(f func(key []byte)) {
	cur := t.root
	start := 0
	var stack = NewStack()
	for cur != nil {
		has := false
		for i := start; i < 256; i++ {
			node := cur.children[byte(i)]
			if node != nil {
				stack.Push(node)
				if node.IsEnd() {
					f(stackToKey(stack))
					back := stack.Top()
					if back != nil {
						stack.Pop()
					}
					back = stack.Top()
					if back == nil {
						cur = nil
						break
					}
					cur = stack.Top()
					start = int(node.b) + 1
					has = true
					break
				}

				cur = node
				start = 0
				has = true
				break
			}
		}
		if !has {
			back := stack.Top()
			if back == nil {
				cur = nil
				break
			}
			cur = stack.Top()
			start = int(cur.b) + 1
		}
	}

}

func stackToKey(stack *Stack) []byte {
	var nodes = stack.ToList()
	var pk = make([]byte, 0, stack.Len())
	for _, v := range nodes {
		pk = append(pk, v.b)
		if len(v.vals) > 0 {
			return pk
		}
	}
	return make([]byte, 0)
}
