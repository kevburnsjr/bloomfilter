package bloomfilter

import (
	"encoding/base64"
	"encoding/binary"
	// "fmt"
	"testing"
)

var jabberwocky = "`Twas brillig, and the slithy toves\n  Did gyre and gimble in the wabe:\nAll mimsy were the borogoves,\n  And the mome raths outgrabe.\n\n\"Beware the Jabberwock, my son!\n  The jaws that bite, the claws that catch!\nBeware the Jubjub bird, and shun\n  The frumious Bandersnatch!\"\n\nHe took his vorpal sword in hand:\n  Long time the manxome foe he sought --\nSo rested he by the Tumtum tree,\n  And stood awhile in thought.\n\nAnd, as in uffish thought he stood,\n  The Jabberwock, with eyes of flame,\nCame whiffling through the tulgey wood,\n  And burbled as it came!\n\nOne, two! One, two! And through and through\n  The vorpal blade went snicker-snack!\nHe left it dead, and with its head\n  He went galumphing back.\n\n\"And, has thou slain the Jabberwock?\n  Come to my arms, my beamish boy!\nO frabjous day! Callooh! Callay!'\n  He chortled in his joy.\n\n`Twas brillig, and the slithy toves\n  Did gyre and gimble in the wabe;\nAll mimsy were the borogoves,\n  And the mome raths outgrabe."

func TestBasic(t *testing.T) {
	f := New(1000, 4)
	n1 := []byte("Bess")
	n2 := []byte("Jane")
	f.Add(n1)
	if !f.Test(n1) {
		t.Fail()
	}
	if f.Test(n2) {
		t.Fail()
	}
}

func TestJabberwocky(t *testing.T) {
	f := New(1000, 4)
	n1 := []byte(jabberwocky)
	n2 := []byte(jabberwocky + "\n")
	f.Add(n1)
	if !f.Test(n1) {
		t.Fail()
	}
	if f.Test(n2) {
		t.Fail()
	}
}

func TestBasicUint32(t *testing.T) {
	f := New(1000, 4)
	n1 := []byte("\u0100")
	n2 := []byte("\u0101")
	n3 := []byte("\u0103")
	f.Add(n1)
	if !f.Test(n1) {
		t.Fail()
	}
	if f.Test(n2) {
		t.Fail()
	}
	if f.Test(n3) {
		t.Fail()
	}
}

func TestWtf(t *testing.T) {
	f := New(1000, 4)
	f.Add([]byte("abc"))
	if f.Test([]byte("wtf")) {
		t.Fail()
	}
}

func TestWorksWithIntegerTypes(t *testing.T) {
	f := New(1000, 4)
	f.AddInt(1)
	if !f.TestInt(1) {
		t.Fail()
	}
	if f.TestInt(2) {
		t.Fail()
	}
}

func TestToFromBytes(t *testing.T) {
	k := 4
	m := 1000
	f := New(m, k)
	f.Add([]byte("abc"))
	f.Add([]byte("def"))
	bb := f.ToBytes()
	f2 := NewFromBytes(bb, k)
	if !f2.Test([]byte("abc")) {
		t.Fail()
	}
	if f2.Test([]byte("ghi")) {
		t.Fail()
	}
}

func TestEstimateParameters(t *testing.T) {
	m, k := EstimateParameters(10000, 1e-6)
	if m != 287552 || k != 20 {
		t.Log(287552, m)
		t.Log(20, k)
		t.Fail()
	}
}

// TestCompatibility compares bloomfilter.js results to this package's results
func TestCompatibility(t *testing.T) {
	expected := "nnjzzw=="
	f := New(32, 21)
	f.Add([]byte("abc"))
	actual := base64.StdEncoding.EncodeToString(f.ToBytes())
	if expected != string(actual) {
		t.Log(expected, string(actual))
		t.Fail()
	}
	expected = "HgDwD4B4A8A="
	f = New(64, 21)
	f.Add([]byte("abc"))
	actual = string(base64.StdEncoding.EncodeToString(f.ToBytes()))
	if expected != actual {
		t.Log(expected, actual)
		t.Fail()
	}
	expected = "HkRD/wB/5AB/APRExEAHxA=="
	f = New(100, 21)
	f.Add([]byte("abc"))
	f.Add([]byte("def"))
	f.Add([]byte("ghi"))
	actual = string(base64.StdEncoding.EncodeToString(f.ToBytes()))
	if expected != actual {
		t.Log(expected, actual)
		t.Fail()
	}
}

func BenchmarkSeparateTestAndAdd(b *testing.B) {
	m, k := EstimateParameters(b.N, 1e-4)
	f := New(m, k)
	key := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		f.Test(key)
		f.Add(key)
	}
}

func BenchmarkSeparateAdd(b *testing.B) {
	m, k := EstimateParameters(b.N, 1e-4)
	f := New(m, k)
	key := make([]byte, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		binary.BigEndian.PutUint32(key, uint32(i))
		f.Add(key)
	}
}
