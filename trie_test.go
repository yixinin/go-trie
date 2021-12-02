package trie

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/huandu/skiplist"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// === RUN   TestTire
// 0.1475573
// --- PASS: TestTire (0.16s)
// PASS
// ok      gotrie/nodetrie 0.216s
func TestTire(t *testing.T) {
	var size = 1000000
	var trie = NewTrie(24, NewHexMap)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := []byte(primitive.NewObjectID().Hex())
		keys = append(keys, key)
		trie.Set(key, key)
	}
	for _, k := range keys {
		log.Println(trie.Get(k))
	}
	fmt.Println(time.Since(start).Seconds())
	// time.Sleep(time.Second * 5)
}

// === RUN   TestSkipList
// 0.1752606
// --- PASS: TestSkipList (0.18s)
// PASS
// ok      gotrie/nodetrie 0.218s

func TestSkipList(t *testing.T) {
	var size = 1000000
	var list = skiplist.New(skiplist.Bytes)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := []byte(primitive.NewObjectID().Hex())
		keys = append(keys, key)
		list.Set(key, key)
	}
	for _, k := range keys {
		list.Get(k)
	}
	fmt.Println(time.Since(start).Seconds())
	// time.Sleep(time.Second * 5)
}

func TestArrayMem(t *testing.T) {
	var array [10000000]*TrieNode

	fmt.Println(len(array), array[0])
	time.Sleep(time.Second * 10)
	for i := range array {
		array[i] = NewTrieNode('a', []byte{1, 2, 3}, "aaa", NewHexMap)
	}
	time.Sleep(time.Second * 10)
}
