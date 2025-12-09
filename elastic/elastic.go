package elastic

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/basebytes/elastic-go/client"
	"github.com/basebytes/elastic-go/service"
)

var (
	elasticService *service.Service
	lock           sync.RWMutex
	once           sync.Once
)

func Init(esCfg *client.Config) {
	once.Do(func() {
		lock.Lock()
		defer lock.Unlock()
		if _service, err := reload(esCfg); err == nil {
			elasticService = _service
		} else {
			panic(err)
		}
	})
}

func Reload(esCfg *client.Config) error {
	if GetService() == nil {
		return uninitializedErr
	}
	_service, err := reload(esCfg)
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		elasticService = _service
	}
	return err
}

func GetService() *service.Service {
	lock.RLock()
	defer lock.RUnlock()
	return elasticService
}

func reload(esCfg *client.Config) (*service.Service, error) {
	_service, err := service.NewService(esCfg)
	if err == nil && !_service.Ping() {
		err = fmt.Errorf("connect to elasticsearch server %s failed ", strings.Join(esCfg.Servers, ","))
	}
	return _service, err
}

var uninitializedErr = errors.New("service uninitialized")
