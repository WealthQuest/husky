package husky

import (
	"strings"

	"github.com/spf13/viper"
)

type Config[C any] interface {
	WithPath(path string) Config[C]
	WithName(name string) Config[C]
	WithType(_type string) Config[C]
	WithEnvPrefix(prefix string) Config[C]
	Load() *C
}

type _Config[C any] struct {
	ins *viper.Viper
}

func (c *_Config[C]) WithPath(path string) Config[C] {
	c.ins.AddConfigPath(path)
	return c
}

func (c *_Config[C]) WithName(name string) Config[C] {
	c.ins.SetConfigName(name)
	return c
}

func (c *_Config[C]) WithType(_type string) Config[C] {
	c.ins.SetConfigType(_type)
	return c
}

func (c *_Config[C]) WithEnvPrefix(prefix string) Config[C] {
	c.ins.SetEnvPrefix(prefix)
	return c
}

func (c *_Config[C]) Load() *C {
	c.ins.AutomaticEnv()
	c.ins.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := c.ins.ReadInConfig(); err != nil {
		panic(err)
	}
	var r C
	if err := c.ins.Unmarshal(&r); err != nil {
		panic(err)
	}
	return &r
}

func NewConfig[C any]() Config[C] {
	ins := viper.New()
	ins.AddConfigPath("./")
	ins.SetConfigName("config")
	ins.SetConfigType("toml")
	ins.SetEnvPrefix("APP")
	return &_Config[C]{ins}
}
