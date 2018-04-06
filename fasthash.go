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

import (
	"reflect"
	"unsafe"
)

func mix(v uint64) uint64 {
	v ^= v >> 23
	v *= 0x2127599bf4325c37
	v ^= v >> 47
	return v
}

func hash64(seed uint64, buf []byte) uint64 {
	const m = 0x880355f21e6d1965
	h := seed ^ (uint64(len(buf)) * m)

	if n := len(buf) / 8; n > 0 {
		hdr := reflect.SliceHeader{
			Data: uintptr(unsafe.Pointer(&buf[0])),
			Len:  n,
			Cap:  n,
		}
		data := *(*[]uint64)(unsafe.Pointer(&hdr))

		for _, v := range data {
			h ^= mix(v)
			h *= m
		}
		buf = buf[n*8:]
	}

	var v uint64
	switch len(buf) {
	case 7:
		v ^= uint64(buf[6]) << 48
		fallthrough
	case 6:
		v ^= uint64(buf[5]) << 40
		fallthrough
	case 5:
		v ^= uint64(buf[4]) << 32
		fallthrough
	case 4:
		v ^= uint64(buf[3]) << 24
		fallthrough
	case 3:
		v ^= uint64(buf[2]) << 16
		fallthrough
	case 2:
		v ^= uint64(buf[1]) << 8
		fallthrough
	case 1:
		v ^= uint64(buf[0])
		h ^= mix(v)
		h *= m
	}

	return mix(h)
}

func Hash32(seed uint32, buf []byte) uint32 {
	h := Hash64(uint64(seed), buf)
	return uint32(h - h>>32)
}
