package main

import (
	"fmt"
	"github.com/asmcos/requests"
	"github.com/gin-gonic/gin"
	"liveRedirect/service"
	"net/http"
	"strconv"
	"strings"
)

func main() {

	//初始化服务列表
	serviceMap := service.GetServiceMap()

	//启动web服务
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	handleService(r, serviceMap)
	r.Run(":5000")
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
			c.String(200, err.Error())
			return
		}

		if key == "huya" {
			processHuya(c, url, roomId)
			return
		}

		c.Redirect(http.StatusFound, url)
	})
}

func processHuya(c *gin.Context, url string, id string) {
	i := strings.LastIndex(url, "/")
	if i > -1 {
		urlPrefix := url[0 : i+1]
		resp, err := requests.Get(strings.TrimSpace(url), requests.Header{"Content-Type": "application/x-www-form-urlencoded",
			"User-Agent": "Mozilla/5.0 (Linux; Android 5.0; SM-G900P Build/LRX21T) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Mobile Safari/537.36 ",
		})
		c.Header("Content-type", "application/vnd.apple.mpegurl")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization")

		text := resp.Text()
		if err != nil {
			service.IncrHuyaCount(id)

			m3u8_content := "#EXTM3U\n#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1\n" + url
			c.Header("Content-Length", strconv.Itoa(len([]rune(m3u8_content))))
			c.Writer.WriteString(m3u8_content)
			return
		}
		if !strings.Contains(text, ".ts") {
			service.IncrHuyaCount(id)
			m3u8_content := text
			c.Header("Content-Length", strconv.Itoa(len([]rune(m3u8_content))))
			c.Writer.WriteString(m3u8_content)
			return
		}
		ss := strings.Split(text, "\n")
		s := ""
		for _, v := range ss {
			if strings.HasPrefix(v, "#") {
				s = s + v + "\r\n"
			} else {
				s = s + urlPrefix + v + "\r\n"
			}
		}
		c.Header("Content-Length", strconv.Itoa(len([]rune(s))))

		c.Writer.WriteString(s)
		return
	}
}
