package cfg

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

const (
	cfgFile = "cfg_%s.json"
	envPrefix = "blog"
)

var vc = viperCfg{}

func init() {
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	if vc.GetBool("production") {
		viper.SetConfigFile(fmt.Sprintf(cfgFile, "prod"))
	} else {
		viper.SetConfigFile(fmt.Sprintf(cfgFile, "dev"))
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

type viperCfg struct {
}

func (c *viperCfg) Get(key string) interface{} {
	return viper.Get(key)
}

func (c *viperCfg) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (c *viperCfg) GetFloat64(key string) float64 {
	return viper.GetFloat64(key)
}

func (c *viperCfg) GetInt(key string) int {
	return viper.GetInt(key)
}


func (c *viperCfg) GetString(key string) string {
	return viper.GetString(key)
}

func (c *viperCfg) GetTime(key string) time.Time {
	return viper.GetTime(key)
}

func (c *viperCfg) GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func (c *viperCfg) IsSet(key string) bool {
	return viper.IsSet(key)
}

func (c *viperCfg) AllSettings() map[string]interface{} {
	return viper.AllSettings()
}
