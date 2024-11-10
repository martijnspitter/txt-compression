package compressor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

type Compressor struct {
	BinaryTree *BinaryTree
	FreqTable  *FrequencyTable
	bufferSize int
}

type HeaderEntry struct {
	Char     rune
	CodeLen  uint8
	CodeBits uint32 // The actual bits of the code
}

func NewCompressor() *Compressor {
	return &Compressor{
		FreqTable:  NewFrequencyTable(),
		bufferSize: 8192,
	}
}

func (c *Compressor) BuildFrequencyTable(reader io.Reader) error {
	c.FreqTable.Create(reader)
	freqMap := c.FreqTable.Get()
	c.BinaryTree = NewBinaryTree(&freqMap)
	c.BinaryTree.GetPrefixCodeTable()

	return nil
}

func (c *Compressor) createHeader() ([]byte, error) {
	var entries []HeaderEntry

	// Convert code table to binary format
	for char, codeStr := range c.BinaryTree.CodeTable {
		codeBits := uint32(0)
		for _, bit := range codeStr {
			codeBits = (codeBits << 1)
			if bit == '1' {
				codeBits |= 1
			}
		}

		entries = append(entries, HeaderEntry{
			Char:     char,
			CodeLen:  uint8(len(codeStr)),
			CodeBits: codeBits,
		})
	}

	// Write to buffer
	buf := make([]byte, 0, len(entries)*8) // Approximate size
	buffer := bytes.NewBuffer(buf)

	// Write number of entries
	if err := binary.Write(buffer, binary.BigEndian, uint32(len(entries))); err != nil {
		return nil, err
	}

	// Write each entry
	for _, entry := range entries {
		if err := binary.Write(buffer, binary.BigEndian, int32(entry.Char)); err != nil {
			return nil, err
		}
		if err := binary.Write(buffer, binary.BigEndian, entry.CodeLen); err != nil {
			return nil, err
		}
		if err := binary.Write(buffer, binary.BigEndian, entry.CodeBits); err != nil {
			return nil, err
		}
	}

	return buffer.Bytes(), nil
}

func (c *Compressor) Compress(reader io.Reader, writer io.Writer) error {
	// First create and write the header
	headerBytes, err := c.createHeader()
	if err != nil {
		return fmt.Errorf("failed to create header: %w", err)
	}

	// Write header length (4 bytes)
	if err := binary.Write(writer, binary.BigEndian, uint32(len(headerBytes))); err != nil {
		return fmt.Errorf("failed to write header length: %w", err)
	}

	// Write header
	if _, err := writer.Write(headerBytes); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Reserve space for padding bits count (we'll write it at the end)
	paddingBitsPos := 4 + len(headerBytes)
	if _, err := writer.Write([]byte{0}); err != nil {
		return fmt.Errorf("failed to write padding placeholder: %w", err)
	}

	// Compress the data
	bitBuffer := NewBitBuffer()
	buf := make([]byte, c.bufferSize)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		// Process the buffer as UTF-8 text
		str := string(buf[:n])
		for len(str) > 0 {
			r, size := utf8.DecodeRuneInString(str)
			if r == utf8.RuneError {
				return fmt.Errorf("invalid UTF-8 sequence")
			}

			code, exists := c.BinaryTree.CodeTable[r]
			if !exists {
				return fmt.Errorf("character not found in code table: %c (hex: %X)", r, r)
			}
			bitBuffer.WriteCode(code)

			// Flush complete bytes
			if bytes := bitBuffer.FlushCompleteBytes(); len(bytes) > 0 {
				if _, err := writer.Write(bytes); err != nil {
					return fmt.Errorf("error writing compressed data: %w", err)
				}
			}

			str = str[size:]
		}
	}

	// Write final bits and get padding
	finalBytes, paddingBits := bitBuffer.FlushFinal()
	if _, err := writer.Write(finalBytes); err != nil {
		return fmt.Errorf("error writing final bytes: %w", err)
	}

	// Write the padding bits count if the writer supports seeking
	if seeker, ok := writer.(io.WriteSeeker); ok {
		if _, err := seeker.Seek(int64(paddingBitsPos), io.SeekStart); err != nil {
			return fmt.Errorf("failed to seek to padding position: %w", err)
		}
		if err := binary.Write(seeker, binary.BigEndian, uint8(paddingBits)); err != nil {
			return fmt.Errorf("failed to write padding bits count: %w", err)
		}
	}

	return nil
}

type BitBuffer struct {
	currentByte byte
	bitCount    int
	bytes       []byte
}

func NewBitBuffer() *BitBuffer {
	return &BitBuffer{
		bytes: make([]byte, 0, 1024),
	}
}

func (bb *BitBuffer) WriteCode(code string) {
	for _, bit := range code {
		bb.currentByte = (bb.currentByte << 1)
		if bit == '1' {
			bb.currentByte |= 1
		}
		bb.bitCount++

		if bb.bitCount == 8 {
			bb.bytes = append(bb.bytes, bb.currentByte)
			bb.currentByte = 0
			bb.bitCount = 0
		}
	}
}

func (bb *BitBuffer) FlushCompleteBytes() []byte {
	if len(bb.bytes) == 0 {
		return nil
	}
	bytes := bb.bytes
	bb.bytes = make([]byte, 0, 1024)
	return bytes
}

func (bb *BitBuffer) FlushFinal() ([]byte, uint8) {
	if bb.bitCount > 0 {
		paddingBits := 8 - bb.bitCount
		bb.currentByte = bb.currentByte << paddingBits
		bb.bytes = append(bb.bytes, bb.currentByte)
		return bb.FlushCompleteBytes(), uint8(paddingBits)
	}
	return bb.FlushCompleteBytes(), 0
}

func (c *Compressor) ReadCompressedFile(reader io.Reader, writer io.Writer) error {
	// Read header length
	var headerLength uint32
	if err := binary.Read(reader, binary.BigEndian, &headerLength); err != nil {
		return fmt.Errorf("failed to read header length: %w", err)
	}

	// Read header
	headerBytes := make([]byte, headerLength)
	if _, err := io.ReadFull(reader, headerBytes); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Parse header and build decode table
	decodeTable, err := c.parseHeader(headerBytes)
	if err != nil {
		return fmt.Errorf("failed to parse header: %w", err)
	}

	// Read padding bits
	var paddingBits uint8
	if err := binary.Read(reader, binary.BigEndian, &paddingBits); err != nil {
		return fmt.Errorf("failed to read padding bits: %w", err)
	}

	// Create decoder and process compressed data
	decoder := NewDecoder(decodeTable)
	return decoder.DecodeStream(reader, writer, paddingBits)
}

func (c *Compressor) parseHeader(headerBytes []byte) (map[string]rune, error) {
	buf := bytes.NewReader(headerBytes)
	decodeTable := make(map[string]rune)

	// Read number of entries
	var numEntries uint32
	if err := binary.Read(buf, binary.BigEndian, &numEntries); err != nil {
		return nil, fmt.Errorf("failed to read number of entries: %w", err)
	}

	// Read each entry
	for i := uint32(0); i < numEntries; i++ {
		var char int32
		var codeLen uint8
		var codeBits uint32

		if err := binary.Read(buf, binary.BigEndian, &char); err != nil {
			return nil, fmt.Errorf("failed to read character: %w", err)
		}
		if err := binary.Read(buf, binary.BigEndian, &codeLen); err != nil {
			return nil, fmt.Errorf("failed to read code length: %w", err)
		}
		if err := binary.Read(buf, binary.BigEndian, &codeBits); err != nil {
			return nil, fmt.Errorf("failed to read code bits: %w", err)
		}

		// Convert bits to string code
		code := ""
		for i := uint8(0); i < codeLen; i++ {
			if codeBits&(1<<(codeLen-1-i)) != 0 {
				code += "1"
			} else {
				code += "0"
			}
		}

		decodeTable[code] = rune(char)
	}

	return decodeTable, nil
}
