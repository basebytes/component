package dict

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/basebytes/component/database/rdb"
	"github.com/basebytes/tools"
)

const defaultDictPath = "dict"

func NewConfig() *Config {
	return &Config{}
}

type Config struct {
	DBName   string                   `json:"dbName,omitempty"`
	DictPath string                   `json:"dictPath,omitempty"`
	Backup   bool                     `json:"backup,omitempty"`
	Action   int                      `json:"action,omitempty"`
	Source   map[string]*SourceConfig `json:"source,omitempty"`
}

func (c *Config) Init(mountPath string) (err error) {
	if c.DictPath == "" {
		c.DictPath = filepath.Join(mountPath, defaultDictPath)
	}
	if c.Action &= actionMask; c.Action == 0 {
		err = fmt.Errorf("invalid dict action[%d]", c.Action)
	}
	if err == nil {
		for name, source := range c.Source {
			if err = source.Init(name, c.DBName); err != nil {
				break
			}
		}
	}
	if err == nil {
		err = tools.CreateDirIfNotExists(c.DictPath, os.ModeDir|0755)
	}
	return
}

type SourceConfig struct {
	Type     string         `json:"type,omitempty"`
	DBName   string         `json:"dbName,omitempty"`
	Filename string         `json:"filename,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

func (s *SourceConfig) Init(name, dbname string) (err error) {
	switch {
	case s.Type == "":
		err = fmt.Errorf("dict source[%s] required configuration[type] not found", name)
	case s.Type != "" && s.Type != SourceTypeMigration && s.Type != SourceTypeRemote && s.Type != SourceTypeLocal:
		err = fmt.Errorf("invalid configuration[type] value[%s] for dict source[%s]", s.Type, name)
	default:
		if s.DBName == "" {
			s.DBName = dbname
		}
		if !rdb.ValidName(s.DBName) {
			err = fmt.Errorf("configuration[dbName] value[%s] for dict source[%s] not found", s.DBName, name)
		}
	}
	if err == nil && s.Filename == "" {
		s.Filename = fmt.Sprintf("%s.json", name)
	}
	return
}

const (
	SourceTypeMigration = "migration"
	SourceTypeRemote    = "remote"
	SourceTypeLocal     = "local"
)

const (
	actionMapping = 1
	actionEnum    = 2
	actionMask    = actionMapping | actionEnum
)
