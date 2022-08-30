package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func GenerateMkdownCmd() cli.Command {
	return cli.Command{
		Name:  "doc",
		Usage: "generate markdown documentation",
		Action: func(c *cli.Context) error {
			fmt.Println(c.App.ToMarkdown())

			return nil
		},
	}
}
