package cmd

import (
	"fmt"
	"strings"
	"os"
	"io"
	"context"
	"github.com/spf13/cobra"
	"github.com/peterh/liner"
)

const exitPromptMessage = "Do you want to exit ([y]/n)? "

type rhubCli struct {
	output io.Writer
	prompt string
	banner string
}

func init() {
	runCommand := &cobra.Command{
		Use:   "cli",
		Short: "Start the rhub cli",
		Long:  "Start the rhub cli",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()
			cli := &rhubCli{
				output : os.Stdout,
				prompt : "> ",
				banner : "Rhub CLI",
			}
			cli.start(ctx)
		},
	}

	RootCmd.AddCommand(runCommand)
}

func (c *rhubCli) start(ctx context.Context) {
	// Initialize the liner library.

	line := liner.NewLiner()
	defer line.Close()
	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)

	if len(c.banner) > 0 {
		fmt.Fprintln(c.output, c.banner)
	}

loop:
	for true {

		input, err := line.Prompt(c.prompt)

		// prompt on ctrl+d
		if err == io.EOF {
			goto exitPrompt
		}

		// reset on ctrl+c
		if err == liner.ErrPromptAborted {
			continue
		}

		// exit on unknown error
		if err != nil {
			fmt.Fprintln(c.output, "error (fatal):", err)
			os.Exit(1)
		}

		c.processCmd(input)
		line.AppendHistory(input)
	}
exitPrompt:
	fmt.Fprintln(c.output)

	for true {
		input, err := line.Prompt(exitPromptMessage)

		// exit on ctrl+d
		if err == io.EOF {
			break
		}

		// reset on ctrl+c
		if err == liner.ErrPromptAborted {
			goto loop
		}

		// exit on unknown error
		if err != nil {
			fmt.Fprintln(c.output, "error (fatal):", err)
			os.Exit(1)
		}

		switch strings.ToLower(input) {
		case "", "y", "yes":
			goto exit
		case "n", "no":
			goto loop
		}
	}

exit:
}

func (c *rhubCli) processCmd(line string)  {
	fmt.Fprintln(c.output, line, "is not supported")

}