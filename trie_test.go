package trie

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/huandu/skiplist"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var size = 10000000

func TestBytesTire(t *testing.T) {
	var trie = NewTrie(12, NewByteMap)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := primitive.NewObjectID()
		k := key[:]
		keys = append(keys, k)
		trie.Set(k, k)
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

func TestHexTire(t *testing.T) {
	var trie = NewTrie(24, NewHexMap)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := []byte(primitive.NewObjectID().Hex())
		k := key[:]
		keys = append(keys, k)
		trie.Set(k, k)
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

	fmt.Println(trie.Get([]byte(primitive.NewObjectID().Hex())))
}

func TestMapTire(t *testing.T) {
	var trie = NewTrie(12, NewLinkmap)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := primitive.NewObjectID()
		k := key[:]
		keys = append(keys, k)
		trie.Set(k, k)
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

func TestMap(t *testing.T) {
	var m = make(map[string]string, size)
	var keys = make([]string, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := primitive.NewObjectID()
		k := key.Hex()
		keys = append(keys, k)
		m[k] = k
	}
	log.Println("set cost", time.Since(start).Seconds())
	for _, k := range keys {
		v, ok := m[k]
		if !ok {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if ok, i := SliceEq([]byte(v), []byte(k)); !ok {
			t.Logf("key %d: %s-%s fails\n", i, k, v)
			t.Fail()
		}
	}
	if len(m) != len(keys) {
		t.Log("size not eq")
		t.Fail()
	}
	log.Println("total cost", time.Since(start).Seconds())
}

func TestSkipList(t *testing.T) {
	var trie = skiplist.New(skiplist.Bytes)
	var keys = make([][]byte, 0, size)
	start := time.Now()
	for i := 0; i < size; i++ {
		key := primitive.NewObjectID()
		k := key[:]
		keys = append(keys, k)
		trie.Set(k, k)
	}
	log.Println("set cost", time.Since(start).Seconds())
	for _, k := range keys {
		v := trie.Get(k)
		if v == nil {
			t.Logf("no key:%v fail\n", k)
			t.Fail()
			continue
		}
		if ok, i := SliceEq(v.Value.([]byte), k); !ok {
			t.Logf("key %d: %s-%s fails\n", i, k, v.Value)
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
