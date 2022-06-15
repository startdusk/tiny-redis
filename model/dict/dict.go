package dict

type Consumer func(key string, val any) bool

type Dict interface {
	Get(key string) (val any, exists bool)
	Len() int
	Put(key string, val any) (result int)
	PutIfAbsent(key string, val any) (result int)
	PutIfExists(key string, val any) (result int)
	Remove(key string) (result int)
	Range(concumer Consumer)
	Keys() (keys []string)
	RandomKeys(limit int) (keys []string)
	RandomDistinctKeys(limit int) (keys []string)
	Clear() error
}
