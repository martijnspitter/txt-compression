package compressor

import (
	"bufio"
	"fmt"
	"io"
)

type Decoder struct {
	decodeTable map[string]byte
	currentCode string
}

func NewDecoder(decodeTable map[string]byte) *Decoder {
	return &Decoder{
		decodeTable: decodeTable,
		currentCode: "",
	}
}

func (d *Decoder) DecodeStream(reader io.Reader, writer io.Writer, paddingBits uint8) error {
	bufferedWriter := bufio.NewWriter(writer)
	defer bufferedWriter.Flush()
	buf := make([]byte, 8192)
	isLastChunk := false

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading compressed data: %w", err)
		}

		// Check if this is the last chunk
		if n < len(buf) {
			isLastChunk = true
		}

		// Process each byte
		for i := 0; i < n; i++ {
			// Skip padding bits in the last byte of the last chunk
			numBits := 8
			if isLastChunk && i == n-1 && paddingBits > 0 {
				numBits = 8 - int(paddingBits)
			}

			// Process each bit
			for bit := 7; bit >= 8-numBits; bit-- {
				if buf[i]&(1<<bit) != 0 {
					d.currentCode += "1"
				} else {
					d.currentCode += "0"
				}

				// Check if current code matches any in the decode table
				if char, ok := d.decodeTable[d.currentCode]; ok {
					if err := d.writeRune(writer, char); err != nil {
						return fmt.Errorf("error writing decoded character: %w", err)
					}
					d.currentCode = ""
				}
			}
		}
	}

	if d.currentCode != "" {
		return fmt.Errorf("invalid compressed data: incomplete code at end")
	}

	return nil
}

func (d *Decoder) writeRune(writer io.Writer, b byte) error {
	_, err := writer.Write([]byte{b})
	return err
}
