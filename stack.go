package trie

type stackNode struct {
	v    *TrieNode
	next *stackNode
}

type Stack struct {
	head *stackNode
	size int
}

func NewStack() *Stack {
	return &Stack{}
}

func (s *Stack) Pop() *TrieNode {
	if s.head == nil {
		return nil
	}
	s.size--
	v := s.head
	s.head = v.next
	return v.v
}

func (s *Stack) Top() *TrieNode {
	if s.head == nil {
		return nil
	}
	return s.head.v
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Push(v *TrieNode) {
	node := &stackNode{
		v: v,
	}
	s.size++

	next := s.head
	s.head = node

	s.head.next = next
}

func (s *Stack) ToList() []*TrieNode {
	var nodes = make([]*TrieNode, 0, s.Len())
	var cur = s.head
	for cur != nil {
		nodes = append(nodes, cur.v)
		cur = cur.next
	}
	return nodes
}
