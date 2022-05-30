package service

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

type HuyaLiveService struct {
}

var (
	mutex sync.Mutex

	// 设置超时时间和清理时间
	ccache       = cache.New(20*time.Second, 21*time.Second)
	HuyaCountMap = make(map[string]int)
)

func IncrHuyaCount(key string) {
	mutex.Lock()
	n := HuyaCountMap[key]
	HuyaCountMap[key] = n + 1
	defer mutex.Unlock()
}

func (s *HuyaLiveService) GetPlayUrl(key string) (string, error) {

	uurls := []string{}
	if urls, ok := ccache.Get(key); ok {
		uurls = urls.([]string)
	} else {
		urls, err := GetHuyaStreamUrls("https://www.huya.com/" + key)
		if err != nil {
			fmt.Print(err.Error())
			return "", err
		}
		HuyaCountMap[key] = 0
		uurls = urls
		ccache.Set(key, urls, 30*time.Second)
	}
	if n, ok := HuyaCountMap[key]; ok {
		if n >= len(uurls) {
			n = 0
			HuyaCountMap[key] = -1
		}
		return uurls[n], nil
	}

	return "", nil
}
func init() {
	RegisterService("huya", new(HuyaLiveService))
}
