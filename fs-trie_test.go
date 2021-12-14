package trie

import (
	"fmt"
	"testing"
)

func TestFsTire(t *testing.T) {
	var filename = "t.db"
	var trie, err = NewFsTrie(filename, 3)
	if err != nil {
		t.Fatal(err)
	}
	err = trie.Set([]byte{1, 2, 3}, []byte{1, 2, 3})
	if err != nil {
		t.Fatal(err)
	}
	err = trie.Set([]byte{1, 2, 4}, []byte{1, 2, 4})
	if err != nil {
		t.Fatal(err)
	}
	val, err := trie.Get([]byte{1, 2, 4})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(val)
}
