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

	compressor := compressor.NewCompressor(fileContent)
	compressor.GenerateFreqTable()

	freqTable := compressor.GetHumanReadableFreqTable()

	fmt.Println(freqTable)
}
