package wnet

import (
	"fmt"
	"sync"
	"winx/wiface"
)

type ConnManager struct {
	connections map[uint32]wiface.IConnection
	lock        sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]wiface.IConnection),
		lock:        sync.RWMutex{},
	}
}

func (c *ConnManager) Add(conn wiface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Remove(conn wiface.IConnection) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.connections, conn.GetConnID())
	fmt.Println("connection remove to ConnManager successfully: conn num = ", c.Len())
}

func (c *ConnManager) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.connections)
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
	fmt.Println("connection clea r to ConnManager successfully: conn num = ", c.Len())
}
