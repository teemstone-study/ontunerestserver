package teemcache

type CacheData struct {
	Datatype string
	Id       string
	Data     string
}

func createHandler() *CacheData {
	cachereturn := &CacheData{}
	return cachereturn
}
