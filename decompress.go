package crtf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Decompress data using RTF compression algorithm.
func Decompress(data []byte) ([]byte, error) {
	if len(data) < 16 {
		return nil, errors.New("data must be at least 16 bytes long")
	}

	idict := make([]byte, maxDictSize)
	copy(idict, initDict)
	for i := initDictSize; i < maxDictSize; i++ {
		idict[i] = ' '
	}

	writeOffset := initDictSize
	var outputBuffer bytes.Buffer
	inStream := bytes.NewReader(data)

	// Read compressed RTF header
	var compSize, rawSize uint32
	var compType [4]byte
	var crcValue uint32

	if err := binary.Read(inStream, binary.LittleEndian, &compSize); err != nil {
		return nil, err
	}

	if err := binary.Read(inStream, binary.LittleEndian, &rawSize); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(inStream, compType[:]); err != nil {
		return nil, err
	}

	if err := binary.Read(inStream, binary.LittleEndian, &crcValue); err != nil {
		return nil, err
	}

	// Get only data
	contents := make([]byte, compSize-12)
	if _, err := io.ReadFull(inStream, contents); err != nil {
		return nil, err
	}

	if compType == compressed {
		// Check CRC
		if crcValue != crc32(contents) {
			return nil, errors.New("CRC is invalid! The file is corrupt")
		}

		contentsReader := bytes.NewReader(contents)
		for {
			controlByte, err := contentsReader.ReadByte()
			if err == io.EOF {
				break
			}

			if err != nil {
				return nil, err
			}

			// Check bits from LSB to MSB
			for i := 0; i < 8; i++ {
				if controlByte&(1<<uint(i)) != 0 {
					// Token is reference (16 bit)
					var token uint16
					if err := binary.Read(contentsReader, binary.BigEndian, &token); err != nil {
						if err == io.EOF {
							return outputBuffer.Bytes(), nil
						}

						return nil, err
					}

					// Extract [12 bit offset][4 bit length]
					offset := (token >> 4) & 0xFFF
					length := (token & 0xF) + 2

					// End indicator
					if int(offset) == writeOffset {
						return outputBuffer.Bytes(), nil
					}

					for j := 0; j < int(length); j++ {
						char := idict[(offset+uint16(j))%maxDictSize]
						outputBuffer.WriteByte(char)
						idict[writeOffset] = char
						writeOffset = (writeOffset + 1) % maxDictSize
					}
				} else {
					// Token is literal (8 bit)
					char, err := contentsReader.ReadByte()
					if err == io.EOF {
						return outputBuffer.Bytes(), nil
					}

					if err != nil {
						return nil, err
					}

					outputBuffer.WriteByte(char)
					idict[writeOffset] = char
					writeOffset = (writeOffset + 1) % maxDictSize
				}
			}
		}
	} else if compType == uncompressed {
		// Check CRC
		if crcValue != 0x00000000 {
			return nil, errors.New("CRC is invalid! Must be 0x00000000")
		}

		return contents[:rawSize], nil
	} else {
		return nil, errors.New("unknown type of RTF compression")
	}

	return outputBuffer.Bytes(), nil
}
