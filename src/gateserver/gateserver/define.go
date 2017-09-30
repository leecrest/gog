package gateserver

import "net"


type GateCfg struct {
	NetID byte
	ServerID uint32			// 服务器组编号
	WorkSpace string		// 服务器组的工作路径
	LogPath string 			// 日志文件路径，与工作路径的相对路径
	RemoteAddr string		// 对外接口
	LocalAddr string		// 对内接口
	LogPrint bool
}


// 客户端链接定义
type ClientConn struct {
	vfd uint32				// 链接编号
	gameID byte				// 所属服务器编号
	seed uint32				// 种子
	conn net.Conn			// 链接对象
	readBuff []byte			// 读缓冲区
	readSize int			// 读缓冲区大小
	writeBuff []byte		// 写缓冲区
	writeSize int			// 写缓冲区大小
}


// 游戏服链接对象
type GameConn struct {
	id byte
	conn net.Conn			// 链接对象
	readBuff []byte			// 读缓冲区
	readSize int			// 读缓冲区大小
	writeBuff []byte		// 写缓冲区
	writeSize int			// 写缓冲区大小
}


// 常量
const MAX_VFD uint32 = 102400
const READ_BUFF_SIZE int = 1024		// 读取来自客户端的请求的缓冲区大小
const WRITE_BUFF_SIZE int = 10240	// 发送到客户端的数据的缓冲区大小
const CLIENT_PACK_HEAD int = 2		// 包头，2字节表现包体的长度
const CLIENT_PACK_MIN int = 3		// 包体最小长度




// 进程之间通信的数据结构
// 包头第1个数据：cmd byte		// 指令类型
// 包头第2个数据：game byte		// 操作gameid
// 包头第3个数据：vfd uint32		// 操作vfd
// 包头第4个数据：data uint32		// 附加数据长度
const CMD_HEAD_SIZE int = 10

// 进程节点之间的通信指令类型
const CMD_N2G_SYNCGSID byte = 0x01	// 由net发送给game，同步gameid
const CMD_N2G_VFD_ADD byte = 0x02	// 通知game，有新的链接需要管理
const CMD_N2G_VFD_DEL byte = 0x03	// 通知game，链接断开
const CMD_N2G_VFD_DATA byte = 0x04	// 通知game，链接收到数据
const CMD_N2G_VFD_CLOSE byte = 0x05	// 通知game，链接断开
const CMD_G2N_VFD_GSID byte = 0x10	// 改变vfd到gsid的映射，新收到的协议将直接发送到gsid中
const CMD_G2N_VFD_SEND byte = 0x11	// 发送数据给vfd
const CMD_G2N_VFD_SENDS byte = 0x12	// 发送广播数据



// 从vfd中获取netid
func VFD2NID(vfd uint32) (byte) {
	return (byte)(vfd >> 30)
}

// 查询vfd在net中的序号
func VFD2IDX(vfd uint32) (uint32) {
	var tmp uint32 = 0xFFFF
	return (tmp >> 2) & vfd
}