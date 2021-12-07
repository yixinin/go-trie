package trie

import (
	"testing"
)

func TestHexmap(t *testing.T) {
	var m = NewHexMap()
	var keys = []byte("0123456789abcdef")
	for _, k := range keys {
		m.Set(k, newTrieNode(k, nil, k, NewHexMap))
	}
	for _, k := range keys {
		v, ok := m.Get(k)
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if k != v.nodeKey {
			t.Logf("key:%v-%v fail\n", k, v.nodeKey)
			t.Fail()
		}
		if k != v.val.(byte) {
			t.Logf("val:%v-%v fail\n", k, v.val)
			t.Fail()
		}
	}
	if ok, i := SliceEq(keys, m.Keys()); !ok {
		t.Logf("key %d:%v\n%v fails\n", i, keys, m.Keys())
		t.Fail()
	}
}

func TestNmap(t *testing.T) {
	var m = NewNmap()
	var keys = []byte("0123456789")
	for _, k := range keys {
		m.Set(k, newTrieNode(k, nil, k, NewNmap))
	}
	for _, k := range keys {
		v, ok := m.Get(k)
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if k != v.nodeKey {
			t.Logf("key:%v-%v fail\n", k, v.nodeKey)
			t.Fail()
		}
		if k != v.val.(byte) {
			t.Logf("val:%v-%v fail\n", k, v.val)
			t.Fail()
		}
	}
	if ok, i := SliceEq(keys, m.Keys()); !ok {
		t.Logf("key %d:%v\n%v fails\n", i, keys, m.Keys())
		t.Fail()
	}
}

func TestLinkmap(t *testing.T) {
	var m = NewLinkmap()
	var keys = make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		keys = append(keys, byte(i))
	}
	for _, k := range keys {
		m.Set(k, newTrieNode(k, nil, k, NewLinkmap))
	}
	for _, k := range keys {
		v, ok := m.Get(k)
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if k != v.nodeKey {
			t.Logf("key:%v-%v fail\n", k, v.nodeKey)
			t.Fail()
		}
		if k != v.val.(byte) {
			t.Logf("val:%v-%v fail\n", k, v.val)
			t.Fail()
		}
	}
	if ok, i := SliceEq(keys, m.Keys()); !ok {
		t.Logf("key %d:%v\n%v fails\n", i, keys, m.Keys())
		t.Fail()
	}
}

func TestBytemap(t *testing.T) {
	var m = NewByteMap()
	var keys = make([]byte, 0, 256)
	for i := 0; i < 256; i++ {
		keys = append(keys, byte(i))
	}
	for _, k := range keys {
		m.Set(k, newTrieNode(k, nil, k, NewByteMap))
	}
	for _, k := range keys {
		v, ok := m.Get(k)
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if k != v.nodeKey {
			t.Logf("key:%v-%v fail\n", k, v.nodeKey)
			t.Fail()
		}
		if k != v.val.(byte) {
			t.Logf("val:%v-%v fail\n", k, v.val)
			t.Fail()
		}
	}
	if ok, i := SliceEq(keys, m.Keys()); !ok {
		t.Logf("key %d:%v\n%v fails\n", i, keys, m.Keys())
		t.Fail()
	}
}
