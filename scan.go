package trie

type Key struct {
	key     []byte
	include bool
}

func Include(key []byte) Key {
	return Key{
		key:     key,
		include: true,
	}
}

func Exclude(key []byte) Key {
	return Key{
		key:     key,
		include: false,
	}
}

func (t *Trie) Scan(l, r Key) []interface{} {
	if t.head == nil {
		return nil
	}
	lnode := t.gt(l.key, l.include)
	rnode := t.lt(r.key, r.include)
	if lnode == nil {
		lnode = t.head
	}
	if rnode == nil {
		rnode = t.tail
	}
	var vals = make([]interface{}, 0, 8)
	for cur := lnode; cur != nil; cur = cur.next {
		vals = append(vals, cur.val)
		if ok, _ := SliceEq(cur.key, rnode.key); ok {
			return vals
		}
	}
	return vals
}
