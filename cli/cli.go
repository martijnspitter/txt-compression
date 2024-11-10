package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd    cobra.Command
	path       string
	outputFile string
}

func NewCLI() *CLI {
	cli := &CLI{}
	cli.rootCmd = cobra.Command{
		Use:   "compress",
		Short: "Supply the file path to compress the file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("File to compress: %s\n", args[0])
			fmt.Printf("File to output to: %s\n", args[1])

			cli.path = args[0]
			cli.outputFile = args[1]
		},
	}
	return cli
}

func (c *CLI) Run() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *CLI) GetPath() string {
	return c.path
}

func (c *CLI) GetOutputFile() string {
	return c.outputFile
}
