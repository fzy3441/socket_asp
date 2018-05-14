package socket_asp

// 主接口
type INetIO interface {
	InitConn(obj *NetIO)                    // 连接成功后初始化操作
	Disconn()                               // 断开连接后操作
	ProcessRead(buf []byte) bool            // 对读到的信息进行外部处理
	ProcessWrite(buf []byte) ([]byte, bool) // 对返回的信息进行外部处理
}
