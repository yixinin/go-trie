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

type FsTrieNode struct {
	self     uint64
	prev     uint64
	next     uint64
	nodeKey  byte
	key      []byte
	val      []byte
	children [256]uint64
	isLeaf   bool
}

func NewFsTrieNode(k byte, offset uint64) *FsTrieNode {
	var node = &FsTrieNode{
		self:    offset,
		nodeKey: k,
		isLeaf:  false,
	}
	return node
}

func NewFsTrieLeaf(k byte, offset uint64, key []byte, val []byte) *FsTrieNode {
	var node = &FsTrieNode{
		self:    offset,
		nodeKey: k,
		key:     key,
		val:     val,
		isLeaf:  true,
	}
	return node
}

func (node *FsTrieNode) unmarshal(buf [FsTrieNodeSize]byte) {
	node.self = binary.BigEndian.Uint64(buf[:8])
	node.prev = binary.BigEndian.Uint64(buf[8:16])
	node.next = binary.BigEndian.Uint64(buf[16:24])
	node.nodeKey = buf[24]
	for i := 0; i < 256; i++ {
		node.children[i] = binary.BigEndian.Uint64(buf[25+8*i : 25+8*i+8])
	}
}
func (node *FsTrieNode) marshal() []byte {
	if node.isLeaf {
		return node.marshalLeaf()
	}
	var buf [FsTrieNodeSize]byte
	binary.BigEndian.PutUint64(buf[:8], node.self)
	binary.BigEndian.PutUint64(buf[8:16], node.prev)
	binary.BigEndian.PutUint64(buf[16:24], node.next)
	buf[24] = node.nodeKey
	for i, v := range node.children {
		binary.BigEndian.PutUint64(buf[25+8*i:25+8*i+8], v)
	}
	return buf[:]
}

func (node *FsTrieNode) unmarshalLeaf(keySize int, buf []byte) {
	// var size = binary.BigEndian.Uint64(buf[:8])
	node.self = binary.BigEndian.Uint64(buf[:8])
	node.prev = binary.BigEndian.Uint64(buf[8:16])
	node.next = binary.BigEndian.Uint64(buf[16:24])
	node.nodeKey = buf[24]
	node.key = buf[25 : 25+keySize]
	node.val = buf[25+keySize:]
}
func (node *FsTrieNode) marshalLeaf() []byte {
	// 计算大小
	var keySize = len(node.key)
	var valSize = len(node.val)
	var size = 25 + keySize + valSize
	var buf = make([]byte, size+8)
	binary.BigEndian.PutUint64(buf[:8], uint64(size))
	binary.BigEndian.PutUint64(buf[8:16], node.self)
	binary.BigEndian.PutUint64(buf[16:24], node.prev)
	binary.BigEndian.PutUint64(buf[24:32], node.next)
	buf[32] = node.nodeKey
	copy(buf[33:33+keySize], node.key)
	copy(buf[33+keySize:33+keySize+valSize], node.val)
	return buf
}

type FsTrie struct {
	root      uint64
	fs        *os.File
	container func() Container
	nodePool  *sync.Pool
	keySize   int
	size      int
	head      uint64
	tail      uint64
	blockSize uint64
	fileSize  uint64
}

func NewFsTrie(filename string, keySize int) (*FsTrie, error) {
	var (
		fstat os.FileInfo
		err   error
	)
	var t = &FsTrie{
		blockSize: 4096,
		keySize:   keySize,
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
	var node = NewFsTrieNode(0, 0)
	if err := t.saveTrieNode(node); err != nil {
		return nil, err
	}
	return t, nil
}
func (t *FsTrie) Set(key []byte, val []byte) error {
	if len(key) != t.keySize {
		return fmt.Errorf("key size should be %s", strconv.Itoa(t.keySize))
	}
	cur, err := t.readNode(t.root)
	if err != nil {
		return err
	}
	for level, k := range key {
		isLeaf := level == t.keySize-1
		if cur.children[k] == 0 {
			// set new node
			var node *FsTrieNode
			if isLeaf {
				node = NewFsTrieLeaf(k, t.fileSize, key, val)
				t.size++
			} else {
				node = NewFsTrieNode(k, t.fileSize)
			}
			cur.children[node.nodeKey] = node.self
			// set prev
			if node.nodeKey > 0 {
				for i := int(node.nodeKey - 1); i >= 0; i-- {
					if cur.children[i] > 0 {
						node.prev = cur.children[i]
						// set prev next to node

						break
					}
				}
			}
			// set next
			if node.nodeKey < 255 {
				for i := int(node.nodeKey + 1); i < 255; i++ {
					if cur.children[i] > 0 {
						node.prev = cur.children[i]
						// set next prev to node
						next, err := t.readTireNode(cur.children[i], byte(i), isLeaf)
						if err != nil {
							return err
						}
						next.prev = node.self
						if err := t.saveTrieNode(next); err != nil {
							return err
						}
						break
					}
				}
			}

			// set prev if node is head
			if node.prev == 0 {
				if cur.prev > 0 {
					parentPrev, err := t.readTireNode(cur.prev, 0, isLeaf)
					if err != nil {
						return err
					}
					for i := 255; i >= 0; i-- {
						if parentPrev.children[i] > 0 {
							node.prev = parentPrev.children[i]
							// set prev next to node
							prev, err := t.readTireNode(parentPrev.children[i], byte(i), isLeaf)
							if err != nil {
								return err
							}
							prev.next = node.self
							if err := t.saveTrieNode(prev); err != nil {
								return err
							}
							break
						}
					}
				}
			}

			//set next if node is tail
			if node.next == 0 {
				if cur.next > 0 {
					parentNext, err := t.readTireNode(cur.next, 0, isLeaf)
					if err != nil {
						return err
					}
					for i := 0; i < 256; i++ {
						if parentNext.children[i] > 0 {
							node.next = parentNext.children[i]
							// set next prev to node
							next, err := t.readTireNode(parentNext.children[i], byte(i), isLeaf)
							if err != nil {
								return err
							}
							next.prev = node.self
							if err := t.saveTrieNode(next); err != nil {
								return err
							}
							break
						}
					}
				}
			}
			if err := t.saveTrieNode(node); err != nil {
				return err
			}
			if err := t.saveTrieNode(cur); err != nil {
				return err
			}
			cur = node
		} else {
			// read exsit node
			cur, err = t.readTireNode(cur.children[k], k, isLeaf)
			if err != nil {
				return err
			}
			if isLeaf {
				cur.val = val
				if err := t.saveTrieNode(cur); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (t *FsTrie) Get(key []byte) ([]byte, error) {
	if len(key) != t.keySize {
		return nil, fmt.Errorf("key size should be %s", strconv.Itoa(t.keySize))
	}
	cur, err := t.readNode(t.root)
	if err != nil {
		return nil, err
	}
	for level, k := range key {
		isLeaf := level == t.keySize-1
		if cur.children[k] == 0 {
			return nil, errors.New("not found")
		} else {
			// read exsit node
			cur, err = t.readTireNode(cur.children[k], k, isLeaf)
			if err != nil {
				return nil, err
			}
			if isLeaf {
				return cur.val, nil
			}
		}
	}
	return nil, errors.New("not found")
}

func (t *FsTrie) readTireNode(offset uint64, k byte, isLeaf bool) (*FsTrieNode, error) {
	var child *FsTrieNode
	var err error
	if isLeaf {
		child, err = t.readLeaf(offset)
	} else {
		child, err = t.readNode(offset)
	}
	if err != nil {
		return nil, err
	}
	if k > 0 && child.nodeKey != k {
		return nil, fmt.Errorf("node not match expect:%d, real:%d", k, child.nodeKey)
	}
	return child, nil
}

func (t *FsTrie) readNode(offset uint64) (*FsTrieNode, error) {
	var buf [FsTrieNodeSize]byte
	n, err := t.fs.ReadAt(buf[:], int64(offset))
	if err != nil {
		return nil, err
	}
	if n != FsTrieNodeSize {
		return nil, errors.New("node size not match")
	}
	var node = new(FsTrieNode)
	node.unmarshal(buf)
	return node, nil
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
	n, err = t.fs.ReadAt(bbuf, int64(offset)+8)
	if err != nil {
		return nil, err
	}
	if n != int(size) {
		return nil, errors.New("leaf node data size not match")
	}
	var node = new(FsTrieNode)
	node.unmarshalLeaf(t.keySize, bbuf)
	return node, nil
}

func (t *FsTrie) saveTrieNode(node *FsTrieNode) error {
	var buf = node.marshal()
	_, err := t.fs.WriteAt(buf[:], int64(node.self))
	if err != nil {
		return err
	}
	if fileSize := node.self + uint64(len(buf)); t.fileSize < fileSize {
		t.fileSize = fileSize
	}
	return nil
}
