# go-fasthash

[![Go Reference](https://pkg.go.dev/badge/github.com/opencoff/go-fasthash.svg)](https://pkg.go.dev/github.com/opencoff/go-fasthash)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0.html)

A Go implementation of **fast-hash**, Zilong Tan's simple, robust, and
efficient general-purpose (non-cryptographic) hash function.

## Attribution

fast-hash was designed by **Zilong Tan**, who derived it using genetic
programming. The original C implementation, `fasthash64` and `fasthash32`,
is at <https://github.com/ztanml/fast-hash> (MIT License, Copyright (C) 2012
Zilong Tan).

This package is a Go implementation of that algorithm. It produces the
same values as the C reference on little-endian machines; the test vectors
in `fasthash_test.go` were checked against the C code.

## Install

    go get github.com/opencoff/go-fasthash

## Usage

```go
import "github.com/opencoff/go-fasthash"

h64 := fasthash.Hash64(seed, data) // seed uint64, data []byte
h32 := fasthash.Hash32(seed, data) // seed uint32, data []byte
```

`Hash32` folds the 64-bit hash to 32 bits the same way the C `fasthash32`
does.

## Quality

According to its author, fast-hash passes every test in the original
[SMHasher](https://github.com/aappleby/smhasher) suite. In Reini Urban's
extended [SMHasher](https://github.com/rurban/smhasher), `fasthash32` is
listed among the fastest hash functions without quality problems. The only
flag on `fasthash32` and `fasthash64` there is "UB", for the C code's
misaligned 64-bit loads. This package reads its input with `encoding/binary`,
so any buffer alignment is safe.

fast-hash is not a cryptographic hash.

## Implementation notes

- On amd64, `Hash64` uses hand-written assembly. Everywhere else it is
  pure Go.
- The Go implementation is endian-independent. It produces the same values
  on little- and big-endian machines, and those values match the C reference
  on little-endian machines.

## Benchmarks

    go test -bench .

Each input size runs in two modes:

- `thru`: independent calls, which the CPU can overlap (hashing a batch of
  keys).
- `lat`: each call's seed is the previous result, so this measures one hash
  from start to finish (e.g. a single table lookup).

On amd64, `BenchmarkHash64Generic` benchmarks the pure-Go code as a baseline
for the assembly.

## License

This package is licensed under the GNU General Public License v3; see
[LICENSE](LICENSE). The original C implementation by Zilong Tan is
MIT-licensed.
