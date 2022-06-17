package dict

import (
	"sync"
)

type SyncDict struct {
	m sync.Map
}

func NewSyncDict() *SyncDict {
	return &SyncDict{}
}

func (d *SyncDict) Get(key string) (val any, exists bool) {
	return d.m.Load(key)
}

func (d *SyncDict) Len() int {
	var length int
	d.m.Range(func(_, _ any) bool {
		length++
		return true
	})
	return length
}

func (d *SyncDict) Put(key string, val any) (result int) {
	_, existed := d.m.Load(key)
	d.m.Store(key, val)
	if existed {
		// update
		return 0
	}
	// new
	return 1
}

func (d *SyncDict) PutIfAbsent(key string, val any) (result int) {
	_, existed := d.m.Load(key)
	if existed {
		return 0
	}
	// not exists
	d.m.Store(key, val)
	return 1
}

func (d *SyncDict) PutIfExists(key string, val any) (result int) {
	_, existed := d.m.Load(key)
	if existed {
		d.m.Store(key, val)
		return 1
	}
	// not exists
	return 0
}

func (d *SyncDict) Remove(key string) (result int) {
	_, existed := d.m.Load(key)
	d.m.Delete(key)
	if existed {
		return 1
	}
	return 0
}

func (d *SyncDict) Range(concumer Consumer) {
	d.m.Range(func(key, value any) bool {
		concumer(key.(string), value)
		return true
	})
}

func (d *SyncDict) Keys() (keys []string) {
	keys = make([]string, d.Len())
	var i int
	d.m.Range(func(key, _ any) bool {
		keys[i] = key.(string)
		i++
		return true
	})
	return
}

func (d *SyncDict) RandomKeys(limit int) (keys []string) {
	if limit <= 0 {
		return
	}
	keys = make([]string, limit)
	for i := 0; i < limit; i++ {
		d.m.Range(func(key, value any) bool {
			keys[i] = key.(string)
			return false
		})
	}
	return keys
}

func (d *SyncDict) RandomDistinctKeys(limit int) (keys []string) {
	if limit <= 0 {
		return
	}
	keys = make([]string, limit)
	var i int
	d.m.Range(func(key, value any) bool {
		keys[i] = key.(string)
		i++
		return !(i == limit)
	})
	return keys
}

func (d *SyncDict) Clear() error {
	*d = *NewSyncDict()
	return nil
}
