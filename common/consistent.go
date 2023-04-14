package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type uints []uint32

// Len : slice length
func (x uints) Len() int {
	return len(x)
}

// Less : compare two uint32 values
func (x uints) Less(i, j int) bool {
	return x[i] < x[j]
}

// Swap : swap two values in the slice
func (x uints) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

// Prompt for error when there is no data on the hash ring
var errEmpty = errors.New("Empty Hash Ring")

type Consistent struct {
	// hash ring: key = hash value, value = node information
	circle map[uint32]string
	// sorted node hash slices
	sortedHashes uints
	// number of virtual nodes: used to increase the balance of hashing
	VirtualNode int
	// Map ReadWriteLock
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		VirtualNode: 25,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		copy(scratch[:], key)
		// using IEEE polynomials to return CRC-32 checksum of data
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHash() {
	hashes := c.sortedHashes[:0]
	// Determine if the slice capacity is too large, and reset if true
	if cap(c.sortedHashes)/(c.VirtualNode) > len(c.circle) {
		hashes = nil
	}
	for k := range c.circle {
		hashes = append(hashes, k)
	}
	// sort for binary search
	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func (c *Consistent) add(element string) {
	// Loop virtual nodes and set replicas
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	c.updateSortedHash()
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHash()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// Find the closest node clockwise
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	//Use binary search to search for the minimum value that meets the conditions for a specified slice
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.Lock()
	defer c.Unlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
