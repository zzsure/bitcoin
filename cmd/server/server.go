package server

import (
	"bitcoin/library/redis"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli"
	"bitcoin/conf"
	"bitcoin/controller/v1"
	"bitcoin/library/db"
	"bitcoin/library/log"
	"bitcoin/middleware"
	//"bitcoin/modules/socket"
	"bitcoin/modules/cron"
	"bitcoin/modules/strategy"
)

var Server = cli.Command{
	Name:  "server",
	Usage: "bitcoin http server",
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
	Action: run,
}

func run(c *cli.Context) {
	conf.Init(c.String("conf"), c.String("args"))
	log.Init()
	db.Init()
	redis.Init()
	strategy.Init()
	//socket.Init()
	cron.Init()

	GinEngine().Run(conf.Config.Server.Listen)
}

func GinEngine() *gin.Engine {
	var r *gin.Engine
	if conf.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(middleware.Recovery)
	} else {
		r = gin.Default()
	}
	r.Use(middleware.Access)
	r.Use(middleware.Auth)
	r.GET("/health")
	V1(r)

	return r
}

func V1(r *gin.Engine) {
	g := r.Group("/v1")
	{
		g.POST("/echo", v1.Echo)
	}
}
