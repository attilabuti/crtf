package crtf

import (
	"bytes"
	"testing"
)

// Test decompression of compressed data.
func TestHitherAndThither(t *testing.T) {
	data := []byte("{\\rtf1\\ansi\\mac\\deff0\\deftab720")
	compressed := Compress(data, true)
	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(decompressed, data) {
		t.Errorf("Hither and thither test failed.\nOriginal: %v\nResult: %v", data, decompressed)
	}
}

// Test decompression of compressed data larger than 4096.
func TestHitherAndThitherLong(t *testing.T) {
	data := []byte("{\\rtf1\\ansi\\ansicpg1252\\pard hello world")
	for len(data) < 4096 {
		data = append(data, []byte("testtest")...)
	}
	data = append(data, '}')

	compressed := Compress(data, true)
	decompressed, err := Decompress(compressed)
	if err != nil {
		t.Fatalf("Decompress failed: %v", err)
	}

	if !bytes.Equal(decompressed, data) {
		t.Errorf("Hither and thither long test failed.\nOriginal length: %d\nResult length: %d", len(data), len(decompressed))
	}
}
