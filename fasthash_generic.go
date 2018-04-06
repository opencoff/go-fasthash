// +build !amd64

package fasthash

func Hash64(seed uint64, buf []byte) uint64 {
	return hash64(seed, buf)
}
