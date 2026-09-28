package fasthash

import (
	"fmt"
	"runtime"
	"testing"
)

var benchSizes = []int{3, 8, 16, 32, 64, 256, 1024, 64 << 10}

var sink uint64

// BenchmarkHash64 measures the exported Hash64 (assembly on amd64).
func BenchmarkHash64(b *testing.B) {
	benchHash(b, Hash64)
}

// BenchmarkHash64Generic measures the pure-Go hash64, the baseline for
// the amd64 assembly. Elsewhere Hash64 is hash64, so it is skipped.
func BenchmarkHash64Generic(b *testing.B) {
	if runtime.GOARCH != "amd64" {
		b.Skip("same as BenchmarkHash64 on", runtime.GOARCH)
	}
	benchHash(b, hash64)
}

// benchHash runs each size in two modes:
//
//	thru: independent calls, which the CPU can overlap (hashing a batch).
//	lat:  each call's seed is the previous result, so this is the time
//	      for one hash end to end (e.g. a single table lookup).
func benchHash(b *testing.B, f func(uint64, []byte) uint64) {
	for _, n := range benchSizes {
		buf := make([]byte, n)
		b.Run(fmt.Sprintf("thru/%d", n), func(b *testing.B) {
			b.SetBytes(int64(n))
			var s uint64
			for b.Loop() {
				s += f(0, buf)
			}
			sink = s
		})
		b.Run(fmt.Sprintf("lat/%d", n), func(b *testing.B) {
			b.SetBytes(int64(n))
			var s uint64
			for b.Loop() {
				s = f(s, buf)
			}
			sink = s
		})
	}
}
