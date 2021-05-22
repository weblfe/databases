package connection

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2"
)

// 链接器
type (
	Connector interface {
		Init() (bool, error)
		Name() string
		GetDb() gorose.IEngin
		GetOrm() gorose.IOrm
		Prefix() string
	}

	ConnectorImpl struct {
		name      string
		cfg       *gorose.ConfigCluster
		db        gorose.IEngin
		errDriver error
	}
)

func NewConnector(name string, cfg ...*gorose.ConfigCluster) *ConnectorImpl {
	var impl = &ConnectorImpl{
		name: name,
	}
	if len(cfg) > 0 && cfg[0] != nil {
		impl.SetConfig(*cfg[0])
	}
	return impl
}

func (impl *ConnectorImpl) SetConfig(cfg gorose.ConfigCluster) *ConnectorImpl {
	if impl.cfg == nil {
		impl.cfg = &cfg
	}
	return impl
}

func (impl *ConnectorImpl) Name() string {
	return impl.name
}

func (impl *ConnectorImpl) Init() (bool, error) {
	if impl.cfg == nil {
		return false, errors.New("empty cfg")
	}
	if impl.db != nil {
		return true, nil
	}
	impl.db, impl.errDriver = gorose.NewEngin(impl.cfg)
	if impl.db != nil && impl.errDriver == nil {
		return true, nil
	}
	return false, impl.errDriver
}

func (impl *ConnectorImpl) GetDb() gorose.IEngin {
	if impl.db == nil {
		_, impl.errDriver = impl.Init()
	}
	return impl.db
}

func (impl *ConnectorImpl) Prefix() string {
	if impl.cfg != nil {
		return impl.cfg.Prefix
	}
	return ""
}

func (impl *ConnectorImpl) GetOrm() gorose.IOrm {
	if impl.db == nil {
		_, impl.errDriver = impl.Init()
	}
	if db, ok := impl.db.(*gorose.Engin); ok {
		db.SetPrefix(impl.Prefix())
		return db.NewOrm()
	}
	return nil
}
