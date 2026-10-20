package dict

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/basebytes/component/database/rdb"
	"github.com/basebytes/tools"
)

const defaultDictPath = "dict"

func NewConfig(dbName string) *Config {
	return &Config{DBName: dbName}
}

type Config struct {
	DBName  string `json:"dbName,omitempty"`
	DictDir string `json:"dictDir,omitempty"`
	Backup  bool   `json:"backup,omitempty"`
}

func (c *Config) Init(mountPath string) (err error) {
	if c.DictDir == "" {
		c.DictDir = filepath.Join(mountPath, defaultDictPath)
	}
	err = tools.CreateDirIfNotExists(c.DictDir, os.ModeDir|0755)
	return
}

type SourceConfig struct {
	DBName   string         `json:"dbName,omitempty"`
	Filename string         `json:"filename,omitempty"`
	Params   map[string]any `json:"params,omitempty"`
}

func (s *SourceConfig) Init(name, dbname string) (err error) {
	if s.DBName == "" {
		s.DBName = dbname
	}
	if !rdb.ValidName(s.DBName) {
		err = fmt.Errorf("configuration[dbName] value[%s] for dict source[%s] not found", s.DBName, name)
	} else if s.Filename == "" {
		s.Filename = fmt.Sprintf("%s.json", name)
	}
	return
}
