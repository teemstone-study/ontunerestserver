package teemcache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// var mapchannel *chan map[string]CacheData
var cachedata map[string]CacheData

func MakeCacheServer() {
	fmt.Println("MakeCacheServer Join")
	cachedata = make(map[string]CacheData)

	//fmt.Println("MakeCacheServer End")
}

func MakeCache(ch *chan string, endch *chan string /*, mapch *chan map[string]CacheData*/) {
	//cachedata = createHandler()
	//mapchannel = mapch
	wg.Add(1)
	go setCache(ch, endch, cachedata)
	clearticker := time.NewTicker(time.Hour * 12)
	go func() {
		for t := range clearticker.C {
			if len(cachedata) > 0 {
				fmt.Println("Cache Map Delete... ", t)
				for k := range cachedata {
					delete(cachedata, k)
				}
			}
		}
	}()
	wg.Wait()
	clearticker.Stop()
}

func GetCachemap() map[string]CacheData {
	return cachedata
}

func setCache(ch *chan string, endch *chan string, cachedatamap map[string]CacheData) {
	for {
		select {
		case data := <-*ch:
			fmt.Println("data input")
			arr := strings.Split(data, "-")
			cachedata := createHandler()
			cachedata.Datatype = arr[0]
			cachedata.Id = arr[1]
			cachedata.Data = arr[2]
			cachedatamap[cachedata.Datatype] = *cachedata
			fmt.Println(cachedatamap[cachedata.Datatype])
		case data := <-*endch:
			fmt.Println("data End")
			if data == "end" {
				return
			}
		}
	}
}
