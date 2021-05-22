package mysql

import (
	"github.com/weblfe/databases/mysql/connection"
	"sync"
)

type (
	DbManager struct {
		Default     string                            `json:"default,default=default"`
		Connections map[string]*connection.ConnConfig `json:"connections"`
	}

	connectionsPool struct {
		locker      sync.RWMutex
		connections map[string]connection.Connector
	}
)

var (
	mgr = &connectionsPool{
		locker:      sync.RWMutex{},
		connections: map[string]connection.Connector{},
	}
)

func (dbMgr DbManager) Get(conn string) connection.Connector {
	if mgr.Exists(conn) {
		return mgr.Get(conn)
	}
	if v, ok := dbMgr.Connections[conn]; ok {
		cfg := v.GetClusterConfig()
		mgr.Set(conn, connection.NewConnector(conn, &cfg))
		return mgr.Get(conn)
	}
	return nil
}

func (mgr *connectionsPool) Exists(name string) bool {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	if _, ok := mgr.connections[name]; ok {
		return true
	}
	return false
}

func (mgr *connectionsPool) Get(name string) connection.Connector {
	mgr.locker.RLock()
	defer mgr.locker.RUnlock()
	if v, ok := mgr.connections[name]; ok {
		return v
	}
	return nil
}

func (mgr *connectionsPool) Set(name string, conn connection.Connector) *connectionsPool {
	if mgr.Exists(name) || conn == nil {
		return mgr
	}
	mgr.locker.Lock()
	defer mgr.locker.Unlock()
	mgr.connections[name] = conn
	return mgr
}
