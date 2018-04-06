package fasthash

import (
	"runtime"
	"testing"
)

var buf = make([]byte, 1024)

func BenchmarkHash64Generic(b *testing.B) {
	b.SetBytes(int64(len(buf)))
	for i := 0; i < b.N; i++ {
		hash64(0, buf)
	}
}

func BenchmarkHash64(b *testing.B) {
	if runtime.GOARCH != "amd64" {
		b.SkipNow()
	}
	b.SetBytes(int64(len(buf)))
	for i := 0; i < b.N; i++ {
		Hash64(0, buf)
	}
}
