package protohash

import (
	"crypto/sha256"
	"fmt"
	"math"
)

const hashLength int = sha256.Size

const (
	// Sorted alphabetically by value.
	boolIdentifier     = `b`
	mapIdentifier      = `d`
	floatIdentifier    = `f`
	intIdentifier      = `i`
	listIdentifier     = `l`
	nilIdentifier      = `n`
	byteIdentifier     = `t`
	unicodeIndentifier = `u`

	// Proto-specific identifiers.
	protoMessageIdentifier = `m`
)

func hash(t string, b []byte) ([]byte, error) {
	h := sha256.New()
	h.Write([]byte(t))
	h.Write(b)
	return h.Sum(nil), nil
}

func hashBool(b bool) ([]byte, error) {
	bb := []byte(`0`)
	if b {
		bb = []byte(`1`)
	}
	return hash(boolIdentifier, bb)
}

func hashBytes(bs []byte) ([]byte, error) {
	return hash(byteIdentifier, bs)
}

func hashFloat(f float64) ([]byte, error) {
	var normalizedFloat string

	switch {
	case math.IsInf(f, 1):
		normalizedFloat = "Infinity"
	case math.IsInf(f, -1):
		normalizedFloat = "-Infinity"
	case math.IsNaN(f):
		normalizedFloat = "NaN"
	default:
		var err error
		normalizedFloat, err = floatNormalize(f)
		if err != nil {
			return nil, err
		}
	}

	return hash(floatIdentifier, []byte(normalizedFloat))
}

func hashInt64(i int64) ([]byte, error) {
	return hash(intIdentifier, []byte(fmt.Sprintf("%d", i)))
}

func hashNil() ([]byte, error) {
	return hash(nilIdentifier, []byte(``))
}

func hashUint64(i uint64) ([]byte, error) {
	return hash(intIdentifier, []byte(fmt.Sprintf("%d", i)))
}

func hashUnicode(s string) ([]byte, error) {
	return hash(unicodeIndentifier, []byte(s))
}
