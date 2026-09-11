package oss

import (
	"fmt"
	"time"

	"github.com/basebytes/types"
)

const (
	defaultEndPoint = "oss-cn-beijing.aliyuncs.com"
	defaultScheme   = SchemeHttps
	SchemeHttp      = "http"
	SchemeHttps     = "https"
)

var (
	defaultReadWriteTimeout = &types.Duration{Duration: time.Minute}
	defaultConnectTimeout   = &types.Duration{Duration: time.Second * 20}
	defaultExpired          = &types.Duration{Duration: time.Hour * 2}
)

type Config struct {
	AccessKeyId     string            `json:"accessKeyId,omitempty"`
	AccessKeySecret string            `json:"accessKeySecret,omitempty"`
	Buckets         map[string]string `json:"buckets,omitempty"`
	EndPoint        string            `json:"endPoint,omitempty"`
	Scheme          string            `json:"scheme,omitempty"`
	Expired         *types.Duration   `json:"expired,omitempty"`
	Timeout         *types.Duration   `json:"timeout,omitempty"`
	ConnectTimeout  *types.Duration   `json:"connectTimeout,omitempty"`
	Visit           string            `json:"visit,omitempty"`
}

func (c *Config) Init() (err error) {
	if c.AccessKeyId == "" || c.AccessKeySecret == "" || len(c.Buckets) == 0 {
		err = fmt.Errorf("oss config required configuration accessKeyId,accessKeySecret or bucket")
	} else {
		if c.Timeout == nil {
			c.Timeout = defaultReadWriteTimeout
		}
		if c.ConnectTimeout == nil {
			c.ConnectTimeout = defaultConnectTimeout
		}
		if c.EndPoint == "" {
			c.EndPoint = defaultEndPoint
		}
		if c.Scheme != SchemeHttp {
			c.Scheme = defaultScheme
		}
		if c.Expired == nil || c.Expired.Duration <= 0 {
			c.Expired = defaultExpired
		}
	}
	return
}

func (c *Config) ConnectTimeoutSec() int64 {
	return c.ConnectTimeout.Milliseconds() / 1000
}

func (c *Config) TimeoutSec() int64 {
	return c.Timeout.Milliseconds() / 1000
}
