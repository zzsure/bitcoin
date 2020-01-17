package main

import (
	"os"

	"github.com/urfave/cli"
	"bitcoin/cmd/server"
	"bitcoin/cmd/tool"
)

func main() {
	app := cli.NewApp()
	app.Name = "bitcoin"
	app.Commands = []cli.Command{
		server.Server,
		tool.InitDB,
		tool.MigrateDB,
		tool.Sale,
		tool.Huobi,
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
