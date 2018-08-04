package main

import (
	"flag"

	"github.com/lightlfyan/buffgate/config"
	"github.com/lightlfyan/buffgate/giant"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"time"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
		c.Next()
	}
}

func handler(c *gin.Context) {
	glog.Info(c.Request.RemoteAddr)
	event := giant.GetEvent()

	if err := c.BindJSON(event); err != nil {
		glog.Info(err)
		c.String(400, "payload error\n")
		return
	}

	event.Timestamp = time.Now().UTC().Add(time.Hour * 8)
	giant.Send(event)
	c.String(200, "ok")
}

func main() {
	flag.Parse()

	glog.MaxSize = 1024 * 1024 * 1024
	go giant.Start()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(Cors())

	r.POST("/collect", handler)
	r.Run(config.Config.Port)
}
