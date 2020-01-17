package tool

import (
	"github.com/urfave/cli"
	"bitcoin/conf"
	"bitcoin/library/db"
	"bitcoin/models"
	//"bitcoin/library/log"
)

var InitDB = cli.Command{
	Name:  "init",
	Usage: "bitcoin init db",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "conf, c",
			Value: "config.toml",
			Usage: "toml配置文件",
		},
		cli.StringFlag{
			Name:  "args",
			Value: "",
			Usage: "multiconfig cmd line args",
		},
	},
	Action: runInitDB,
}

func runInitDB(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	//log.Init()
	db.Init()
	models.CreateTable()
}
