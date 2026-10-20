package dict

import (
	"fmt"
	"sync"
)

var (
	lock sync.RWMutex
	once sync.Once
	dict map[string]Source
)

func Add(name string, e any) {
	if d, ok := getDict(name); ok {
		d.Add(e)
	}
}

func Update(name string, e any) {
	if d, ok := getDict(name); ok {
		d.Update(e)
	}
}

func Remove(name string, e any) {
	if d, ok := getDict(name); ok {
		d.Remove(e)
	}
}

func ExistKey(name, category, key string) (OK bool) {
	if d, ok := getDict(name); ok {
		OK = d.ExistKey(category, key)
	}
	return
}

func AllKVs(name string) map[string]map[string]string {
	if d, ok := getDict(name); ok {
		return d.AllKVs()
	}
	return emptySource.AllKVs()
}

func AllReversedKVs(name string) map[string]map[string]string {
	if d, ok := getDict(name); ok {
		return d.AllReversedKVs()
	}
	return emptySource.AllReversedKVs()
}

func GetKV(name string, category string) map[string]string {
	if d, ok := getDict(name); ok {
		return d.GetKV(category)
	}
	return emptySource.GetKV(category)
}

func GetReversedKV(name string, category string) map[string]string {
	if d, ok := getDict(name); ok {
		return d.GetReversedKV(category)
	}
	return emptySource.GetReversedKV(category)
}

func Get(name string, key any) any {
	if d, ok := getDict(name); ok {
		return d.Get(key)
	}
	return emptySource.Get(key)
}

func AllDictItems(name string) map[string]*Items {
	if d, ok := getDict(name); ok {
		return d.AllDictItems()
	}
	return emptySource.AllDictItems()
}

func DictItems(name, category string) *Items {
	if d, ok := getDict(name); ok {
		return d.DictItems(category)
	}
	return emptySource.DictItems(category)
}

func DictItem(name, category, key string) Item {
	if d, ok := getDict(name); ok {
		return d.DictItem(category, key)
	}
	return emptySource.DictItem(category, key)
}

func AddSource(source Source, config *SourceConfig) (err error) {
	if err = source.Load(config); err == nil {
		lock.Lock()
		defer lock.Unlock()
		if _, OK := dict[source.Name()]; OK {
			err = fmt.Errorf("duplicate source %s", source.Name())
		} else {
			dict[source.Name()] = source
		}
	}
	return
}

func Init(config *Config, sources ...Source) {
	lock.Lock()
	defer lock.Unlock()
	once.Do(func() {
		if sourceMap, err := reload(config, sources); err == nil {
			dict = sourceMap
		} else {
			panic(err)
		}
	})
}

func Reload(config *Config, sources ...Source) error {
	sourceMap, err := reload(config, sources)
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		dict = sourceMap
	}
	return err
}

func reload(config *Config, sources []Source) (sourceMap map[string]Source, err error) {
	sourceMap = make(map[string]Source)
	for _, source := range sources {
		if err = source.Load(&SourceConfig{DBName: config.DBName}); err != nil {
			break
		}
		sourceMap[source.Name()] = source
	}
	return
}

func getDict(name string) (s Source, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	s, ok = dict[name]
	return
}

func SourceExists(name string) (ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	_, ok = dict[name]
	return
}

var emptySource EmptySource
