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
	output := cli.GetOutputFile()

	reader := reader.NewReader(path)

	fileContent, err := reader.ReadFile()
	if err != nil {
		fmt.Println("READ ERROR: ", err)
		return
	}

	freqTable := getFrequencyTable(fileContent)
	codeTable := getCodeTable(freqTable)

	writer := getWriter(output, codeTable, "compressed")
	err = writer.WriteFile()
	if err != nil {
		fmt.Println("WRITE ERROR: ", err)
		return
	}

	fmt.Println(codeTable)
}

func getFrequencyTable(text string) map[rune]int {
	newFreqTable := compressor.NewFrequencyTable()
	newFreqTable.Create(text)
	freqTable := newFreqTable.Get()
	return freqTable
}

func getCodeTable(freqTable map[rune]int) string {
	newBinaryTree := compressor.NewBinaryTree(&freqTable)
	newBinaryTree.GetPrefixCodeTable()
	return newBinaryTree.GetCodeTableAsString()
}

func getWriter(headerText, compressedText, output string) *reader.Writer {
	return reader.NewWriter(headerText, compressedText, output)
}
