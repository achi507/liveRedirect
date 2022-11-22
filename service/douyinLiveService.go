package service

import (
	"errors"
	"github.com/asmcos/requests"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"regexp"
)

type DouyuYinService struct {
}

func init() {
	RegisterService("douyin", new(DouyuYinService))
}
func (s *DouyuYinService) GetPlayUrl(key string) (string, error) {
	return parseUrl(key)
}
func parseUrl(room string) (string, error) {
	u := room
	url := "https://v.douyin.com/" + u + "/"
	id, _ := parseRoomId(url)

	roomUrl := "https://webcast.amemv.com/webcast/room/reflow/info/"
	p := requests.Params{
		"type_id": "0",
		"live_id": "1",
		"room_id": id,
		"app_id":  "1128",
		"X-Bogus": "1",
	}
	resp, _ := requests.Get(roomUrl, p, requests.Header{"authority": "webcast.amemv.com",
		"cookie": "_tea_utm_cache_1128={%22utm_source%22:%22copy%22%2C%22utm_medium%22:%22android%22%2C%22utm_campaign%22:%22client_share%22}",
	})
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	durl := json.Get([]byte(resp.Text()), "data", "room", "stream_url", "rtmp_pull_url")

	url = durl.ToString()
	if url == "" {
		return "", errors.New("房间不存在")
	}
	return url, nil
}
func parseRoomId(url string) (string, error) {
	//获取roomId
	req, _ := http.NewRequest("HEAD", url, nil)
	// 比如说设置个token
	req.Header.Set("authority", "v.douyin.com")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_1 like Mac OS X) AppleWebKit/603.1.30 (KHTML, like Gecko) Version/10.0 Mobile/14E304 Safari/602.1")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	url = resp.Header.Get("Location")
	re := regexp.MustCompile(`\d{19}`)
	res := re.FindStringSubmatch(url)
	if len(res) > 0 {
		return res[0], nil
	} else {
		return "", errors.New("没用可以用的roomId")
	}
}
