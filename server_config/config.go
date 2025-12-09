package server_config

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/basebytes/component/database/rdb"
	"github.com/basebytes/config-manager-go/config"
)

type serverConfig struct {
	ServerName  string      `json:"serverName,omitempty"`
	StartUp     string      `json:"startUp,omitempty"`
	TablePrefix string      `json:"tablePrefix,omitempty"`
	FromDB      *rdb.Config `json:"fromDB,omitempty"`
}

func Load(configPath string, cfg any, options ...config.Option) (fromDB *rdb.Config, err error) {
	options = append(options, config.WithDefaultConfig("startUp", "startup"), config.WatchConfigFile(false))
	var (
		loader     = newLoader(withDateTimeDecoder)
		cfgManager = config.New(loader, configPath, options...)
		ins        *rdb.Instance
		serverCfg  = &serverConfig{}
	)
	if err = cfgManager.ReadConfig(cfg); err == nil {
		if serverCfg.FromDB != nil {
			if err = serverCfg.FromDB.Init(); err != nil {
				err = fmt.Errorf("init config db items failed:%s", err.Error())
			} else if ins, err = rdb.NewInstance("", serverCfg.FromDB); err == nil {
				sCfg := NewConfig(serverCfg.ServerName, serverCfg.StartUp)
				if err = ins.FindFirstByCondition(sCfg).Error; err == nil {
					if len(sCfg.Content) == 0 {
						err = errors.New("server config is empty")
					} else {
						err = loader.ReadConfig(bytes.NewReader(sCfg.Content))
					}
				}
			}
			if err == nil {
				fromDB = serverCfg.FromDB
				err = loader.Unmarshal(&cfg)
			}
		}
	}
	return
}
