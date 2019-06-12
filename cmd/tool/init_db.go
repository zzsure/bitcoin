package tool

import (
	"github.com/urfave/cli"
	"gitlab.azbit.cn/web/bitcoin/conf"
	"gitlab.azbit.cn/web/bitcoin/library/db"
	"gitlab.azbit.cn/web/bitcoin/library/log"
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
	log.Init()
	db.Init()

	// 初始化db
	db.DB.LogMode(true)
	// drop table
	db.DB.DropTableIfExists()

	// create table
	create := db.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	create.CreateTable()
}

