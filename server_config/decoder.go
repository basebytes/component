package server_config

import (
	"reflect"
	"time"

	"github.com/basebytes/types"
	"github.com/mitchellh/mapstructure"
)

var withDateTimeDecoder = func(c *mapstructure.DecoderConfig) {
	c.DecodeHook = func(f reflect.Type, t reflect.Type, data any) (newData any, err error) {
		if f.Kind() == reflect.String {
			switch t {
			case reflect.TypeOf(zeroDur.Duration):
				newData, err = time.ParseDuration(data.(string))
			case reflect.TypeOf(&zeroDur):
				var _d time.Duration
				if _d, err = time.ParseDuration(data.(string)); err == nil {
					newData = &types.Duration{Duration: _d}
				}
			case reflect.TypeOf(&zeroT):
				var (
					_t      time.Time
					timeFmt = "2006-01-02 15:04:05"
				)
				str := data.(string)
				if len(str) < len(timeFmt) {
					timeFmt = "2006-01-02"
					str = str[:10]
				}
				if _t, err = time.ParseInLocation(timeFmt, str, time.Local); err == nil {
					newData = types.Time{Time: _t}
				}
			case reflect.TypeOf(&zeroD):
				var (
					_t      time.Time
					dateFmt = "2006-01-02"
				)
				str := data.(string)
				if len(str) > len(dateFmt) {
					str = str[:10]
				}
				if _t, err = time.ParseInLocation(dateFmt, str, time.Local); err == nil {
					newData = types.Date{Time: _t}
				}
			default:
				return data, nil
			}
			return
		}
		return data, nil
	}
}

var (
	zeroD   types.Date
	zeroT   types.Time
	zeroDur types.Duration
)
