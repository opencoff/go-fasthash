// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package fasthash provides fast-hash (a simple, robust, and efficient general-purpose hash function) implementation in Go.
package fasthash

import "encoding/binary"

const (
	m    = 0x880355f21e6d1965
	mixK = 0x2127599bf4325c37
)

// mixMul is mixK held in a variable rather than a constant: the compiler
// then keeps it in a register across the hash64 loop instead of rebuilding
// the 64-bit immediate on every iteration (MOVD + 3x MOVK on arm64).
var mixMul uint64 = mixK

func mix(v, k uint64) uint64 {
	v ^= v >> 23
	v *= k
	v ^= v >> 47
	return v
}

func hash64(seed uint64, buf []byte) uint64 {
	k := mixMul
	h := seed ^ (uint64(len(buf)) * m)

	// This loop shape (i < len-7, buf[i:i+8]) lets the compiler prove
	// every access in bounds, so the loop has no bounds checks.
	for i := 0; i < len(buf)-7; i += 8 {
		h ^= mix(binary.LittleEndian.Uint64(buf[i:i+8]), k)
		h *= m
	}

	if buf = buf[len(buf)&^7:]; len(buf) > 0 {
		h ^= mix(tail(buf), k)
		h *= m
	}
	return mix(h, k)
}

// tail returns the 1..7 bytes of b as a little-endian uint64.
func tail(b []byte) uint64 {
	n := len(b)
	if n >= 4 {
		// two overlapping 4-byte loads: b[0:4] and b[n-4:n]
		lo := uint64(binary.LittleEndian.Uint32(b))
		hi := uint64(binary.LittleEndian.Uint32(b[n-4:]))
		return lo | (hi>>(8*(8-n)))<<32
	}
	// b[0], b[n/2] and b[n-1] cover every byte for n = 1..3
	return uint64(b[0]) | uint64(b[n>>1])<<(8*(n>>1)) | uint64(b[n-1])<<(8*(n-1))
}

func Hash32(seed uint32, buf []byte) uint32 {
	h := Hash64(uint64(seed), buf)
	return uint32(h - h>>32)
}
