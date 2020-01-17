package tool

import (
	"github.com/urfave/cli"
	"bitcoin/conf"
	"bitcoin/library/db"
	"bitcoin/models"
)

var MigrateDB = cli.Command{
	Name:  "migrate",
	Usage: "migrate db",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "./config.toml",
			Usage: "toml配置文件",
		},
		cli.StringFlag{
			Name:  "args, a",
			Value: "",
			Usage: "multiconfig cmd line args",
		},
	},
	Action: runMigrateDB,
}

func runMigrateDB(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	db.Init()
	models.MigrateTable()
}
