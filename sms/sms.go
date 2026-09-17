package sms

import (
	"errors"
	"fmt"
	"sync"
)

var (
	clients       map[string]Client
	newClientFunc = make(map[string]newClient)
	lock          sync.RWMutex
)

func Init(configs []*Config) {
	if err := Reload(configs); err != nil {
		panic(err)
	}
}

func Reload(configs []*Config) (err error) {
	_clients := make(map[string]Client)
	for _, cfg := range configs {
		var client Client
		if newFunc, ok := newClientFunc[cfg.Name]; ok {
			client, err = newFunc(cfg)
		} else {
			err = fmt.Errorf("invalid sms config name %s", cfg.Name)
		}
		if err != nil {
			break
		}
		_clients[cfg.Name] = client
	}
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		clients = _clients
	}
	return
}

func SendTemplateSms(clientName string, templateName string, receivers string, params map[string]string, result any) (err error) {
	var client Client
	if client, err = getClient(clientName); err == nil {
		err = client.SendTemplateSms(templateName, receivers, params, result)
	}
	return
}

func GetTemplate(clientName, templateId string, result any) (err error) {
	var client Client
	if client, err = getClient(clientName); err == nil {
		err = client.GetTemplate(templateId, result)
	}
	return
}

func PullMessageStatus(clientName, taskId string, result any) (err error) {
	var client Client
	if client, err = getClient(clientName); err == nil {
		err = client.PullMessageStatus(taskId, &result)
	}
	return
}

func getClient(clientName string) (Client, error) {
	lock.RLock()
	defer lock.RUnlock()
	if client, ok := clients[clientName]; ok {
		return client, nil
	}
	return nil, notFoundError
}

func RegisterNewClientFunc(name string, newClient newClient) (ok bool) {
	if _, ok = newClientFunc[name]; !ok {
		newClientFunc[name] = newClient
		ok = true
	}
	return
}

var notFoundError = errors.New("sms client not found")
