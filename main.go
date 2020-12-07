package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"webssh/controller"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
)

var (
	port     = flag.Int("p", 5032, "服务运行端口")
	timeout  int
	savePass bool
)

func init() {
	flag.IntVar(&timeout, "t", 60, "ssh连接超时时间(min)")
	flag.BoolVar(&savePass, "s", true, "是否保存ssh密码")
	flag.Parse()
}

func staticRouter(router *gin.Engine) {
	box := packr.New("websshBox", "./web/dist")
	router.Use(func(c *gin.Context) {
		requestURL := c.Request.URL.Path
		if box.Has(requestURL) || requestURL == "/" {
			http.FileServer(box).ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
}

func main() {
	server := gin.Default()
	server.Use(gzip.Gzip(gzip.DefaultCompression))
	staticRouter(server)
	server.GET("/term", func(c *gin.Context) {
		controller.TermWs(c, time.Duration(timeout)*time.Minute)
	})
	server.GET("/check", func(c *gin.Context) {
		envVal, ok := os.LookupEnv("savePass")
		if ok {
			b, err := strconv.ParseBool(envVal)
			if err != nil {
				savePass = false
			} else {
				savePass = b
			}
		}
		c.JSON(200, map[string]interface{}{
			"savePass": savePass,
			"result":   controller.CheckSSH(c),
		})
	})
	file := server.Group("/file")
	{
		file.GET("/list", func(c *gin.Context) {
			c.JSON(200, controller.FileList(c))
		})
		file.GET("/download", func(c *gin.Context) {
			controller.DownloadFile(c)
		})
		file.POST("/upload", func(c *gin.Context) {
			c.JSON(200, controller.UploadFile(c))
		})
	}
	server.Run(fmt.Sprintf(":%d", *port))
}
