package trie

import (
	"fmt"
	"log"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var size = 10000000

func TestTire(t *testing.T) {
	var trie = NewTrie(24, NewHexMap)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := []byte(primitive.NewObjectID().Hex())
		keys = append(keys, key)
		trie.Set([]byte(key), key)
	}
	log.Println("set cost", time.Since(start).Seconds())
	for _, k := range keys {
		v, ok := trie.Get(k)
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if ok, i := SliceEq(v.([]byte), k); !ok {
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

func TestGt(t *testing.T) {
	var trie = NewTrie(3, NewNmap)
	for i := 101; i < 200; i += 10 {
		key := []byte(fmt.Sprint(i))
		trie.Set(key, key)
	}
	v := trie.Gte([]byte("111"))
	fmt.Printf("%s\n", v)
}

func TestLt(t *testing.T) {
	var trie = NewTrie(3, NewNmap)
	for i := 101; i < 200; i += 10 {
		key := []byte(fmt.Sprint(i))
		trie.Set(key, key)
	}
	v := trie.Lte([]byte("191"))
	fmt.Printf("%s\n", v)
}

func TestScan(t *testing.T) {
	var trie = NewTrie(3, NewNmap)
	for i := 101; i < 200; i += 10 {
		key := []byte(fmt.Sprint(i))
		trie.Set(key, key)
	}
	trie.Del([]byte("131"))
	vals := trie.Scan(Include([]byte("121")), Exclude([]byte("221")))
	for _, v := range vals {
		fmt.Printf("%s\n", v)
	}

}
