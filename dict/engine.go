package dict

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/basebytes/component/database/rdb"
)

func Init(config *Config, sources []Source) {
	if err := Reload(config, sources); err != nil {
		panic(err)
	}
}

func Reload(config *Config, sources []Source) (err error) {
	_manager := newManager(config)
	if err = _manager.Init(sources); err == nil {
		if err = _manager.load(); err == nil {
			err = _manager.construct()
		}
	}
	return
}

func newManager(config *Config) *manager {
	return &manager{config: config, biz: NewNormal[*BizDict](sourceBiz)}
}

type manager struct {
	sources map[string]Source
	biz     Source
	config  *Config
}

func (m *manager) Init(sources []Source) (err error) {
	m.sources = make(map[string]Source, len(sources))
	for _, source := range sources {
		if _, ok := m.sources[source.Name()]; ok {
			err = fmt.Errorf("duplicate source[%s]", source.Name())
			return
		}
		m.sources[source.Name()] = source
	}
	for name := range m.config.Source {
		if _, ok := m.sources[name]; !ok {
			err = fmt.Errorf("dict config source[%s] instance not found", name)
			break
		}
	}
	return
}

func (m *manager) load() (err error) {
	for name, config := range m.config.Source {
		if source, ok := m.sources[name]; ok {
			if err = source.Load(config.DBName, config.Params); err == nil && m.config.Backup {
				err = SaveFile(source.Values(), m.config.DictPath, config.Filename)
			}
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		err = m.biz.Load(m.config.DBName, nil)
	}
	return
}

func (m *manager) construct() (err error) {
	if m.config.Action&actionMapping == actionMapping {
		err = m.constructMapping()
	}
	if err == nil && m.config.Action&actionEnum == actionEnum {
		m.constructEnums()
	}
	return
}

func (m *manager) constructMapping() (err error) {
	disableMap, enableMap, _mappings := m.prepare()
	if err = m.upsert(disableMap, enableMap); err == nil {
		resetMapping(_mappings)
	}
	return
}

func (m *manager) prepare() (map[Category]map[string]*BizDict, map[Category]map[string]*BizDict, *mappings) {
	disableMap := make(map[Category]map[string]*BizDict)
	enableMap := make(map[Category]map[string]*BizDict)
	_mappings := newMappings()
	for _, d := range m.biz.Values() {
		if _d, OK := d.(*BizDict); OK {
			if d.GetStatus() == StatusEnable {
				addDict(enableMap, _d)
			} else {
				addDict(disableMap, _d)
			}
		}
	}
	return disableMap, enableMap, _mappings
}

func (m *manager) upsert(disableMap map[Category]map[string]*BizDict, enableMap map[Category]map[string]*BizDict) (err error) {
	for name, source := range m.config.Source {
		if source.Type != SourceTypeMigration {
			continue
		}
		var (
			creates []*BizDict
			updates []rdb.Data
		)
		for _, data := range m.sources[name].Values() {
			if d, OK := data.(*BizDict); OK {
				status := d.GetStatus()
				if _d, ok := findDict(enableMap, d); ok {
					if _d.MappingKey == "" && _d.Id > 0 && _d.updateMigration(d) {
						if _d.GetStatus() == StatusDisable {
							removeDict(enableMap, _d.Category, _d.Key)
							addDict(disableMap, _d)
						}
						updates = append(updates, d)
					}
				} else if status == StatusEnable {
					if _d, ok = findDict(disableMap, d); ok {
						if _d.MappingKey == "" && _d.Id > 0 {
							_d.updateMigration(d)
							removeDict(disableMap, _d.Category, _d.Key)
							addDict(enableMap, _d)
							updates = append(updates, d)
						}
					} else {
						addDict(enableMap, d)
						creates = append(creates, d)
					}
				}
			}
		}
		if ins, ok := rdb.GetConnection(m.config.DBName); ok {
			if len(updates) > 0 {
				err = ins.BatchUpdatesNotEmpty(updates)
			}
			if err == nil && len(creates) > 0 {
				err = ins.Create(creates).Error
			}
		} else {
			err = dbNotFoundErr
		}
	}
	return
}

func (m *manager) constructEnums() {
	_enums := make(map[Category]*Enums)
	_enumMap := make(map[string]*Enum)
	for name, config := range m.config.Source {
		if config.Type == SourceTypeMigration {
			continue
		}
		if source, ok := m.sources[name]; ok {
			enums(_enums, _enumMap, source.Values())
		}
	}
	enums(_enums, _enumMap, m.biz.Values())
	resetEnum(_enums, _enumMap)
	return
}

func enums[T Dict](_enums map[Category]*Enums, _enumMap map[string]*Enum, values []T) {
	for _, d := range values {
		enum := d.Enum()
		key := enum.Unique()
		if d.GetMappingKey() != "" {
			continue
		}
		if _, ok := _enumMap[key]; ok {
			continue
		}
		_enumMap[key] = enum
		if _, ok := _enums[d.GetCategory()]; !ok {
			_enums[d.GetCategory()] = NewEnums(10)
		}
		_enums[d.GetCategory()].Append(enum)
	}
}

// ====================

//func loadPhone(_ *entities.RawDict, config *Config) (err error) {
//	err = loadPhoneData(filepath.Join(config.DictPath, phoneDataFile))
//	return
//}

func LoadFile(o any, paths ...string) error {
	content, e := os.ReadFile(filepath.Join(paths...))
	if e == nil {
		e = json.Unmarshal(content, &o)
	}
	return e
}

func SaveFile(obj any, paths ...string) (err error) {
	filename := filepath.Join(paths...)
	var content []byte
	if content, err = json.MarshalIndent(&obj, "", " "); err == nil {
		err = os.WriteFile(filename, content, os.ModePerm)
	}
	return
}

func findDict(enableMap map[Category]map[string]*BizDict, target *BizDict) (d *BizDict, ok bool) {
	if _, ok = enableMap[target.Category]; ok {
		d, ok = enableMap[target.Category][target.Key]
	}
	return
}

func addDict(dictMap map[Category]map[string]*BizDict, d *BizDict) {
	if _, ok := dictMap[d.Category]; !ok {
		dictMap[d.Category] = make(map[string]*BizDict)
	}
	dictMap[d.Category][d.Key] = d
}

func removeDict(dictMap map[Category]map[string]*BizDict, category Category, key string) {
	delete(dictMap[category], key)
	if len(dictMap[category]) == 0 {
		delete(dictMap, category)
	}
}

const sourceBiz = "biz"

var dbNotFoundErr = errors.New("source db connection not found")
