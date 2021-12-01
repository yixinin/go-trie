package nodetrie

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/huandu/skiplist"
)

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
	time.Sleep(time.Second * 5)
}

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
	time.Sleep(time.Second * 5)
}
