package fasthash

import (
	"math/rand/v2"
	"testing"
)

var testCases64 = []struct {
	in   string
	want uint64
}{
	{"0", 7749275010220701263},
	{"01", 10872400197931041294},
	{"012", 15346974464947871338},
	{"0123", 4981619345643336833},
	{"01234", 13689749573412676407},
	{"012345", 12368623765455266986},
	{"0123456", 12533847505615611486},
	{"01234567", 7278149109652352471},
	{"012345678", 732621452303077734},
	{"0123456789", 11460214490870572832},
	{"01234567890", 17034508007883094207},
	{"012345678901", 2613408247548219540},
	{"0123456789012", 9592538707932359806},
	{"01234567890123", 8026059100768657110},
	{"012345678901234", 10326712318572838968},
	{"0123456789012345", 8542297193267372205},
	{"01234567890123456", 9007936939456149167},
	{"012345678901234567", 4883792730719261022},
	{"0123456789012345678", 13174688786326369759},
	{"01234567890123456789", 4588891807043452244},
}

var testCases32 = []struct {
	in   string
	want uint32
}{
	{"0", 259604927},
	{"01", 2393799598},
	{"012", 4214879860},
	{"0123", 1780920347},
	{"01234", 3847594115},
	{"012345", 3234084578},
	{"0123456", 2201041528},
	{"01234567", 1780891594},
	{"012345678", 1093682706},
	{"0123456789", 3004500676},
	{"01234567890", 3775167992},
	{"012345678901", 1139498912},
	{"0123456789012", 526295389},
	{"01234567890123", 3740828},
	{"012345678901234", 1725916693},
	{"0123456789012345", 1408780964},
	{"01234567890123456", 509908932},
	{"012345678901234567", 3009501634},
	{"0123456789012345678", 1484279865},
	{"01234567890123456789", 3621222539},
}

func TestHash64(t *testing.T) {
	for _, tc := range testCases64 {
		got := Hash64(0, []byte(tc.in))
		if got != tc.want {
			t.Errorf("Hash64(0, %q): want %d, got %d", tc.in, tc.want, got)
		}
	}
}

func TestHash32(t *testing.T) {
	for _, tc := range testCases32 {
		got := Hash32(0, []byte(tc.in))
		if got != tc.want {
			t.Errorf("Hash32(0, %q): want %d, got %d", tc.in, tc.want, got)
		}
	}
}

// testCasesLong were produced by the original reflect.SliceHeader-based
// hash64, over buf[i] = byte(i*131 + 7).
var testCasesLong = []struct {
	seed uint64
	n    int
	want uint64
}{
	{0x0, 63, 14668168003056199637},
	{0x0, 64, 5797897923055584387},
	{0x123456789abcdef, 65, 5385230469984326604},
	{0x123456789abcdef, 255, 2723948467434838504},
	{0x2a, 1000, 16971330316647270956},
}

func TestHash64Long(t *testing.T) {
	buf := make([]byte, 1000)
	for i := range buf {
		buf[i] = byte(i*131 + 7)
	}
	for _, tc := range testCasesLong {
		got := Hash64(tc.seed, buf[:tc.n])
		if got != tc.want {
			t.Errorf("Hash64(%#x, buf[:%d]): want %d, got %d", tc.seed, tc.n, tc.want, got)
		}
	}
}

// refHash64 is a direct byte-at-a-time transcription of fast-hash.
func refHash64(seed uint64, buf []byte) uint64 {
	mix := func(v uint64) uint64 {
		v ^= v >> 23
		v *= 0x2127599bf4325c37
		v ^= v >> 47
		return v
	}
	h := seed ^ (uint64(len(buf)) * m)
	for len(buf) > 0 {
		n := min(len(buf), 8)
		var v uint64
		for i := n - 1; i >= 0; i-- {
			v = v<<8 | uint64(buf[i])
		}
		h ^= mix(v)
		h *= m
		buf = buf[n:]
	}
	return mix(h)
}

// TestHash64Ref checks Hash64 (assembly on amd64) and hash64 against
// refHash64 for every length up to 1 KiB, at every alignment.
func TestHash64Ref(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	buf := make([]byte, 1024+8)
	for i := range buf {
		buf[i] = byte(r.Uint32())
	}
	for off := 0; off < 8; off++ {
		for n := 0; off+n <= len(buf); n++ {
			b := buf[off : off+n]
			seed := r.Uint64()
			want := refHash64(seed, b)
			if got := Hash64(seed, b); got != want {
				t.Fatalf("Hash64 off=%d len=%d: want %d, got %d", off, n, want, got)
			}
			if got := hash64(seed, b); got != want {
				t.Fatalf("hash64 off=%d len=%d: want %d, got %d", off, n, want, got)
			}
		}
	}
}

func TestHash64NoAlloc(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		var b [64]byte
		sink = Hash64(0, b[:])
	})
	if allocs != 0 {
		t.Errorf("Hash64 allocates %v times per call", allocs)
	}
}
