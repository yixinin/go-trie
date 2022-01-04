package trie

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestFsBytesTire(t *testing.T) {
	var filename = "trie.db"
	if err := os.Remove(filename); err != nil {
		if !strings.Contains(err.Error(), "cannot find") {
			t.Error(err)
			return
		}
	}
	var trie, err = NewFsTrie(filename, 12, NewFsLinkMap)
	if err != nil {
		t.Error(err)
		return
	}
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := primitive.NewObjectID()
		k := key[:]
		keys = append(keys, k)
		trie.Set(k, k)
	}
	log.Println("set cost", time.Since(start).Seconds())
	log.Println("key size", len(keys)*12)
	for _, k := range keys {
		v, err := trie.Get(k)
		if err != nil {
			t.Error(err)
			continue
		}
		if ok, i := SliceEq(v, k); !ok {
			t.Logf("key %d: %s-%s fails\n", i, k, v)
			t.Fail()
		}
	}
	if trie.Len() != len(keys) {
		t.Log("size not eq")
		t.Fail()
	}
	log.Println("total cost", time.Since(start).Seconds())
}

func TestMemTest(t *testing.T) {
	var x = make([]*FsTrieNode, size)
	for i := range x {
		x[i] = newDiskNode([]byte{0}, []byte{0}, 0)
	}
	fmt.Println(len(x))
	time.Sleep(time.Second * 5)
}

func TestBpTree(t *testing.T) {
	t.Skip()
	// var filename = "bptree.db"
	// if err := os.Remove(filename); err != nil {
	// 	if !strings.Contains(err.Error(), "cannot find") {
	// 		t.Error(err)
	// 		return
	// 	}
	// }
	// var trie, err = bptree.NewTree(filename)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }
	// var keys = make([]uint64, 0, size)
	// start := time.Now()
	// for i := 0; i < size; i++ {
	// 	key := primitive.NewObjectID()
	// 	k := key.Timestamp().Unix()
	// 	keys = append(keys, uint64(k))
	// 	trie.Insert(uint64(k), key.Hex())
	// }
	// log.Println("set cost", time.Since(start).Seconds())
	// for _, k := range keys {
	// 	v, err := trie.Find(k)
	// 	if err != nil {
	// 		t.Error(err)
	// 		continue
	// 	}
	// 	if k != k {
	// 		t.Logf("key %d: %s-%s fails\n", k, v, v)
	// 		t.Fail()
	// 	}
	// }

	// log.Println("total cost", time.Since(start).Seconds())
}
