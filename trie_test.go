package trie

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	var keys = []string{
		"1000001",
		"2000001",
		// "2000002",
		// "1000003",
		// "1000004",
		"2000005",
		"4000006",
	}
	var trie = NewTrie()
	for _, key := range keys {
		trie.Insert([]byte(key), "v")
	}
	var l, r []byte
	var key = "2100002"
	l = trie.FindPrev([]byte(key))
	r = trie.FindNext([]byte(key))
	fmt.Printf("key:%s, l:%s, r:%s\n", key, l, r)
	trie.Foreach(func(key []byte) {
		fmt.Println(string(key))
	})
}
