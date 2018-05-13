package socket_asp

import (
	"convert"
	"encoding/json"
	"errors"
	"net"
)

type NetIO struct {
	Inter   INetIO
	Conn    net.Conn
	DisConn bool
}

// 创建连接
func NewConn(inter INetIO, conn net.Conn) *NetIO {
	netio := &NetIO{Inter: inter, Conn: conn}
	if netio.Inter != nil {
		netio.Inter.Disconn(netio)
	}
	return netio
}

// tcp 连接创建
func DailTcp(inter INetIO, address string) (*NetIO, error) {
	conn, err := net.Dial("tcp", address)
	netio := &NetIO{inter, conn, false}
	return netio, err
}

func (obj *NetIO) Read(process func(buf []byte) bool) (bool, error) {
	buf := make([]byte, 1024)
	tmpBuf := make([]byte, 0)
	for {
		// 连接是不否中断
		if obj.DisConn {
			if obj.Inter != nil {
				obj.Inter.InitConn(obj)
			}
			return false, nil
		}

		_, err := obj.Conn.Read(buf)

		// 当彰连接己断开
		if err != nil && err.Error() == "EOF" {
			obj.DisConn = true
			continue
		}
		// 消息为空
		if buf == nil && err == nil {
			continue
		}

		tmpBuf = append(tmpBuf, buf...)

		for l, err := convert.Bytes2Dec(tmpBuf[0:4]); err == nil && (l+4) <= len(tmpBuf); l, err = convert.Bytes2Dec(tmpBuf[0:4]) {
			if l <= 0 {
				tmpBuf = make([]byte, 0)
				break
			}
			buff := tmpBuf[4 : 4+l]

			if obj.Inter != nil {
				flag := obj.Inter.ProcessRead(buff)
				if !flag {
					return false, errors.New("args function \"ProcessRead\" execution error")
				}
			}

			tmpBuf = tmpBuf[4+l:]
		}
	}
	return false, nil
}

// 发送消息
func (obj *NetIO) Write(b []byte) error {
	if obj.Inter != nil {
		flag := obj.Inter.ProcessWrite(b)
		if !flag {
			return errors.New("args function \"ProcessWrite\" execution error")
		}
	}

	l := len(b)
	head, _ := convert.Int2Bytes(l)
	write := append(head, b...)
	_, err := obj.Conn.Write(write)
	return err
}

// 输出字符串
func (obj *NetIO) WriteString(str string) {
	obj.Write([]byte(str))
}

// 输出对象并转换为json格式
func (obj *NetIO) WriteJson(param interface{}) error {
	buf, err := json.Marshal(param)
	if err == nil {
		obj.Write(buf)
	}
	return err
}

// 断开连接
func (obj *NetIO) Close() {
	obj.Conn.Close()
}
