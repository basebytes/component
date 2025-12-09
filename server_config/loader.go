package server_config

import "github.com/spf13/viper"

func newLoader(opts ...viper.DecoderConfigOption) *configLoader {
	return &configLoader{Viper: viper.NewWithOptions(), opts: opts}
}

type configLoader struct {
	*viper.Viper
	opts []viper.DecoderConfigOption
}

func (cl *configLoader) Load(path string, v any) (err error) {
	cl.SetConfigFile(path)
	err = cl.ReadInConfig()
	if err == nil {
		err = cl.Unmarshal(v)
	}
	return err
}

func (cl *configLoader) Save(path string, v any) error {
	cl.SetConfigFile(path)
	return cl.WriteConfig()
}

func (cl *configLoader) Unmarshal(v any) error {
	return cl.Viper.Unmarshal(v, cl.opts...)
}

func (cl *configLoader) UnmarshalKey(key string, v any) error {
	return cl.Viper.UnmarshalKey(key, v, cl.opts...)
}
