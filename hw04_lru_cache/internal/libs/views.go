package libs

type CacheItem struct {
	Value interface{}
	Key   Key
}

func NewCacheItem(value interface{}, key Key) *CacheItem {
	return &CacheItem{
		Value: value,
		Key:   key,
	}
}
