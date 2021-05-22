package connection

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/gorose/v2"
	"github.com/weblfe/databases/mysql/types"
)

type ConnConfig struct {
	Driver          string        `json:"driver,default=mysql"`
	Host            string        `json:"host,default=127.0.0.1"`
	Port            int           `json:"port,default=3306"`
	User            string        `json:"user,default=root"`
	Password        string        `json:"password,default=root"`
	Database        string        `json:"database,default=test"`
	SetMaxOpenConns int           `json:"setMaxOpenConns"` // (连接池)最大打开的连接数，默认值为0表示不限制
	SetMaxIdleConns int           `json:"setMaxIdleConns"` // (连接池)闲置的连接数, 默认0
	Options         string        `json:"options,default=\"charset=utf8&parseTime=true\""`
	Prefix          string        `json:"prefix,default=''"`
	Master          []*ConnConfig `json:"master"`
	Slave           []*ConnConfig `json:"slave"`
}

var (
	errC error
)

func (c ConnConfig) String() string {
	v, err := json.Marshal(&c)
	if err != nil {
		errC = err
		return ""
	}
	return string(v)
}

// "root:root@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true"
func (c ConnConfig) GetDsnUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		c.getUser(), c.getPassword(), c.getHost(), c.getPort(),
		c.getDb(), c.getOptions(),
	)
}

func (c ConnConfig) getDriver() string {
	if c.Driver == "" {
		return types.DriverName
	}
	return c.Driver
}

func (c ConnConfig) getUser() string {
	if c.User == "" {
		return types.DefaultUser
	}
	return c.User
}

func (c ConnConfig) getPassword() string {
	if c.Password == "" {
		return types.DefaultPassword
	}
	return c.Password
}

func (c ConnConfig) getHost() string {
	if c.Host == "" {
		return types.DefaultHost
	}
	return c.Host
}

func (c ConnConfig) getPort() int {
	if c.Port <= 0 || c.Port > 65535 {
		return types.DefaultPort
	}
	return c.Port
}

func (c ConnConfig) getDb() string {
	if c.Database == "" {
		return types.DefaultDbName
	}
	return c.Database
}

func (c ConnConfig) getOptions() string {
	if c.Options == "" {
		return types.DefaultOptions
	}
	return c.Options
}

func (c ConnConfig) GetConfig() gorose.Config {
	return gorose.Config{
		Driver:          c.getDriver(),
		Dsn:             c.GetDsnUrl(),
		SetMaxIdleConns: c.SetMaxIdleConns,
		SetMaxOpenConns: c.SetMaxOpenConns,
		Prefix:          c.Prefix,
	}
}

func (c ConnConfig) GetClusterConfig() gorose.ConfigCluster {
	var cfg = gorose.ConfigCluster{
		Prefix: c.Prefix,
		Driver: c.getDriver(),
	}
	if len(c.Master) > 0 && len(c.Slave) > 0 {
		for _, v := range c.Master {
			it:=v.GetConfig()
			if it.Prefix == "" && c.Prefix!= "" {
				it.Prefix = c.Prefix
			}
			if it.Driver == "" && c.Driver!= "" {
				it.Driver = c.Driver
			}
			cfg.Master = append(cfg.Master,it)
		}
		for _, v := range c.Slave {
			it:=v.GetConfig()
			if it.Prefix == "" && c.Prefix!= "" {
				it.Prefix = c.Prefix
			}
			if it.Driver == "" && c.Driver!= "" {
				it.Driver = c.Driver
			}
			cfg.Slave = append(cfg.Slave, it)
		}
	} else {
		cfg.Slave = append(cfg.Slave, c.GetConfig())
		cfg.Master = append(cfg.Master, c.GetConfig())
	}
	return cfg
}

func GetConnConfigError() error {
	return errC
}

func NewConfigByBytes(data []byte) (*ConnConfig, error) {
	var cfg = &ConnConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
