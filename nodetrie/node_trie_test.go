package nodetrie

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/huandu/skiplist"
)

// === RUN   TestTire
// 0.1475573
// --- PASS: TestTire (0.16s)
// PASS
// ok      gotrie/nodetrie 0.216s
func TestTire(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var trie = NewNodeTrie(7)
	var keys = make([][]byte, 0, 10000000)
	start := time.Now()
	for i := 0; i < 100000; i++ {
		key := []byte(strconv.Itoa(rand.Intn(9000000) + 999999))
		keys = append(keys, key)
		trie.Set(key, "a")
	}
	for _, k := range keys {
		trie.Get(k)

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
	var list = skiplist.New(skiplist.String)
	var keys = make([]string, 0, 10000000)
	start := time.Now()
	for i := 0; i < 100000; i++ {
		key := strconv.Itoa(rand.Intn(9000000) + 999999)
		keys = append(keys, key)
		list.Set(key, "a")
	}
	for _, k := range keys {
		list.Get(k)
	}
	fmt.Println(time.Since(start).Seconds())
	// time.Sleep(time.Second * 5)
}

func TestArrayMem(t *testing.T) {
	var array [10000000]byte
	fmt.Println(len(array), array[0])
	time.Sleep(time.Second * 10)
}
