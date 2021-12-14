package trie

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
)

const (
	FsTrieNodeSize = 256*8 + 24 + 1
)

type MemNode struct {
	prev     *FsTrieNode
	next     *FsTrieNode
	nodeKey  byte
	children FsContainer
}

type DiskNode struct {
	self uint64
	prev uint64
	next uint64
	key  []byte
	val  []byte
}
type FsTrieNode struct {
	*MemNode
	*DiskNode
	offset *uint64
	isLeaf bool
}

func newMemNode(k byte, container func() FsContainer) *FsTrieNode {
	return &FsTrieNode{
		MemNode: &MemNode{
			nodeKey:  k,
			children: container(),
		},
	}
}

func newDiskNode(key, val []byte, offset uint64) *FsTrieNode {
	return &FsTrieNode{
		DiskNode: &DiskNode{
			self: offset,
			key:  key,
			val:  val,
		},
		MemNode: &MemNode{
			nodeKey: key[len(key)-1],
		},
		offset: &offset,
		isLeaf: true,
	}
}

func (n *FsTrieNode) getNodeKey() byte {
	if n.isLeaf {
		return n.DiskNode.key[len(n.DiskNode.key)-1]
	}
	return n.MemNode.nodeKey
}

func (n *FsTrieNode) initFsTrieLeaf(offset uint64, key []byte, val []byte) *FsTrieNode {
	node := n.DiskNode
	node.self = offset
	node.key = key
	node.val = val
	n.isLeaf = true
	return n
}

func (n *FsTrieNode) unmarshalLeaf(keySize int, buf []byte) {
	node := n.DiskNode
	node.self = binary.BigEndian.Uint64(buf[8:16])
	node.prev = binary.BigEndian.Uint64(buf[16:24])
	node.next = binary.BigEndian.Uint64(buf[24:32])
	node.key = buf[32 : 32+keySize]
	node.val = buf[32+keySize:]
}
func (n *FsTrieNode) marshalLeaf() []byte {
	// 计算大小
	node := n.DiskNode
	var keySize = len(node.key)
	var valSize = len(node.val)
	var size = 32 + keySize + valSize
	var buf = make([]byte, size)
	binary.BigEndian.PutUint64(buf[:8], uint64(size))
	binary.BigEndian.PutUint64(buf[8:16], node.self)
	binary.BigEndian.PutUint64(buf[16:24], node.prev)
	binary.BigEndian.PutUint64(buf[24:32], node.next)
	copy(buf[32:32+keySize], node.key)
	copy(buf[32+keySize:32+keySize+valSize], node.val)
	return buf
}

type FsTrie struct {
	root      *FsTrieNode
	fs        *os.File
	nodePool  *sync.Pool
	keySize   int
	size      int
	head      uint64
	tail      uint64
	blockSize uint64
	fileSize  uint64
	container func() FsContainer
}

func NewFsTrie(filename string, keySize int, container func() FsContainer) (*FsTrie, error) {
	var (
		fstat os.FileInfo
		err   error
	)
	var t = &FsTrie{
		blockSize: 4096,
		keySize:   keySize,
		root:      newMemNode(0, container),
		container: container,
	}
	t.nodePool = &sync.Pool{
		New: func() interface{} {
			return &FsTrieNode{}
		},
	}
	if t.fs, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return nil, err
	}

	if fstat, err = t.fs.Stat(); err != nil {
		return nil, err
	}
	t.fileSize = uint64(fstat.Size())
	if t.fileSize > 0 {
		return t, nil
	}

	return t, nil
}

func (t *FsTrie) Set(key []byte, val []byte) error {
	if len(key) != t.keySize {
		return fmt.Errorf("key size should be %s", strconv.Itoa(t.keySize))
	}
	cur := t.root
	var ok bool
	for level, nodeKey := range key {
		var isLeaf = level == t.keySize-1
		var node *FsTrieNode
		if node, ok = cur.children.Get(nodeKey); !ok {
			if isLeaf {
				node = newDiskNode(key, val, t.fileSize)
				t.size++
			} else {
				node = newMemNode(nodeKey, t.container)
			}

			cur.children.Set(nodeKey, node)
			if tail := cur.children.Tail(); tail == node {
				if next := cur.MemNode.next; next != nil {
					if nextHead := next.children.Head(); nextHead != nil {
						nextHead.MemNode.prev = node
						node.MemNode.next = nextHead
					}
				}
			}
			if head := cur.children.Head(); head == node {
				if prev := cur.MemNode.prev; prev != nil {
					if prevTail := prev.children.Tail(); prevTail != nil {
						prevTail.MemNode.next = node
						node.MemNode.prev = prevTail
					}
				}
			}
		} else if isLeaf {
			node.val = val
		}
		if isLeaf {
			t.saveTrieNode(node)
			if t.head == 0 {
				t.head = node.DiskNode.self
				t.tail = node.DiskNode.self
				return nil
			}
			if t.head == node.self {
				t.head = node.self
			}
			if t.tail == node.self {
				t.tail = node.self
			}
			node.DiskNode = nil
		}
		cur = node
	}
	return nil
}

func (t *FsTrie) Get(key []byte) ([]byte, error) {
	if len(key) != t.keySize {
		return nil, fmt.Errorf("key size should be %s", strconv.Itoa(t.keySize))
	}
	cur := t.root
	var ok bool
	for level, k := range key {
		isLeaf := level == t.keySize-1
		cur, ok = cur.children.Get(k)
		if !ok {
			return nil, errors.New("not found")
		}
		if isLeaf {
			leaf, err := t.readLeaf(*cur.offset)
			if err != nil {
				return nil, err
			}
			val := leaf.val
			leaf.DiskNode = nil
			return val, nil
		}
	}
	return nil, errors.New("not found")
}

func (t *FsTrie) readLeaf(offset uint64) (*FsTrieNode, error) {
	var buf [8]byte
	n, err := t.fs.ReadAt(buf[:], int64(offset))
	if err != nil {
		return nil, err
	}
	if n != 8 {
		return nil, errors.New("leaf node head size not match")
	}
	var size = binary.BigEndian.Uint64(buf[:])
	var bbuf = make([]byte, size)
	n, err = t.fs.ReadAt(bbuf, int64(offset))
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, errors.New("leaf node data size not match")
	}
	var node = &FsTrieNode{DiskNode: &DiskNode{}}
	node.unmarshalLeaf(t.keySize, bbuf)
	return node, nil
}

func (t *FsTrie) saveTrieNode(node *FsTrieNode) error {
	var buf = node.marshalLeaf()
	_, err := t.fs.WriteAt(buf[:], int64(node.self))
	if err != nil {
		return err
	}
	if fileSize := node.self + uint64(len(buf)); t.fileSize < fileSize {
		t.fileSize = fileSize
	}
	return nil
}

func (t *FsTrie) Len() int {
	return t.size
}
