package main

import (
    "github.com/urfave/cli"
    "gitlab.azbit.cn/web/bitcoin/cmd/server"
    "gitlab.azbit.cn/web/bitcoin/cmd/tool"
    "os"
)

func main() {
    app := cli.NewApp()
    app.Name = "bitcoin"
    app.Commands = []cli.Command{
        server.Server,
        tool.InitDB,
    }
    err := app.Run(os.Args)
    if err != nil {
        panic(err)
    }
}
