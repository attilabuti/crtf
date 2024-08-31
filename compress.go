package crtf

import (
	"bytes"
	"encoding/binary"
	"io"
)

// Compress data using RTF compression algorithm.
// If compr (compressed) flag is false, data will be written uncompressed.
func Compress(data []byte, compr bool) []byte {
	var outputBuffer bytes.Buffer
	var compType [4]byte
	crcValue := make([]byte, 4)

	// Set init dict
	idict := make([]byte, maxDictSize)
	copy(idict, initDict)
	for i := initDictSize; i < maxDictSize; i++ {
		idict[i] = ' '
	}

	writeOffset := initDictSize

	if compr {
		compType = compressed

		inStream := bytes.NewReader(data)

		// Init params
		controlByte := byte(0)
		controlBit := 1
		tokenOffset := 0
		tokenBuffer := new(bytes.Buffer)
		dictOffset := 0
		longestMatch := 0

		for {
			// Find longest match
			dictOffset, longestMatch, writeOffset = findLongestMatch(idict, inStream, writeOffset)

			numBytesToRead := 1
			if longestMatch > 1 {
				numBytesToRead = longestMatch
			}

			char := make([]byte, numBytesToRead)
			_, err := inStream.Read(char)

			// End of input
			if err != nil {
				// Update params
				controlByte |= 1 << (controlBit - 1)
				controlBit += 1
				tokenOffset += 2

				// Add dict reference
				dictRef := (writeOffset & 0xfff) << 4
				binary.Write(tokenBuffer, binary.BigEndian, uint16(dictRef))

				// Add to output
				outputBuffer.WriteByte(controlByte)
				outputBuffer.Write(tokenBuffer.Bytes()[:tokenOffset])

				break
			} else {
				if longestMatch > 1 {
					// Update params
					controlByte |= 1 << (controlBit - 1)
					controlBit += 1
					tokenOffset += 2

					// Add dict reference
					dictRef := (dictOffset&0xfff)<<4 | (longestMatch-2)&0xf
					binary.Write(tokenBuffer, binary.BigEndian, uint16(dictRef))
				} else {
					// Character is not found in dictionary
					if longestMatch == 0 {
						idict[writeOffset] = char[0]
						writeOffset = (writeOffset + 1) % maxDictSize
					}

					// Update params
					controlByte |= 0 << (controlBit - 1)
					controlBit += 1
					tokenOffset += 1

					// Add literal
					binary.Write(tokenBuffer, binary.LittleEndian, char)
				}

				longestMatch = 0
				if controlBit > 8 {
					// Add to output
					outputBuffer.WriteByte(controlByte)
					outputBuffer.Write(tokenBuffer.Bytes()[:tokenOffset])

					// Reset params
					controlByte = 0
					controlBit = 1
					tokenOffset = 0
					tokenBuffer = new(bytes.Buffer)
				}
			}
		}

		binary.LittleEndian.PutUint32(crcValue, crc32(outputBuffer.Bytes()))
	} else {
		compType = uncompressed
		binary.Write(&outputBuffer, binary.LittleEndian, data)
		binary.LittleEndian.PutUint32(crcValue, 0x00000000)
	}

	// Write compressed RTF header
	result := bytes.Buffer{}
	binary.Write(&result, binary.LittleEndian, uint32(len(outputBuffer.Bytes())+12))
	binary.Write(&result, binary.LittleEndian, uint32(len(data)))
	result.Write(compType[:])
	result.Write(crcValue)
	result.Write(outputBuffer.Bytes())

	return result.Bytes()
}

// Find the longest match.
func findLongestMatch(idict []byte, stream *bytes.Reader, writeOffset int) (int, int, int) {
	// Read the first char
	char, err := stream.ReadByte()
	if err != nil {
		return 0, 0, writeOffset
	}

	prevWriteOffset := writeOffset
	dictIndex := 0
	matchLen := 0
	longestMatchLen := 0
	dictOffset := 0

	// Find the first char
	for {
		if idict[dictIndex%maxDictSize] == char {
			matchLen++

			// If found longest match
			if matchLen <= 17 && matchLen > longestMatchLen {
				dictOffset = dictIndex - matchLen + 1

				// Add to dictionary and update longest match
				idict[writeOffset] = char
				writeOffset = (writeOffset + 1) % maxDictSize
				longestMatchLen = matchLen
			}

			// Read the next char
			char, err = stream.ReadByte()
			if err != nil {
				pos, err := stream.Seek(0, io.SeekCurrent)
				if err != nil {
					return 0, 0, writeOffset
				}

				stream.Seek(pos-int64(matchLen), 0)
				return dictOffset, longestMatchLen, writeOffset
			}
		} else {
			pos, err := stream.Seek(0, io.SeekCurrent)
			if err != nil {
				return 0, 0, writeOffset
			}

			stream.Seek(pos-int64(matchLen)-1, 0)
			matchLen = 0

			// Read the first char
			char, err = stream.ReadByte()
			if err != nil {
				break
			}
		}

		dictIndex++
		if dictIndex >= prevWriteOffset+longestMatchLen {
			break
		}
	}

	pos, err := stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, 0, writeOffset
	}

	stream.Seek(pos-int64(matchLen)-1, 0)

	return dictOffset, longestMatchLen, writeOffset
}
