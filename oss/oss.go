package oss

import (
	"errors"
	"sync"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	buckets map[string]map[string]*Oss
	visits  map[string]string
	lock    sync.RWMutex
	once    sync.Once
)

func Init(ossConfigs map[string]*Config) {
	once.Do(func() {
		lock.Lock()
		defer lock.Unlock()
		if _buckets, _visits, err := load(ossConfigs); err == nil {
			buckets = _buckets
			visits = _visits
		} else {
			panic(err)
		}
	})
}

func load(ossConfigs map[string]*Config) (_buckets map[string]map[string]*Oss, _visits map[string]string, err error) {
	_buckets = make(map[string]map[string]*Oss, len(ossConfigs))
	_visits = make(map[string]string)
	for name, config := range ossConfigs {
		var c *oss.Client
		if c, err = oss.New(config.Scheme+"://"+config.EndPoint, config.AccessKeyId, config.AccessKeySecret,
			oss.Timeout(config.ConnectTimeoutSec(), config.TimeoutSec())); err == nil {
			_buckets[name] = make(map[string]*Oss, len(config.Buckets))
			if config.Visit != "" {
				_visits[name] = config.Visit
			}
			for k, v := range config.Buckets {
				var bucket *oss.Bucket
				if bucket, err = c.Bucket(v); err != nil {
					return
				}
				_buckets[name][k] = NewOss(bucket, Scheme(config.Scheme), Endpoint(config.EndPoint),
					Domain(config.Visit), ExpiredInSec(config.Expired.Milliseconds()/1000))
			}
		}
	}
	return
}

func Reload(ossConfigs map[string]*Config) error {
	if buckets == nil {
		return uninitializedErr
	}
	_buckets, _visits, err := load(ossConfigs)
	if err == nil {
		lock.Lock()
		defer lock.Unlock()
		buckets = _buckets
		visits = _visits
	}
	return err
}

func GetDefaultBucket(key string) (*Oss, error) {
	return GetBucket(DefaultClientName, key)
}

func GetDefaultVisit() string {
	return GetVisit(DefaultClientName)
}

func GetVisit(name string) string {
	lock.RLock()
	defer lock.RUnlock()
	return visits[name]
}

func GetBucket(name, key string) (s *Oss, e error) {
	lock.RLock()
	defer lock.RUnlock()
	_buckets, ok := buckets[name]
	if ok {
		s, ok = _buckets[key]
	}
	if !ok {
		e = bucketNotFoundErr
	}
	return
}

var (
	bucketNotFoundErr = errors.New("bucket not found")
	uninitializedErr  = errors.New("oss instance uninitialized")
)

const DefaultClientName = "default"
