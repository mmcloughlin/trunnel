package main

import (
	"os"

	"github.com/mmcloughlin/trunnel/gen"
	"github.com/mmcloughlin/trunnel/parse"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "trunnel"
	app.Usage = "Code generator for binary parsing"
	app.Version = ""

	app.Commands = []cli.Command{
		build,
	}

	app.Run(os.Args)
}

// build command
var (
	build = cli.Command{
		Name:      "build",
		Usage:     "generate go code from trunnel",
		ArgsUsage: "trunnelfile",
		Action: func(c *cli.Context) error {
			filename := c.Args().First()
			if filename == "" {
				return cli.NewExitError("missing trunnel filename", 1)
			}

			f, err := parse.File(filename)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			src, err := gen.File("pkg", f)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			os.Stdout.Write(src)

			return nil
		},
	}
)
