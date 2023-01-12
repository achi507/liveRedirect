package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"liveRedirect/service"
	"net/http"
	"strconv"
)

// 实际中应该用更好的变量名
var (
	port int
	help bool
)

func main() {
	flag.BoolVar(&help, "h", false, "help info")
	flag.IntVar(&port, "p", 5000, "listen port")
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	//初始化服务列表
	serviceMap := service.GetServiceMap()

	//启动web服务
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	handleService(r, serviceMap)

	fmt.Println("服务启动监听:", port)
	r.Run(":" + strconv.Itoa(port))
}

func handleService(r *gin.Engine, serviceMap map[string]service.LiveService) gin.IRoutes {
	return r.GET("/:key/:id", func(c *gin.Context) {
		key := c.Param("key")
		roomId := c.Param("id")

		_, ok := serviceMap[key]
		if !ok {
			key = "huya"
		}

		url, err := serviceMap[key].GetPlayUrl(roomId)
		if err != nil {
			fmt.Println(err.Error())
			c.String(500, err.Error())
			return
		}

		c.Redirect(http.StatusFound, url)
	})
}
