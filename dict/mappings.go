package dict

import "sync"

var (
	_m          *mappings
	mappingLock sync.RWMutex
)

func GetMappingKey(category Category, oriKey string) (mappingKey string) {
	mappingLock.RLock()
	defer mappingLock.RUnlock()
	return _m.GetValue(category, oriKey)
}

func SetMappingKey(category Category, oriKey, mappingKey string) {
	mappingLock.Lock()
	defer mappingLock.Unlock()
	_m.AppendChild(category, oriKey, mappingKey)
}

func resetMapping(nm *mappings) {
	mappingLock.Lock()
	defer mappingLock.Unlock()
	_m = nm
}

func newMappings() *mappings {
	return &mappings{}
}

type mappings map[Category]map[string]string

func (d *mappings) AppendChild(category Category, k, v string) {
	if _, OK := (*d)[category]; !OK {
		(*d)[category] = make(map[string]string)
	}
	(*d)[category][k] = v
}

func (d *mappings) GetValue(category Category, k string) (v string) {
	if _, OK := (*d)[category]; OK {
		v = (*d)[category][k]
	}
	return
}
