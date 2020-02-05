package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uints []uint32

func (x uints) Len() int {
	return len(x)
}

func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("the hash circle is empty")

type Consistent struct {
	circle      map[uint32]string
	sortedHash  uints
	virtualNode int
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		virtualNode: 20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashkey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	for i := 0; i < c.virtualNode; i++ {
		c.circle[c.hashkey(c.generateKey(element, i))] = element
	}
	c.updateSortedHash()
}

func (c *Consistent) updateSortedHash() {
	sortedHash := c.sortedHash[:0]
	if cap(sortedHash)/c.virtualNode > len(c.circle) {
		sortedHash = nil
	}
	for k := range c.circle {
		sortedHash = append(sortedHash, k)
	}

	sort.Sort(sortedHash)
	c.sortedHash = sortedHash
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.virtualNode; i++ {
		delete(c.circle, c.hashkey(c.generateKey(element, i)))
	}
	c.updateSortedHash()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	c.RUnlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}

	key := c.hashkey(name)
	index := c.search(key)
	return c.circle[c.sortedHash[index]], nil

}

func (c *Consistent) search(key uint32) int {
	f := func(i int) bool { return c.sortedHash[i] > key }
	i := sort.Search(len(c.sortedHash), f)
	if i >= len(c.sortedHash) {
		return 0
	}
	return i
}
