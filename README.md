# socket_asp
socket 自定义包

## 使用方法

### 引用包
```go
import (
	so "socket_asp"
)
```

### 主接口INetIO实现 
``` go
// socket_asp 主接口
type INetIO interface {
	InitConn(obj *NetIO)                    // 连接成功后初始化操作
	Disconn()                               // 断开连接后操作
	ProcessRead(buf []byte) bool            // 对读到的信息进行外部处理
	ProcessWrite(buf []byte) ([]byte, bool) // 对返回的信息进行外部处理
}
```
>接口实现
```go
 // 实现IConnEvent接口
type Achieve struct {
	*Event
}
// 初始化连接执行
func (obj *Achieve) InitConn(param *so.NetIO) {
	obj.NetIO = param // 把连接对象赋值给当前对象
	// 对实现IConnEvent接口方法，进行调用
	if obj.ConnEvent != nil {
		obj.ConnEvent.Connection()
	}
}

// 断开连接时执行
func (obj *Achieve) Disconn() {
	// 对实现IConnEvent接口方法，进行调用
	if obj.ConnEvent != nil {
		obj.ConnEvent.Disconnect()
	}
}

// 读取到消息时执行
func (obj *Achieve) ProcessRead(buf []byte) bool {
	// 对实现IConnEvent接口方法，进行调用
	if obj.ConnEvent != nil {
		obj.ConnEvent.Reading(buf)
	}
	read := &ReqMessage{}
	json.Unmarshal(buf, read)

	flag := true
	func() {
		defer func() {
			if err := recover(); err != nil {
				flag = false
				obj.Close() // 断开连接
			}
		}()
		obj.Request <- read
	}()
	return flag
}

// 发送消息时执行
func (obj *Achieve) ProcessWrite(buf []byte) ([]byte, bool) {
	// 对实现IConnEvent接口方法，进行调用
	if obj.ConnEvent != nil {
		obj.ConnEvent.Writing(buf)
	}
	return buf, true
}
```

### 监听连接

```go
import (
	"fmt"
	"os"
	so "socket_asp"
)

listener,err := so.Listen(network, address)
if err != nil {
	fmt.Printf("打开服务器%s端口: %s 失败", network, address)
	os.Exit(0)
}
event := NewEvent(inter, 30) // 创建事件对象 30秒客户端无连接自动断开
achieve := &Achieve{event}   // 创建NetIO主接口实现对象
obj.Listener.Accept(achieve) // 等待连接
event.Init()                 // 初始化事件对象

go event.WaitRead()  // 等待客户端消息
go event.Events() // 当前连接事件处理

```

### 端口连接
```go
import (
	so "socket_asp"
)
event := NewEvent(inter, 30) // 创建事件对象
achieve := &Achieve{event}   // 创建NetIO主接口实现对象
so.DailTcp(achieve, address) // 主动连接
event.Init()                 // 初始化事件对象

go event.WaitRead()
go event.Events()
```
