package rdb

import (
	"errors"
	"sync"
)

var (
	instanceMap map[string]*Instance
	lock        sync.RWMutex
	once        sync.Once
)

func Init(configs map[string]*Config) {
	once.Do(func() {
		lock.Lock()
		defer lock.Unlock()
		if instances, err := load(configs); err == nil {
			instanceMap = instances
		} else {
			panic(err)
		}
	})
}

func Reload(configs map[string]*Config) error {
	if instanceCount() == 0 {
		return uninitializedErr
	}
	instances, err := load(configs)
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		instanceMap = instances
	}
	return err
}

func load(configs map[string]*Config) (instanceMap map[string]*Instance, err error) {
	instanceMap = make(map[string]*Instance)
	for name, config := range configs {
		var ins *Instance
		if ins, err = NewInstance(name, config); err != nil {
			break
		}
		instanceMap[name] = ins
	}
	return
}

func GetConnection(name string) (ins *Instance, ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	ins, ok = instanceMap[name]
	return
}

func GetDBName(name string) (dbName string) {
	if ins, ok := GetConnection(name); ok {
		return ins.DBName()
	}
	return
}

func ValidName(name string) (ok bool) {
	lock.RLock()
	defer lock.RUnlock()
	_, ok = instanceMap[name]
	return
}

func instanceCount() (count int) {
	lock.RLock()
	defer lock.RUnlock()
	if instanceMap != nil {
		count = len(instanceMap)
	}
	return
}

var uninitializedErr = errors.New("rdb instance uninitialized")
