package main

import (
	"fmt"
	"txt-compression/cli"
	"txt-compression/compressor"
	"txt-compression/reader"
)

func main() {
	cli := cli.NewCLI()
	cli.Run()
	path := cli.GetPath()

	reader := reader.NewReader(path)

	fileContent, err := reader.ReadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	freqTable := getFrequencyTable(fileContent)
	// binaryTree := getBinaryTree(freqTable)

	// compressor := compressor.NewCompressor(binaryTree)

	fmt.Println(freqTable)
}

func getFrequencyTable(text string) map[rune]int {
	newFreqTable := compressor.NewFrequencyTable()
	newFreqTable.Create(text)
	freqTable := newFreqTable.Get()
	return freqTable
}

func getCodeTable(freqTable map[rune]int) *compressor.CodeTable {
	newBinaryTree := compressor.NewBinaryTree(&freqTable)
	return newBinaryTree.GetPrefixCodeTable()
}
