package main

import (
	"fmt"
	"txt-compression/cli"
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

	fmt.Println(fileContent)
}
