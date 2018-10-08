package bloomfilter

import (
	"encoding/binary"
	"math"
)

type BloomFilter struct {
	m       uint32
	k       int
	buckets []uint32
}

// New creates a new bloom filter. m should specify the number of bits.
// m is rounded up to the nearest multiple of 32.
// k specifies the number of hashing functions.
func New(m, k int) BloomFilter {
	var n = uint32(math.Ceil(float64(m) / 32))
	bf := BloomFilter{
		m:       n * 32,
		k:       k,
		buckets: make([]uint32, n),
	}
	return bf
}

// NewFromUint32Slice creates a new bloom filter from a int32 slice.
// b is a bucket set.
// k specifies the number of hashing functions.
func NewFromUint32Slice(ii []uint32, k int) BloomFilter {
	bf := BloomFilter{
		m:       uint32(len(ii) * 32),
		k:       k,
		buckets: ii,
	}
	return bf
}

// NewFromBytes creates a new bloom filter from a byte slice.
// b is a byte slice exported from another bloomfilter.
// k specifies the number of hashing functions.
func NewFromBytes(bb []byte, k int) BloomFilter {
	ii := make([]uint32, len(bb)/4)
	for i := range ii {
		ii[i] = uint32(binary.BigEndian.Uint32(bb[i*4 : (i+1)*4]))
	}
	return NewFromUint32Slice(ii, k)
}

// EstimateParameters estimates requirements for m and k.
// https://github.com/willf/bloom
func EstimateParameters(n int, p float64) (m int, k int) {
	m = int(math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k = int(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	if m%32 > 0 {
		m += 32 - m%32
	}
	return
}

func (bf BloomFilter) locations(v []byte) []uint32 {
	var r = make([]uint32, bf.k)
	var a = fnv_1a(v, 0)
	var b = fnv_1a(v, 1576284489)
	var x = a % uint32(bf.m)
	for i := 0; i < bf.k; i++ {
		if x < 0 {
			r[i] = x + bf.m
		} else {
			r[i] = x
		}
		x = (x + b) % bf.m
	}
	return r
}

// Add adds a byte array to the bloom filter
func (bf BloomFilter) Add(v []byte) {
	var l = bf.locations(v)
	for i := 0; i < bf.k; i++ {
		bf.buckets[int(math.Floor(float64(uint32(l[i])/32)))] |= 1 << (uint32(l[i]) % 32)
	}
}

// AddInt adds an int to the bloom filter
func (bf BloomFilter) AddInt(v int) {
	var a = make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(v))
	bf.Add(a)
}

// Test evaluates a byte array to determine whether it is (probably) in the bloom filter
func (bf BloomFilter) Test(v []byte) bool {
	var l = bf.locations(v)
	for i := 0; i < bf.k; i++ {
		if (bf.buckets[int(math.Floor(float64(uint32(l[i])/32)))] & (1 << (uint32(l[i]) % 32))) == 0 {
			return false
		}
	}
	return true
}

// TestInt evaluates an int to determine whether it is (probably) in the bloom filter
func (bf BloomFilter) TestInt(v int) bool {
	var a = make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(v))
	return bf.Test(a)
}

// ToBytes returns the bloom filter as a byte slice
func (bf BloomFilter) ToBytes() []byte {
	var bb = []byte{}
	for i := 0; i < len(bf.buckets); i++ {
		var a = make([]byte, 4)
		binary.BigEndian.PutUint32(a, bf.buckets[i])
		bb = append(bb, a...)
	}
	return bb
}

// ToUint32Slice returns the bloom filter as a uint32 slice
func (bf BloomFilter) ToUint32Slice() []uint32 {
	return bf.buckets
}

// Fowler/Noll/Vo hashing.
// Nonstandard variation: this function optionally takes a seed value that is incorporated
// into the offset basis. According to http://www.isthe.com/chongo/tech/comp/fnv/index.html
// "almost any offset_basis will serve so long as it is non-zero".
func fnv_1a(v []byte, seed int) uint32 {
	var a = uint32(2166136261 ^ seed)
	var n = len(v)
	for i := 0; i < n; i++ {
		var c = uint32(v[i])
		var d = c & 0xff00
		if d != 0 {
			a = fnv_multiply(a ^ d>>8)
		}
		a = fnv_multiply(a ^ c&0xff)
	}
	return fnv_mix(a)
}

// a * 16777619 mod 2**32
func fnv_multiply(a uint32) uint32 {
	return a + (a << 1) + (a << 4) + (a << 7) + (a << 8) + (a << 24)
}

// See https://web.archive.org/web/20131019013225/http://home.comcast.net/~bretm/hash/6.html
func fnv_mix(a uint32) uint32 {
	a += a << 13
	a ^= a >> 7
	a += a << 3
	a ^= a >> 17
	a += a << 5
	return a & 0xffffffff
}
