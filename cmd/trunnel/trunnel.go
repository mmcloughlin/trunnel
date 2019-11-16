package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/mmcloughlin/trunnel/gen"
	"github.com/mmcloughlin/trunnel/meta"
	"github.com/mmcloughlin/trunnel/parse"
)

func main() {
	app := cli.NewApp()
	app.Name = "trunnel"
	app.Usage = "Code generator for binary parsing"
	app.Version = meta.GitSHA

	app.Commands = []*cli.Command{
		build,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// build command
var (
	cfg gen.Config

	build = &cli.Command{
		Name:      "build",
		Usage:     "Generate go package from trunnel",
		ArgsUsage: "<trunnelfile>...",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "pkg",
				Aliases:     []string{"p"},
				Usage:       "package name",
				Destination: &cfg.Package,
			},
			&cli.StringFlag{
				Name:        "dir",
				Aliases:     []string{"d"},
				Usage:       "output directory",
				Value:       ".",
				Destination: &cfg.Dir,
			},
		},
		Action: func(c *cli.Context) error {
			filenames := c.Args().Slice()
			if len(filenames) == 0 {
				return cli.NewExitError("missing trunnel filenames", 1)
			}

			fs, err := parse.Files(filenames)
			if err != nil {
				return cli.NewExitError(err, 1)
			}

			if err = gen.Package(cfg, fs); err != nil {
				return cli.NewExitError(err, 1)
			}

			return nil
		},
	}
)
