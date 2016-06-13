package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"clog"
	"config"
	"orm"
	"pongo"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/engine/standard"
)

var (
	configPath string
	rootPath   string
)

const (
	CONFIG_PATH = "src/config/config.toml"
	INDEX_PATH  = "src/views/index.html"
)

func main() {
	config.ParseConfig(configPath)

	clog.InitLog(config.Conf.Log)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.File(INDEX_PATH)
	})

	// 设置render
	render := pongo.GetRenderer(pongo.PongorOption{Directory: config.Conf.Pongo.Directory,
		Reload: config.Conf.Pongo.Reload})
	e.SetRenderer(render)

	// 设置错误处理函数
	e.SetHTTPErrorHandler(func(err error, c echo.Context) {
		c.Render(http.StatusOK, "error.html", map[string]interface{}{"msg": err.Error()})
	})

	if err := orm.Open(config.Conf.DB); err != nil {
		clog.Error("打开数据库失败, %s", err)
		return
	}

	orm.CreateTable()

	// 注册路由
	router(e)

	var server engine.Server

	if config.Conf.Web.Fasthttp {
		server = fasthttp.New(config.Conf.Web.Listen)
	} else {
		server = standard.New(config.Conf.Web.Listen)
	}

	clog.Info("服务器监听的地址为%s", config.Conf.Web.Listen)
	e.Run(server)
}

func init() {
	rootPath, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("获取文件路径失败, %s\n", err)
		os.Exit(1)
	}

	flag.StringVar(&configPath, "-c", filepath.Join(rootPath, CONFIG_PATH), "-c /configPath")
}
