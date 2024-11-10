package main

import (
	"fmt"
	"os"
	"txt-compression/cli"
	"txt-compression/compressor"
)

func main() {
	cli := cli.NewCLI()
	cli.Run()
	inputPath := cli.GetPath()
	outputPath := cli.GetOutputFile()
	isDecompress := cli.IsDecompress()

	var err error
	if isDecompress {
		err = decompressFile(inputPath, outputPath)
	} else {
		err = compressFile(inputPath, outputPath)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func compressFile(inputPath, outputPath string) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Create new compressor
	comp := compressor.NewCompressor()

	// First pass: build frequency table
	fmt.Println("Building frequency table...")
	if err := comp.BuildFrequencyTable(inputFile); err != nil {
		return fmt.Errorf("failed to build frequency table: %w", err)
	}

	// Reset file pointer to beginning for second pass
	if _, err := inputFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file position: %w", err)
	}

	// Second pass: compress and write to output file
	fmt.Println("Compressing file...")
	if err := comp.Compress(inputFile, outputFile); err != nil {
		return fmt.Errorf("failed to compress file: %w", err)
	}

	// Print compression statistics
	inputInfo, _ := inputFile.Stat()
	outputInfo, _ := outputFile.Stat()
	ratio := float64(outputInfo.Size()) / float64(inputInfo.Size()) * 100

	fmt.Printf("Compression complete:\n")
	fmt.Printf("Original size: %d bytes\n", inputInfo.Size())
	fmt.Printf("Compressed size: %d bytes\n", outputInfo.Size())
	fmt.Printf("Compression ratio: %.2f%%\n", ratio)

	return nil
}

func decompressFile(inputPath, outputPath string) error {
	// Open input file
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// Write UTF-8 BOM
	bom := []byte{0xEF, 0xBB, 0xBF}
	if _, err := outputFile.Write(bom); err != nil {
		return fmt.Errorf("failed to write BOM: %w", err)
	}

	// Create compressor and decompress
	comp := compressor.NewCompressor()
	if err := comp.ReadCompressedFile(inputFile, outputFile); err != nil {
		return fmt.Errorf("decompression failed: %w", err)
	}

	// Print decompression statistics
	inputInfo, _ := inputFile.Stat()
	outputInfo, _ := outputFile.Stat()

	fmt.Printf("Decompression complete:\n")
	fmt.Printf("Compressed size: %d bytes\n", inputInfo.Size())
	fmt.Printf("Decompressed size: %d bytes\n", outputInfo.Size())

	return nil
}
