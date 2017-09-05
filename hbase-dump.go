package main

import (
	"os"

	"github.com/urfave/cli"
	"hbase-dump/cmd"
)

func makeApp() *cli.App {
	app := cli.NewApp()

	app.Name = ""
	app.Usage = "Dump file (JSON, CSV, TSV) from hbase table. Edit"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:    "dump",
			Aliases: []string{"d"},
			Usage:   "dump file",
			Action:  cmd.CmdDump,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "access, a",
					Usage: "Connect to access point",
				},
				cli.StringFlag{
					Name: "index, i",
					Usage: "Dump index",
				},
				cli.StringFlag{
					Name: "table, t",
					Usage: "Dump table",
				},
				cli.BoolFlag{
					Name: "child, c",
					Usage: "Table is child",
					Hidden: false,
				},
			},
		},
	}

	return app
}

func main() {
	app := makeApp()
	app.Run(os.Args)
}
