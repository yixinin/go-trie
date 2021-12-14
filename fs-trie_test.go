package trie

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestFsTire(t *testing.T) {
	var filename = "t.db"
	if err := os.Remove(filename); err != nil {
		if !strings.Contains(err.Error(), "cannot find") {
			t.Error(err)
			return
		}
	}
	var trie, err = NewFsTrie(filename, 3)
	if err != nil {
		t.Error(err)
		return
	}

	var keys = []string{
		"123",
		"124",
	}
	err = trie.Set([]byte(keys[0]), []byte(keys[0]))
	if err != nil {
		t.Error(err)
		return
	}
	err = trie.Set([]byte(keys[1]), []byte(keys[1]))
	if err != nil {
		t.Error(err)
		return
	}
	val, err := trie.Get([]byte(keys[0]))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(string(val))
}
