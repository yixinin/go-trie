package nodetrie

import (
	"fmt"
	"testing"
)

func TestTire(t *testing.T) {
	var trie = NewNodeTrie()
	var keys = []string{

		"1000002",
		"1000003",
		"1000004",
		"2000001",
		"1000001",
		"3000002",
		"4000003",
		"9000004",
		"5000004",
	}
	for _, key := range keys {
		trie.Insert([]byte(key))
	}
	trie.Foreach(func(k []byte) {
		fmt.Println(string(k))
	})
}
