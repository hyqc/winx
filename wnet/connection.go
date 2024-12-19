package wnet

import (
	"errors"
	"io"
	"net"
	"winx/global"
	"winx/wiface"
)

type Connection struct {
	// 当前连接属于哪个Server
	Server wiface.IServer
	//当前连接
	Conn *net.TCPConn
	// 当前连接的ID号，全局唯一
	ConnID uint32
	// 当前连接的关闭状态
	isClosed bool
	//通知退出/停止的channel
	ExitBuffChan chan bool
	//处理方法
	MsgHandler wiface.IMsgHandler
	//无缓冲通道，用于读写goroutine之间的通信
	msgChan chan []byte
	// 缓冲队列，用于读写goroutine之间的通信
	msgBuffChan chan []byte
}

// NewConnection 创建连接的方法
func NewConnection(ser wiface.IServer, conn *net.TCPConn, connID uint32, handler wiface.IMsgHandler) *Connection {

	c := &Connection{
		Server:       ser,
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   handler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, global.Conf.MaxWorkerTaskLen),
	}
	// 将当前连接加入管理器
	c.Server.GetConnMgr().Add(c)

	return c
}

// startReader 处理读数据
func (c *Connection) startReader() {
	SysPrintInfo("Reader Goroutine is running")
	defer SysPrintInfo("connID: ", c.ConnID, " Reader is exit, remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()
	for {
		//拆包
		dp := NewDataPack()
		// 读取客户端消息header
		headerData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetConnection(), headerData); err != nil {
			SysPrintError("read header error:", err)
			c.ExitBuffChan <- true
			continue
		}
		SysPrintInfo("unpack header: ", string(headerData))
		//拆解包头
		msg, err := dp.Unpack(headerData)
		if err != nil {
			SysPrintError("unpack error:", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆解body
		var body []byte
		if msg.GetDataLen() > 0 {
			body = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetConnection(), body); err != nil {
				SysPrintError("read body error:", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(body)

		SysPrintInfo("read from client msgId: ", msg.GetMsgID())
		SysPrintInfo("read from client msgDataLen: ", msg.GetDataLen())
		SysPrintInfo("read from client msgData: ", string(msg.GetData()))
		req := NewRequest(c, msg)
		if global.Conf.WorkerPoolSize > 0 {
			Print(1111)
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

func (c *Connection) startWriter() {
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				SysPrintInfo("send msg error:", err)
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					SysPrintError("send msg error:", err)
					return
				}
			} else {
				SysPrintWarn("msgBuffChan is Closed")
				break
			}

		case <-c.ExitBuffChan:
			return
		}
	}
}

// Stop 停止连接
func (c *Connection) Stop() {
	SysPrintInfo("Conn Stop()...ConnID = ", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	c.Server.CallBeforeConnStop(c)

	c.Conn.Close()
	c.ExitBuffChan <- true

	c.Server.GetConnMgr().Remove(c)

	close(c.ExitBuffChan)
	close(c.msgChan)
	close(c.msgBuffChan)
}

func (c *Connection) Start() {
	go c.startReader()
	go c.startWriter()

	c.Server.CallAfterConnStart(c)

	select {
	case <-c.ExitBuffChan:
		SysPrintInfo("connection exit!")
		return
	}
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetConnection() net.Conn {
	return c.Conn
}

func (c *Connection) SendMsg(msgId uint32, body []byte) error {
	if c.isClosed {
		SysPrintWarn("connection closed when send msg")
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, body))
	if err != nil {
		SysPrintError("pack error:", err)
		return err
	}
	// 写回客户端
	c.msgChan <- msg
	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.isClosed {
		SysPrintWarn("connection closed when send msg")
		return errors.New("connection closed when send msg")
	}
	dp := NewDataPack()
	msg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		SysPrintError("pack error:", err)
		return err
	}
	c.msgBuffChan <- msg
	return nil
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
