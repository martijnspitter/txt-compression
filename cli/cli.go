package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type CLI struct {
	rootCmd cobra.Command
	path    string
}

func NewCLI() *CLI {
	cli := &CLI{}
	cli.rootCmd = cobra.Command{
		Use:   "askme",
		Short: "Ask me anything",
		Run: func(cmd *cobra.Command, args []string) {
			cli.path = getFilePathFromUser()
			fmt.Printf("You entered file path: %s\n", cli.path)
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

func getFilePathFromUser() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter the file path: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return ""
	}

	// Clean the input by removing newline and spaces
	return strings.TrimSpace(input)
}

func (c *CLI) GetPath() string {
	return c.path
}
