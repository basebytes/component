package dict

import (
	"fmt"
	"sync"
)

var (
	_nem     map[string]*Enum
	_e       map[Category]*Enums
	enumLock sync.RWMutex
)

func GetEnum(category Category) *Enums {
	enumLock.RLock()
	defer enumLock.RUnlock()
	return _e[category]
}

func GetEnums() map[Category]*Enums {
	enumLock.RLock()
	defer enumLock.RUnlock()
	return _e
}

func AddEnum(dict Dict) {
	enumLock.Lock()
	defer enumLock.Unlock()
	uniKey := dictUniqueKey(dict)
	if _, ok := _nem[uniKey]; !ok {
		enum := dict.Enum()
		_nem[uniKey] = enum
		if dict.GetMappingKey() == "" {
			if _, ok = _e[enum.category]; !ok {
				_e[enum.category] = NewEnums(2)
			}
			_e[enum.category].Append(enum)
		}
	}
}

func UpdateEnum(dict Dict) {
	enumLock.Lock()
	defer enumLock.Unlock()
	newKey := dictUniqueKey(dict)
	enum, ok := _nem[newKey]
	updateFlag := dict.UpdateFlag()
	if ok {
		enum.update(dict.Enum(), updateFlag)
	} else if updateFlag&UpdateFlagStatus == UpdateFlagStatus {
		old := dict.Enum()
		old.Status ^= StatusDisable
		oldKey := old.Unique()
		if enum, ok = _nem[oldKey]; ok {
			enum.update(dict.Enum(), updateFlag)
			delete(_nem, oldKey)
			_nem[newKey] = enum
		}
	}
	if ok && (updateFlag&UpdateFlagMapping == UpdateFlagMapping) {
		if dict.GetMappingKey() == "" {
			if _, ok = _e[enum.category]; !ok {
				_e[enum.category] = NewEnums(2)
			}
			_e[enum.category].Append(enum)
		} else if _enums, OK := _e[enum.category]; OK {
			_enums.Remove(newKey)
			if _enums.Len() == 0 {
				delete(_e, enum.category)
			}
		}
	}
}

func RemoveEnum(key string) {
	enumLock.Lock()
	defer enumLock.Unlock()
	if enum, ok := _nem[key]; ok {
		delete(_nem, key)
		if _enums, OK := _e[enum.category]; OK {
			_enums.Remove(key)
			if _enums.Len() == 0 {
				delete(_e, enum.category)
			}
		}
	}
}

func resetEnum(ne map[Category]*Enums, nem map[string]*Enum) {
	enumLock.Lock()
	defer enumLock.Unlock()
	_e = ne
	_nem = nem
}

func dictUniqueKey(d Dict) string {
	return fmt.Sprintf("%s/%s/%d", d.GetCategory(), d.GetKey(), d.GetStatus())
}
