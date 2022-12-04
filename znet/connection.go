package znet

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/golang-framework/tcpx/confs"
	"github.com/golang-framework/tcpx/z"
	"github.com/golang-framework/tcpx/ziface"
)

// Connection 链接
type Connection struct {
	//当前Conn属于哪个Server
	TCPServer ziface.IServer
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//消息管理MsgID和对应处理方法的消息管理模块
	MsgHandler ziface.IMsgHandle
	//告知该链接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc
	//有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	sync.RWMutex
	//链接属性
	property map[string]interface{}
	////保护当前property的锁
	propertyLock sync.Mutex
	//当前连接的关闭状态
	isClosed bool
}

// NewConnection 创建连接的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TCPServer:   server,
		Conn:        conn,
		ConnID:      connID,
		isClosed:    false,
		MsgHandler:  msgHandler,
		msgBuffChan: make(chan []byte, confs.MaxMsgChanLen),
		property:    nil,
	}

	//将新创建的Conn添加到链接管理中
	c.TCPServer.GetConnMgr().Add(c)
	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("»» » writer goroutine is running «")
	defer fmt.Println("»» » conn writer exit « ", c.RemoteAddr().String())

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("»» send buf data error:, ", err, " conn writer exit")
					return
				}
			} else {
				fmt.Println("»» msg buf chan is closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	fmt.Println("»» » reader goroutine is running «")
	defer fmt.Println("»» » conn reader exit « ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			//-<=======================================================================
			// Self Define Vehicle Terminals
			// Message Header & Message Border
			// Sample
			// - GPS DATA :7e 0200 0019 100000716234 0003 04000000 00000108 011e8b38 0457c217 0000 03 220809145240 2e 7e
			// - 7e ... 7e 包头包尾
			// = 消息头 [加上包头标识符1字节共12+1=13字节][固定长度]
			// - 0200			[位置汇报][2字节]
			// - 0019			[包长度][2字节]
			// - 100000716234	[终端识别号][6字节]
			// - 0003			[消息流水号][2字节]
			// = 消息体 [字节数浮动, 根据包长度0019+包尾标识符2字节组成]
			// - 04000000		[报警标识位]
			// - 00000108		[状态标识位]
			// - 011e8b38		[纬度]
			// - 0457c217		[经度]
			// - 0000			[速度]
			// - 03				[方向]
			// - 220809145240	[时间]
			// - 2e				[BCC 校验码]
			//-<=======================================================================

			bufPre := make([]byte, 1)
			_, errTag := io.ReadFull(c.Conn, bufPre)
			if errTag != nil {
				go z.SetOErrToZ(fmt.Sprintf("reader Prefix err:%v", errTag.Error()))
				break
			}

			if bytes.Equal(bufPre, []byte{0x7e}) == false {
				if bytes.Equal(bufPre, []byte{0x00}) == false {
					go z.SetOErrToZ(fmt.Sprintf("Tag != 0x7e, sys break, Tag:%v", hex.EncodeToString(bufPre)))
				}
				break
			}

			var buf bytes.Buffer
			c.ioReader(&buf)
			res := buf.Bytes()
			buf.Reset()

			fmt.Println(hex.EncodeToString(res))

			//msgHeaderLen := 12
			//bufSrcHeader := make([]byte, msgHeaderLen)
			//_, errSrcHeader := io.ReadFull(c.Conn, bufSrcHeader)
			//if errSrcHeader != nil {
			//	go z.SetOErrToZ(fmt.Sprintf("reader Msg Source errSrcHeader err:%v", errSrcHeader.Error()))
			//	break
			//}
			//
			//msgDatumxLen, _ := ztools.BytesToIntU(bufSrcHeader[2:4])
			//bufSrcDatumx := make([]byte, msgDatumxLen)
			//_, errSrcDatumx := io.ReadFull(c.Conn, bufSrcDatumx)
			//if errSrcDatumx != nil {
			//	go z.SetOErrToZ(fmt.Sprintf("reader Msg Source errSrcDatumx err:%v", errSrcDatumx.Error()))
			//	break
			//}
			//
			//bufSuf := make([]byte, 1)
			//_, errSuf := io.ReadFull(c.Conn, bufSuf)
			//if errSuf != nil {
			//	go z.SetOErrToZ(fmt.Sprintf("reader Suffix err:%v", errSuf.Error()))
			//	break
			//}

			//bufSource := make([]byte, 1024)
			//numSource, errSource := c.Conn.Read(bufSource)
			//if errSource != nil {
			//	return
			//}
			//
			//c.SendMsgToLogQueue(
			//	hex.EncodeToString(bufSource[1:3]),
			//	hex.EncodeToString(bufSource))
			//
			//arrSource := bytes.Split(bufSource[:numSource], []byte{0x7e})
			//if len(arrSource) < 3 {
			//	break
			//}
			//
			//startNum := 0
			//if len(arrSource[startNum]) != 0 {
			//	arrSource = arrSource[1:]
			//}
			//
			//endedNum := len(arrSource) - 1
			//if len(arrSource[endedNum]) != 0 {
			//	arrSource = arrSource[:endedNum-1]
			//}
			//
			//if len(arrSource) <= 0 {
			//	break
			//}
			//
			//res := make([][]byte, 0)
			//
			//for _, src := range arrSource {
			//	if len(src) != 0 {
			//		res = append(res, src)
			//	}
			//}
			//
			//if len(res) < 1 {
			//	break
			//}
			//
			//for _, v := range res {
			//	if bytes.Contains(v, []byte{0x7d, 0x02}) {
			//		v = bytes.Replace(v, []byte{0x7d, 0x02}, []byte{0x7e}, -1)
			//	}
			//
			//	if bytes.Contains(v, []byte{0x7d, 0x01}) {
			//		v = bytes.Replace(v, []byte{0x7d, 0x01}, []byte{0x7d}, -1)
			//	}
			//
			//	c.SendReqToTaskQueue(c.MsgType(hex.EncodeToString(v[:2])), v, uint32(len(v)))
			//}
		}
	}
}

func (c *Connection) ioReader(buf *bytes.Buffer) {
	bufSrc := make([]byte, 1)
	_, errSrc := io.ReadFull(c.Conn, bufSrc)
	if errSrc != nil {
		go z.SetOErrToZ(fmt.Sprintf("add io buf err:%v", errSrc.Error()))
		return
	}

	if bytes.Equal(bufSrc, []byte{0x7e}) {
		return
	} else {
		_, err := (*buf).Write(bufSrc)
		if err != nil {
			go z.SetOErrToZ(fmt.Sprintf("add io buf err:%v", err.Error()))
		}

		c.ioReader(buf)
	}
}

func (c *Connection) SendMsgToLogQueue(no string, err string) {
	switch no {
	case "9210":
		go z.Set9210ToZ(err)
		return

	default:
		return
	}
}

func (c *Connection) MsgType(no string) uint32 {
	switch no {
	case "0200":
		return uint32(0)

	case "0b05":
		return uint32(1)

	case "0b03":
		return uint32(2)

	case "0b04":
		return uint32(3)

	case "0203":
		return uint32(4)

	case "9210":
		return uint32(5)

	default:
		return uint32(9999)
	}
}

func (c *Connection) SendReqToTaskQueue(msgId uint32, data []byte, dataLen uint32) {
	var msg ziface.IMessage = &Message{}

	msg.SetMsgID(msgId)
	msg.SetData(data)
	msg.SetDataLen(dataLen)

	//得到当前客户端请求的Request数据
	req := Request{
		conn: c,
		msg:  msg,
	}

	if confs.MaxWorkerPoolSize > 0 {
		//已经启动工作池机制，将消息交给Worker处理
		c.MsgHandler.SendMsgToTaskQueue(&req)
	} else {
		//从绑定好的消息和对应的处理方法中执行对应的Handle方法
		go c.MsgHandler.DoMsgHandler(&req)
	}
}

// Start 启动连接，让当前连接开始工作
func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())
	//1 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	//2 开启用于写回客户端数据流程的Goroutine
	go c.StartWriter()
	//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TCPServer.CallOnConnStart(c)

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
	}
}

// Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	c.cancel()
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// SendMsg 直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	if c.isClosed == true {
		return errors.New("»» connection closed when send msg")
	}

	_ = msgID

	//将data封包，并且发送
	//dp := c.TCPServer.Packet()
	//msg, err := dp.Pack(NewMsgPackage(msgID, data))
	//if err != nil {
	//	fmt.Println("Pack error msg ID = ", msgID)
	//	return errors.New("Pack error msg ")
	//}

	msg := data

	//写回客户端
	_, err := c.Conn.Write(msg)
	return err
}

// SendBuffMsg  发生BuffMsg
func (c *Connection) SendBuffMsg(msgID uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()
	idleTimeout := time.NewTimer(500 * time.Millisecond)
	defer idleTimeout.Stop()

	if c.isClosed == true {
		return errors.New("»» connection closed when send buff msg")
	}

	_ = msgID

	//将data封包，并且发送
	//dp := c.TCPServer.Packet()
	//msg, err := dp.Pack(NewMsgPackage(msgID, data))
	//if err != nil {
	//	fmt.Println("Pack error msg ID = ", msgID)
	//	return errors.New("Pack error msg ")
	//}

	msg := data

	// 发送超时
	select {
	case <-idleTimeout.C:
		return errors.New("»» send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}
	//写回客户端
	//c.msgBuffChan <- msg

	//return nil
}

// SetProperty 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

// GetProperty 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("»» no property found")
}

// RemoveProperty 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

// Context 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *Connection) Context() context.Context {
	return c.ctx
}

func (c *Connection) finalizer() {
	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TCPServer.CallOnConnStop(c)

	c.Lock()
	defer c.Unlock()

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	fmt.Println("»» conn stop()... connID = ", c.ConnID)

	// 关闭socket链接
	_ = c.Conn.Close()

	//将链接从连接管理器中删除
	c.TCPServer.GetConnMgr().Remove(c)

	//关闭该链接全部管道
	close(c.msgBuffChan)
	//设置标志位
	c.isClosed = true
}
