package wnet

import (
	"fmt"
	"sync"
	"sync/atomic"
	"winx/wiface"
)

type ConnManager struct {
	connections map[uint32]wiface.IConnection
	lock        *sync.RWMutex
	length      atomic.Int64
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]wiface.IConnection),
		lock:        &sync.RWMutex{},
		length:      atomic.Int64{},
	}
}

func (c *ConnManager) Add(conn wiface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[conn.GetConnID()] = conn
	c.length.Add(1)
	SysPrintInfo("connection add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn wiface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.connections, conn.GetConnID())
	c.length.Add(-1)
	SysPrintInfo("connection remove to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Len() int {
	return int(c.length.Load())
}

func (c *ConnManager) Get(connID uint32) (wiface.IConnection, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, fmt.Errorf("connection not found")
}

func (c *ConnManager) ClearConn() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	SysPrintInfo("connection clea r to ConnManager successfully: conn num = ", c.Len())
}
