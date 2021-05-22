package mysql

import (
	"encoding/json"
	"fmt"
	"github.com/weblfe/databases/mysql/connection"
	"github.com/weblfe/databases/mysql/types"
	"os"
	"strconv"
)

func DebugOn() bool {
	var v = os.Getenv(types.DebugEnvKey)
	if v != "" && (v == "1" || v == "on" || v == "true") {
		return true
	}
	return false
}

func NewDbManagerByEnv(envPrefix ...string) (*DbManager, error) {
	var mgr = &DbManager{}
	if len(envPrefix) <= 0 {
		mgr.Default = getEnvWithDefault("DATABASES_DEFAULT", "default")
		mgr.Connections = map[string]*connection.ConnConfig{
			mgr.Default: NewConnConfigByEnv(makeEnvKey(mgr.Default, "")),
		}
	} else {
		prefix := envPrefix[0]
		mgr.Default = getEnvWithDefault(makeEnvKey("DATABASES_DEFAULT", prefix), "default")
		mgr.Connections = map[string]*connection.ConnConfig{
			mgr.Default: NewConnConfigByEnv(makeEnvKey(mgr.Default, prefix)),
		}
	}
	return mgr, nil
}

func getEnvWithDefault(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func makeEnvKey(key string, prefix string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%s_%s", prefix, key)
}

func NewConnConfigByEnv(prefix string) *connection.ConnConfig {
	var cfg = &connection.ConnConfig{
		Driver:          getEnvWithDefault(makeEnvKey("DRIVER", prefix), types.DriverName),
		Host:            getEnvWithDefault(makeEnvKey("HOST", prefix), types.DefaultHost),
		Port:            str2IntWithDefault(getEnvWithDefault(makeEnvKey("PORT", prefix), "3306"), types.DefaultPort),
		User:            getEnvWithDefault(makeEnvKey("USER", prefix), types.DefaultUser),
		Password:        getEnvWithDefault(makeEnvKey("PASSWORD", prefix), types.DefaultPassword),
		Database:        getEnvWithDefault(makeEnvKey("DATABASE", prefix), types.DefaultDbName),
		SetMaxOpenConns: str2IntWithDefault(getEnvWithDefault(makeEnvKey("SET_MAX_OPEN_CONNS", prefix), "0"), 0),
		SetMaxIdleConns: str2IntWithDefault(getEnvWithDefault(makeEnvKey("SET_MAX_IDLE_CONNS", prefix), "0"), 0),
		Options:         getEnvWithDefault(makeEnvKey("OPTIONS", prefix), types.DefaultOptions),
		Prefix:          getEnvWithDefault(makeEnvKey("PREFIX", prefix), ""),
		Master:          strJsonToCfgArr(getEnvWithDefault(makeEnvKey("MASTER", prefix), "[]")),
		Slave:           strJsonToCfgArr(getEnvWithDefault(makeEnvKey("SLAVE", prefix), "[]")),
	}
	return cfg
}

func str2IntWithDefault(str string, def int) int {
	n, err := strconv.Atoi(str)
	if err != nil {
		return def
	}
	return n
}

func strJsonToCfgArr(str string) []*connection.ConnConfig {
	if str == "" || str == "[]" || str == "{}" {
		return nil
	}
	var arr []*connection.ConnConfig
	if err := json.Unmarshal([]byte(str), &arr); err != nil {
		return nil
	}
	return arr
}
