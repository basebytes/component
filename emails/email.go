package emails

import (
	"strings"
	"sync"

	"github.com/basebytes/email"
)

var (
	clients   = make(map[string]*email.Client)
	receivers = make(map[string]map[string]string)
	lock      sync.RWMutex
	once      sync.Once
)

func Init(configs []*EmailConfig) {
	once.Do(func() {
		lock.Lock()
		defer lock.Unlock()
		if _clients, _receivers, err := reload(configs); err == nil {
			clients, receivers = _clients, _receivers
		} else {
			panic(err)
		}
	})
}

func Reload(configs []*EmailConfig) error {
	_clients, _receivers, err := reload(configs)
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		clients, receivers = _clients, _receivers
	}
	return err
}

func GetDefaultClient() *email.Client {
	return GetClient(DefaultClientName)
}

func GetDefaultReceivers() map[string]string {
	return GetReceivers(DefaultClientName)
}

func GetDefaultReceiver(key string) (receiver string) {
	return GetReceiver(DefaultClientName, key)
}

func GetClient(name string) *email.Client {
	lock.RLock()
	defer lock.RUnlock()
	return clients[name]
}

func GetReceivers(name string) map[string]string {
	lock.RLock()
	defer lock.RUnlock()
	if rec, OK := receivers[name]; OK {
		return rec
	}
	return emptyMap
}

func GetReceiver(name string, key string) (receiver string) {
	lock.RLock()
	defer lock.RUnlock()
	if rec, OK := receivers[name]; OK {
		receiver = rec[strings.ToUpper(key)]
	}
	return
}

func reload(configs []*EmailConfig) (map[string]*email.Client, map[string]map[string]string, error) {
	_clients := make(map[string]*email.Client)
	_receivers := make(map[string]map[string]string)
	for _, cfg := range configs {
		if client, err := email.NewClient(&email.Config{
			Server:   cfg.Server,
			User:     cfg.User,
			Password: cfg.Password,
		}); err == nil {
			_clients[cfg.Name] = client
			upped := make(map[string]string, len(cfg.Receivers))
			for k, v := range cfg.Receivers {
				upped[strings.ToUpper(k)] = v
			}
			_receivers[cfg.Name] = upped
		} else {
			return nil, nil, err
		}
	}
	return _clients, _receivers, nil
}

const DefaultClientName = "default"

var emptyMap = make(map[string]string)
