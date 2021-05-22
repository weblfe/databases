package mysql

import (
	"github.com/gohouse/gorose/v2"
)

type DataMapper interface {
	TableName() string
	New() interface{}
	Collection() interface{}
}

type Model struct {
	gorose.IOrm
	db         *DbManager
	connection string
	mapper     DataMapper
}

func (m *Model) SetConnection(name string) *Model {
	if m.connection == "" {
		m.connection = name
	}
	return m
}

func (m *Model) SetDbManger(mgr *DbManager) *Model {
	if m.db == nil && mgr != nil {
		m.db = mgr
	}
	return m
}

func (m *Model) Bind(mapper DataMapper) *Model {
	if m.mapper == nil && mapper != nil {
		m.mapper = mapper
	}
	return m
}

func (m *Model) GetMapper() DataMapper {
	if m.mapper == nil {
		return nil
	}
	return m.mapper
}

func (m *Model) One() (interface{}, error) {
	var mapper = m.mapper.New()
	if err := m.GetOrm().Table(mapper).Limit(1).Select(); err != nil {
		return nil, err
	}
	return mapper, nil
}

func (m *Model) Find() (interface{}, error) {
	var mapper = m.mapper.New()
	if err := m.GetOrm().Table(mapper).Limit(1).Select(); err != nil {
		return nil, err
	}
	return mapper, nil
}

func (m *Model) All() (interface{}, error) {
	var mapperArr = m.mapper.Collection()
	if err := m.GetOrm().Table(&mapperArr).Select(); err != nil {
		return nil, err
	}
	return mapperArr, nil
}

func (m *Model) Get() (interface{}, error) {
	var mapperArr = m.mapper.Collection()
	if err := m.GetOrm().Table(&mapperArr).Select(); err != nil {
		return nil, err
	}
	return mapperArr, nil
}

func (m *Model) GetOrm() gorose.IOrm {
	return m.InitOrm().IOrm
}

func (m *Model) InitOrm() *Model {
	if m.IOrm == nil {
		m.IOrm = m.db.Get(m.connection).GetOrm()
	}
	return m
}

func (m *Model) Init(db ...*DbManager) {
	if m.db == nil && len(db) > 0 && db[0] != nil {
		m.SetDbManger(db[0]).SetConnection(db[0].Default).InitOrm()
	}
}

func (m *Model) Reset() gorose.IOrm {
	return m.GetOrm().Reset()
}

func (m *Model) Query(binds ...interface{}) gorose.IOrm {
	if len(binds) > 0 && binds[0] != nil {
		return m.GetOrm().Reset().Table(binds[0])
	}
	return m.GetOrm().Reset().Table(m.GetMapper().TableName())
}

func (m *Model) TableName() string {
	if m.mapper != nil {
		return m.db.Get(m.connection).Prefix() + m.mapper.TableName()
	}
	return ""
}
